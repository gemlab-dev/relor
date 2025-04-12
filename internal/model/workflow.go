package model

import (
	"fmt"

	"github.com/google/uuid"

	pb "github.com/gemlab-dev/relor/gen/pb/db"
)

type WorkflowStatus string

const (
	WorkflowStatusUnknown   WorkflowStatus = "unknown"
	WorkflowStatusPending   WorkflowStatus = "pending"
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusFailed    WorkflowStatus = "failed"
)

type Workflow struct {
	ID               uuid.UUID
	CurrentNode      string
	Status           WorkflowStatus
	Graph            *Graph
	LastTransitionId *int64
	// NextActionAt time.Time
}

func NewWorkflow(id uuid.UUID, g *Graph) *Workflow {
	return &Workflow{
		ID: id,
		// TODO: set status to WorkflowStatusPending to stage workflows before running.
		Status:      WorkflowStatusRunning,
		CurrentNode: g.Head(),
		Graph:       g,
	}
}

func (w Workflow) ToProto() (*pb.Workflow, error) {
	if w.Graph == nil {
		return nil, fmt.Errorf("graph is not initialised")
	}
	pbGraph, err := w.Graph.ToProto()
	if err != nil {
		return nil, fmt.Errorf("failed to convert graph to proto: %w", err)
	}
	return &pb.Workflow{
		Id:               w.ID.String(),
		CurrentNode:      w.CurrentNode,
		Status:           workflowStatusToProto(w.Status),
		Graph:            pbGraph,
		LastTransitionId: w.LastTransitionId,
	}, nil
}

func (w *Workflow) FromProto(pbWorkflow *pb.Workflow) error {
	if w == nil {
		return nil
	}
	var err error
	w.ID, err = uuid.Parse(pbWorkflow.Id)
	if err != nil {
		return fmt.Errorf("failed to parse workflow ID: %w", err)
	}
	w.CurrentNode = pbWorkflow.CurrentNode
	w.Status = protoToWorkflowStatus(pbWorkflow.Status)
	g := &Graph{}
	if err := g.FromProto(pbWorkflow.Graph); err != nil {
		return fmt.Errorf("failed to convert graph from proto: %w", err)
	}
	w.Graph = g
	w.LastTransitionId = pbWorkflow.LastTransitionId
	return nil
}

func protoToWorkflowStatus(status pb.WorkflowStatus) WorkflowStatus {
	switch status {
	case pb.WorkflowStatus_WORKFLOW_STATUS_PENDING:
		return WorkflowStatusPending
	case pb.WorkflowStatus_WORKFLOW_STATUS_RUNNING:
		return WorkflowStatusRunning
	case pb.WorkflowStatus_WORKFLOW_STATUS_COMPLETED:
		return WorkflowStatusCompleted
	case pb.WorkflowStatus_WORKFLOW_STATUS_FAILED:
		return WorkflowStatusFailed
	default:
		return WorkflowStatusUnknown
	}
}

func workflowStatusToProto(status WorkflowStatus) pb.WorkflowStatus {
	switch status {
	case WorkflowStatusPending:
		return pb.WorkflowStatus_WORKFLOW_STATUS_PENDING
	case WorkflowStatusRunning:
		return pb.WorkflowStatus_WORKFLOW_STATUS_RUNNING
	case WorkflowStatusCompleted:
		return pb.WorkflowStatus_WORKFLOW_STATUS_COMPLETED
	case WorkflowStatusFailed:
		return pb.WorkflowStatus_WORKFLOW_STATUS_FAILED
	default:
		return pb.WorkflowStatus_WORKFLOW_STATUS_UNSPECIFIED
	}
}
