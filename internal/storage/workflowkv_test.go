package storage

import (
	"context"
	"os"
	"strings"
	"testing"

	pb "github.com/gemlab-dev/relor/gen/pb/graph"
	"github.com/gemlab-dev/relor/internal/model"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

func TestWorkflowKVStorage(t *testing.T) {
	// Create a temporary file for the BoltDB database
	tempFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Errorf("failed to remove temp file: %v", err)
		}
	}()

	// Initialize WorkflowKVStorage
	storage, err := NewWorkflowKVStorage(tempFile.Name(), "testBucket")
	if err != nil {
		t.Fatalf("failed to initialize storage: %v", err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			t.Errorf("failed to close storage: %v", err)
		}
	}()

	txt := `
		start: "a"
		nodes { id: "a" }
		`
	gpb := &pb.Graph{}
	if err := prototext.Unmarshal([]byte(txt), gpb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	g := &model.Graph{}
	if err := g.FromProto(gpb); err != nil {
		t.Fatalf("failed to load graph: %v", err)
	}

	workflowID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	t.Run("Create and Get Workflow", func(t *testing.T) {
		workflow := model.Workflow{
			ID:          workflowID,
			CurrentNode: "a",
			Status:      model.WorkflowStatusRunning,
			Graph:       g,
		}
		ctx := context.Background()

		// Create a workflow
		if err := storage.CreateWorkflow(ctx, workflow); err != nil {
			t.Fatalf("failed to create workflow: %v", err)
		}

		// Retrieve the workflow
		retrievedWorkflow, err := storage.GetWorkflow(ctx, workflowID)
		if err != nil {
			t.Fatalf("failed to retrieve workflow: %v", err)
		}
		if retrievedWorkflow.ID != workflow.ID {
			t.Errorf("expected workflow ID %v, got %v", workflow.ID, retrievedWorkflow.ID)
		}
		if retrievedWorkflow.CurrentNode != workflow.CurrentNode {
			t.Errorf("expected current node %v, got %v", workflow.CurrentNode, retrievedWorkflow.CurrentNode)
		}
		if retrievedWorkflow.Status != workflow.Status {
			t.Errorf("expected status %v, got %v", workflow.Status, retrievedWorkflow.Status)
		}

		gotGraphPb, err := retrievedWorkflow.Graph.ToProto()
		if err != nil {
			t.Fatalf("failed to convert graph to proto: %v", err)
		}
		if !proto.Equal(gpb, gotGraphPb) {
			t.Errorf("expected graph %v, got %v", gpb, gotGraphPb)
		}

	})

	t.Run("Create Duplicate Workflow", func(t *testing.T) {
		ctx := context.Background()
		workflowID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		workflow := model.Workflow{
			ID:          workflowID,
			CurrentNode: "a",
			Status:      model.WorkflowStatusRunning,
			Graph:       g,
		}
		// Attempt to create a duplicate workflow
		err := storage.CreateWorkflow(ctx, workflow)
		if err == nil {
			t.Fatalf("expected error for duplicate workflow, got nil")
		}
		if !strings.Contains(err.Error(), "workflow already exists") {
			t.Errorf("expected error 'workflow already exists', got %v", err.Error())
		}
	})

	t.Run("Get Non-Existent Workflow", func(t *testing.T) {
		ctx := context.Background()
		nonExistentID := uuid.New()

		_, err := storage.GetWorkflow(ctx, nonExistentID)
		if err == nil {
			t.Fatalf("expected error for non-existent workflow, got nil")
		}
		if err.Error() != "workflow not found" {
			t.Errorf("expected error 'workflow not found', got %v", err.Error())
		}
	})
}
