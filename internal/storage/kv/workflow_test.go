package kv

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	pb "github.com/gemlab-dev/relor/gen/pb/graph"
	"github.com/gemlab-dev/relor/internal/model"
	"github.com/gemlab-dev/relor/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/prototext"
)

func TestWorkflowKVStorageOpenFailure(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Errorf("failed to remove temp file: %v", err)
		}
	}()

	kv1, err := NewWorkflowStorage(tempFile.Name(), "testBucket", time.Now)
	if err != nil {
		t.Fatalf("failed to initialize storage: %v", err)
	}
	defer func() {
		if err := kv1.Close(); err != nil {
			t.Errorf("failed to close storage: %v", err)
		}
	}()
	// Attempt to open multiple instances of the same file.
	kv2, err := NewWorkflowStorage(tempFile.Name(), "testBucket", time.Now)
	if err == nil {
		t.Fatalf("expected error when opening multiple instances, got nil")
	}
	if kv2 != nil {
		t.Fatalf("expected nil storage instance, got %v", kv2)
	}
}

func TestWorkflowKVStorage(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Errorf("failed to remove temp file: %v", err)
		}
	}()

	ts := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	now := func() time.Time {
		return ts
	}

	const testBucket = "testBucket"
	kv, err := NewWorkflowStorage(tempFile.Name(), testBucket, now)
	if err != nil {
		t.Fatalf("failed to initialize storage: %v", err)
	}
	defer func() {
		if err := kv.Close(); err != nil {
			t.Errorf("failed to close storage: %v", err)
		}
	}()

	txt := `
		start: "a"
		nodes { id: "a"  }
		nodes { id: "b"  }
		edges { from_id: "a" to_id: "b" condition { operation_result: "ok" } }
		edges { from_id: "a" to_id: "a" condition { operation_result: "retry" } }
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

	t.Run("Get Workflow before creating any workflows", func(t *testing.T) {
		ctx := context.Background()
		wf, err := kv.GetWorkflow(ctx, workflowID)
		if err != nil {
			t.Fatalf("failed to retrieve workflow: %v", err)
		}
		if wf != nil {
			t.Fatalf("expected nil workflow, got %v", wf)
		}
	})

	t.Run("Get Workflow with wrong bucket", func(t *testing.T) {
		ctx := context.Background()
		kv.bucketName = "wrongBucket"
		wf, err := kv.GetWorkflow(ctx, workflowID)
		if err == nil {
			t.Fatalf("expected error for wrong bucket, got nil")
		}
		if wf != nil {
			t.Fatalf("expected nil workflow, got %v", wf)
		}
	})
	kv.bucketName = testBucket

	t.Run("Create and Get Workflow", func(t *testing.T) {
		workflow := model.NewWorkflow(workflowID, g, ts)
		ctx := context.Background()

		// This will create the bucket as well.
		if err := kv.CreateWorkflow(ctx, *workflow); err != nil {
			t.Fatalf("failed to create workflow: %v", err)
		}
		na := storage.NextAction{
			ID:                  workflowID,
			Label:               "retry",
			LastKnownTransition: 0,
		}
		if err := kv.UpdateNextAction(ctx, na); err != nil {
			t.Fatalf("failed to update next action: %v", err)
		}
		na.Label = "ok"
		na.LastKnownTransition = 1
		if err := kv.UpdateNextAction(ctx, na); err != nil {
			t.Fatalf("failed to update next action: %v", err)
		}

		retrievedWorkflow, err := kv.GetWorkflow(ctx, workflowID)
		if err != nil {
			t.Fatalf("failed to retrieve workflow: %v", err)
		}
		if retrievedWorkflow.ID != workflow.ID {
			t.Errorf("expected workflow ID %v, got %v", workflow.ID, retrievedWorkflow.ID)
		}
		if retrievedWorkflow.CurrentNode() != "b" {
			t.Errorf("expected current node b, got %v", retrievedWorkflow.CurrentNode())
		}
		if retrievedWorkflow.Status != model.WorkflowStatusCompleted {
			t.Errorf("expected status %v, got %v", model.WorkflowStatusCompleted, retrievedWorkflow.Status)
		}
		if retrievedWorkflow.LastTransitionId() != 2 {
			t.Errorf("expected last transition ID 2, got %v", retrievedWorkflow.LastTransitionId())
		}

		tn, err := kv.GetHistory(ctx, workflowID)
		if err != nil {
			t.Fatalf("failed to get history: %v", err)
		}
		if tn.Next() == nil {
			t.Fatal("expected next transition, got nil")
		}
		if tn.Next().Label() != "retry" {
			t.Errorf("expected label 'retry', got %v", tn.Next().Label())
		}
		if tn.Next().Next() != nil {
			t.Fatal("expected no next transition, got non-nil")
		}
		if tn.Label() != "ok" {
			t.Errorf("expected label 'ok', got %v", tn.Label())
		}
	})

	t.Run("Create Duplicate Workflow", func(t *testing.T) {
		ctx := context.Background()
		workflowID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		workflow := model.NewWorkflow(workflowID, g, ts)
		// Attempt to create a duplicate workflow.
		err := kv.CreateWorkflow(ctx, *workflow)
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

		wf, err := kv.GetWorkflow(ctx, nonExistentID)
		if err != nil {
			t.Fatalf("failed to retrieve workflow: %v", err)
		}
		if wf != nil {
			t.Fatalf("expected nil workflow, got %v", wf)
		}
	})

	t.Run("Update next action time", func(t *testing.T) {
		ctx := context.Background()

		if w, err := kv.GetWorkflow(ctx, workflowID); err != nil {
			t.Fatalf("failed to retrieve workflow: %v", err)
		} else if w.NextActionAt.IsZero() {
			t.Fatalf("expected next action time to be set, got zero value")
		} else if w.NextActionAt.Sub(ts) > 1*time.Second {
			t.Errorf("expected next action time to be %v, got %v", ts, w.NextActionAt)
		}

		if err := kv.UpdateTimeout(ctx, workflowID, 10*time.Second); err != nil {
			t.Fatalf("failed to update timeout: %v", err)
		}
		w, err := kv.GetWorkflow(ctx, workflowID)
		if err != nil {
			t.Fatalf("failed to retrieve workflow: %v", err)
		}
		if w.NextActionAt.Sub(ts.Add(10*time.Second)) > 1*time.Second {
			t.Errorf("expected next action time to be %v, got %v", ts.Add(10*time.Second), w.NextActionAt)
		}
	})

	t.Run("Update non-existent next action time", func(t *testing.T) {
		ctx := context.Background()
		nonExistentID := uuid.New()

		if err := kv.UpdateTimeout(ctx, nonExistentID, 10*time.Second); err == nil {
			t.Fatalf("expected error for non-existent workflow, got nil")
		}
	})
}
