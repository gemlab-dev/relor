package kv

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	pb "github.com/gemlab-dev/relor/gen/pb/graph"
	"github.com/gemlab-dev/relor/internal/model"
	"github.com/gemlab-dev/relor/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/prototext"
)

const (
	graphTxt = `
		start: "a"
		nodes { id: "a"  }
		nodes { id: "b"  }
		edges { from_id: "a" to_id: "b" condition { operation_result: "ok" } }
		edges { from_id: "a" to_id: "a" condition { operation_result: "retry" } }
		`
)

func TestGetNextWorkflows(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Errorf("failed to remove temp file: %v", err)
		}
	}()

	// Initialise the storage.
	currentTime := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	now := func() time.Time {
		return currentTime
	}
	kvStore, err := NewWorkflowStorage(tempFile.Name(), "testBucket", now, 10)
	if err != nil {
		t.Fatalf("failed to initialize storage: %v", err)
	}
	defer func() {
		if err := kvStore.Close(); err != nil {
			t.Errorf("failed to close storage: %v", err)
		}
	}()

	// Create a workflow.
	gpb := &pb.Graph{}
	if err := prototext.Unmarshal([]byte(graphTxt), gpb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &model.Graph{}
	if err := g.FromProto(gpb); err != nil {
		t.Fatalf("failed to load graph: %v", err)
	}
	workflowID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	startTime := time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC) // Set a future start time
	workflow := model.NewWorkflow(workflowID, g, startTime)
	ctx := context.Background()

	if err := kvStore.CreateWorkflow(ctx, *workflow); err != nil {
		t.Fatalf("failed to create workflow: %v", err)
	}

	// Scheduled workflow should not be returned before the start time.
	wfs, err := kvStore.GetNextWorkflows(ctx)
	if err != nil {
		t.Fatalf("failed to get next workflows: %v", err)
	}
	if len(wfs) != 0 {
		t.Fatalf("expected no workflows, got %d", len(wfs))
	}

	// Move time forward to schedule the workflow.
	currentTime = time.Date(2025, time.January, 3, 0, 0, 0, 0, time.UTC)

	wfs, err = kvStore.GetNextWorkflows(ctx)
	if err != nil {
		t.Fatalf("failed to get next workflows: %v", err)
	}
	if len(wfs) != 1 {
		t.Fatalf("expected 1 workflow, got %d", len(wfs))
	}
	if wfs[0].ID != workflowID {
		t.Fatalf("expected workflow ID %s, got %s", workflowID, wfs[0].ID)
	}

	// Update the workflow's timeout to a future date.
	err = kvStore.UpdateTimeout(ctx, workflowID, 23*time.Hour)
	if err != nil {
		t.Fatalf("failed to update timeout: %v", err)
	}
	wfs, err = kvStore.GetNextWorkflows(ctx)
	if err != nil {
		t.Fatalf("failed to get next workflows: %v", err)
	}
	if len(wfs) != 0 {
		t.Fatalf("expected no workflows, got %d", len(wfs))
	}

	// Move time forward again and check if the workflow is returned.
	currentTime = time.Date(2025, time.January, 4, 0, 0, 0, 0, time.UTC)
	wfs, err = kvStore.GetNextWorkflows(ctx)
	if err != nil {
		t.Fatalf("failed to get next workflows: %v", err)
	}
	if len(wfs) != 1 {
		t.Fatalf("expected 1 workflow, got %d", len(wfs))
	}
	if wfs[0].ID != workflowID {
		t.Fatalf("expected workflow ID %s, got %s", workflowID, wfs[0].ID)
	}

	// Move time forward again and check if the workflow is still returned.
	currentTime = time.Date(2025, time.January, 5, 0, 0, 0, 0, time.UTC)
	wfs, err = kvStore.GetNextWorkflows(ctx)
	if err != nil {
		t.Fatalf("failed to get next workflows: %v", err)
	}
	if len(wfs) != 1 {
		t.Fatalf("expected 1 workflow, got %d", len(wfs))
	}
	if wfs[0].ID != workflowID {
		t.Fatalf("expected workflow ID %s, got %s", workflowID, wfs[0].ID)
	}
}

