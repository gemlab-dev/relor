package model

import (
	"testing"

	wpb "github.com/gemlab-dev/relor/gen/pb/db"
	gpb "github.com/gemlab-dev/relor/gen/pb/graph"
	"google.golang.org/protobuf/encoding/prototext"
)

func TestWorkflowFromProto(t *testing.T) {
	txt := `
	start: "a"
	nodes { id: "a"  }
	nodes { id: "b"  }
	edges { from_id: "a" to_id: "b" condition { operation_result: "ok" } }
	edges { from_id: "a" to_id: "a" condition { operation_result: "retry" } }
	`
	pb := &wpb.Workflow{
		Id:     "00000000-0000-0000-0000-000000000001",
		Graph:  &gpb.Graph{},
		Status: wpb.WorkflowStatus_WORKFLOW_STATUS_RUNNING,
	}
	if err := prototext.Unmarshal([]byte(txt), pb.Graph); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	wf := &Workflow{}

	// Default next node.
	if err := wf.FromProto(pb, "", 0); err != nil {
		t.Fatalf("failed to convert from proto: %v", err)
	}
	if wf.CurrentNode() != "a" {
		t.Errorf("expected current node to be 'a', got %s", wf.CurrentNode())
	}
	if wf.LastTransitionId() != 0 {
		t.Errorf("expected last transition ID to be 0, got %d", wf.LastTransitionId())
	}
	if wf.Status != WorkflowStatusRunning {
		t.Errorf("expected status to be 'running', got %s", wf.Status)
	}

	// Specified next node.
	if err := wf.FromProto(pb, "b", 1); err != nil {
		t.Fatalf("failed to convert from proto: %v", err)
	}
	if wf.CurrentNode() != "b" {
		t.Errorf("expected current node to be 'b', got %s", wf.CurrentNode())
	}
	if wf.LastTransitionId() != 1 {
		t.Errorf("expected last transition ID to be 1, got %d", wf.LastTransitionId())
	}
}

func TestWorkflowFromProtoUUIDError(t *testing.T) {
	txt := `
	start: "a"
	nodes { id: "a"  }
	`
	pb := &wpb.Workflow{
		Id:    "broken",
		Graph: &gpb.Graph{},
	}
	if err := prototext.Unmarshal([]byte(txt), pb.Graph); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	wf := &Workflow{}

	if err := wf.FromProto(pb, "", 0); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
