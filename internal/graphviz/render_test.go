package graphviz

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/gemlab-dev/relor/internal/model"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/prototext"

	gpb "github.com/gemlab-dev/relor/gen/pb/graph"
)

func TestRenderSVG(t *testing.T) {
	// 1. Setup a workflow and transition history
	workflowID := uuid.New()

	textpb := `
	start: "start"
	nodes: {
		id: "start"
		op: { name: "start_op" }
	}
	nodes: {
		id: "in_progress"
		op: { name: "in_progress_op" }
	}
	nodes: {
		id: "done"
		op: { name: "done_op" }
	}
	edges: {
		from_id: "start"
		to_id: "in_progress"
		condition: { operation_result: "begin" }
	}
	edges: {
		from_id: "in_progress"
		to_id: "done"
		condition: { operation_result: "finish" }
	}
	`
	var pbGraph gpb.Graph
	if err := prototext.Unmarshal([]byte(textpb), &pbGraph); err != nil {
		t.Fatalf("failed to unmarshal text proto: %v", err)
	}

	modelGraph := &model.Graph{}
	if err := modelGraph.FromProto(&pbGraph); err != nil {
		t.Fatalf("failed to create graph from proto: %v", err)
	}
	workflow := model.NewWorkflow(workflowID, modelGraph, time.Now())
	workflow.SetCurrentNode("in_progress")

	history, _ := model.NewTransitionHistory(time.Now(), []model.RawTransition{
		{From: "start", To: "in_progress", Label: "begin"},
	})

	// 2. Call RenderSVG
	svgBytes, err := RenderSVG(context.Background(), *workflow, history)
	if err != nil {
		t.Fatalf("RenderSVG failed: %v", err)
	}

	// 3. Assert the output
	svgString := string(svgBytes)

	if !strings.Contains(svgString, "<svg") {
		t.Error("output does not look like an SVG")
	}

	// Check for nodes
	if !strings.Contains(svgString, ">start</text>") {
		t.Error("start node not found in SVG")
	}
	if !strings.Contains(svgString, ">in_progress</text>") {
		t.Error("in_progress node not found in SVG")
	}
	if !strings.Contains(svgString, ">done</text>") {
		t.Error("done node not found in SVG")
	}

	// Check for edges
	if !strings.Contains(svgString, ">begin (1)</text>") {
		t.Error("begin edge not found in SVG")
	}
	if !strings.Contains(svgString, ">finish</text>") {
		t.Error("finish edge not found in SVG")
	}
}

func TestRenderSVG_NilGraph(t *testing.T) {
	// 1. Setup a workflow with a nil graph
	workflow := model.NewWorkflow(uuid.New(), nil, time.Now())
	history, _ := model.NewTransitionHistory(time.Now(), nil)

	// 2. Call RenderSVG
	_, err := RenderSVG(context.Background(), *workflow, history)

	// 3. Assert that an error is returned
	if err == nil {
		t.Fatal("RenderSVG should have returned an error for a nil graph, but it didn't")
	}
}