func TestGetNextWorkflowsBatchSize(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Errorf("failed to remove temp file: %v", err)
		}
	}()

	// Initialise the storage.
	currentTime := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	now := func() time.Time {
		return currentTime
	}
	kvStore, err := NewWorkflowStorage(tempFile.Name(), "testBucket", now, 2)
	if err != nil {
		t.Fatalf("failed to initialize storage: %v", err)
	}
	defer func() {
		if err := kvStore.Close(); err != nil {
			t.Errorf("failed to close storage: %v", err)
		}
	}()

	wfs, err := kvStore.GetNextWorkflows(context.Background())
	if err != nil {
		t.Fatalf("failed to get next workflows: %v", err)
	}
	if len(wfs) != 0 {
		t.Fatalf("expected no workflows, got %d", len(wfs))
	}

	gpb := &pb.Graph{}
	if err := prototext.Unmarshal([]byte(graphTxt), gpb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &model.Graph{}
	if err := g.FromProto(gpb); err != nil {
		t.Fatalf("failed to load graph: %v", err)
	}

	for i := 0; i < 5; i++ {
		ctx := context.Background()
		workflowID := uuid.MustParse("00000000-0000-0000-0000-00000000000" + strconv.Itoa(i+1))
		startTime := time.Date(2024, time.January, 1+i, 0, 0, 0, 0, time.UTC) // Set past start time
		workflow := model.NewWorkflow(workflowID, g, startTime)
		if err := kvStore.CreateWorkflow(ctx, *workflow); err != nil {
			t.Fatalf("failed to create workflow: %v", err)
		}
	}
	ctx := context.Background()
	wfs, err = kvStore.GetNextWorkflows(ctx)
	if err != nil {
		t.Fatalf("failed to get next workflows: %v", err)
	}
	if len(wfs) != 2 {
		t.Fatalf("expected 2 workflows, got %d", len(wfs))
	}
	if wfs[0].ID != uuid.MustParse("00000000-0000-0000-0000-000000000001") {
		t.Fatalf("expected workflow ID %s, got %s", "00000000-0000-0000-0000-000000000001", wfs[0].ID)
	}
	if wfs[1].ID != uuid.MustParse("00000000-0000-0000-0000-000000000002") {
		t.Fatalf("expected workflow ID %s, got %s", "00000000-0000-0000-0000-000000000002", wfs[1].ID)
	}

	err = kvStore.UpdateTimeout(ctx, uuid.MustParse("00000000-0000-0000-0000-000000000001"), 23*time.Hour)
	if err != nil {
		t.Fatalf("failed to update timeout: %v", err)
	}
	wfs, err = kvStore.GetNextWorkflows(ctx)
	if err != nil {
		t.Fatalf("failed to get next workflows: %v", err)
	}
	if len(wfs) != 2 {
		t.Fatalf("expected 2 workflows, got %d", len(wfs))
	}
	if wfs[0].ID != uuid.MustParse("00000000-0000-0000-0000-000000000002") {
		t.Fatalf("expected workflow ID %s, got %s", "00000000-0000-0000-0000-000000000002", wfs[0].ID)
	}
	if wfs[1].ID != uuid.MustParse("00000000-0000-0000-0000-000000000003") {
		t.Fatalf("expected workflow ID %s, got %s", "00000000-0000-0000-0000-000000000003", wfs[1].ID)
	}
}

func TestWorkflowCompletion(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Errorf("failed to remove temp file: %v", err)
		}
	}()

	// Initialise the storage.
	currentTime := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	now := func() time.Time {
		return currentTime
	}
	kvStore, err := NewWorkflowStorage(tempFile.Name(), "testBucket", now, 10)
	if err != nil {
		t.Fatalf("failed to initialize storage: %v", err)
	}
	defer func() {
		if err := kvStore.Close(); err != nil {
			t.Errorf("failed to close storage: %v", err)
		}
	}()

	gpb := &pb.Graph{}
	if err := prototext.Unmarshal([]byte(graphTxt), gpb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &model.Graph{}
	if err := g.FromProto(gpb); err != nil {
		t.Fatalf("failed to load graph: %v", err)
	}

	ctx := context.Background()

	wfs, err := kvStore.GetNextWorkflows(ctx)
	if err != nil {
		t.Fatalf("failed to get next workflows: %v", err)
	}
	if len(wfs) != 0 {
		t.Fatalf("expected no workflows, got %d", len(wfs))
	}

	workflowID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	startTime := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC) // Set past start time
	workflow := model.NewWorkflow(workflowID, g, startTime)
	if err := kvStore.CreateWorkflow(ctx, *workflow); err != nil {
		t.Fatalf("failed to create workflow: %v", err)
	}

	wfs, err = kvStore.GetNextWorkflows(ctx)
	if err != nil {
		t.Fatalf("failed to get next workflows: %v", err)
	}
	if len(wfs) != 1 {
		t.Fatalf("expected 1 workflow, got %d", len(wfs))
	}

	if err := kvStore.UpdateTimeout(ctx, workflowID, 1*time.Minute); err != nil {
		t.Fatalf("failed to update timeout: %v", err)
	}
	na := storage.NextAction{
		ID:                  workflowID,
		Label:               "retry",
		LastKnownTransition: 0,
	}
	if err := kvStore.UpdateNextAction(ctx, na); err != nil {
		t.Fatalf("failed to update next action: %v", err)
	}

	currentTime = currentTime.Add(2 * time.Minute)
	wfs, err = kvStore.GetNextWorkflows(ctx)
	if err != nil {
		t.Fatalf("failed to get next workflows: %v", err)
	}
	if len(wfs) != 1 {
		t.Fatalf("expected 1 workflow, got %d", len(wfs))
	}

	if err := kvStore.UpdateTimeout(ctx, workflowID, 1*time.Minute); err != nil {
		t.Fatalf("failed to update timeout: %v", err)
	}
	na.Label = "ok"
	na.LastKnownTransition = 1
	if err := kvStore.UpdateNextAction(ctx, na); err != nil {
		t.Fatalf("failed to update next action: %v", err)
	}

	currentTime = currentTime.Add(2 * time.Minute)
	wfs, err = kvStore.GetNextWorkflows(ctx)
	if err != nil {
		t.Fatalf("failed to get next workflows: %v", err)
	}
	if len(wfs) != 0 {
		t.Fatalf("expected no workflows, got %d", len(wfs))
	}
}
