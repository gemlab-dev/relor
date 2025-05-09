package workflow

import (
	"context"
	"time"

	pb "github.com/gemlab-dev/relor/gen/pb/api"
	"github.com/gemlab-dev/relor/internal/graphviz"
	"github.com/gemlab-dev/relor/internal/model"
	"github.com/gemlab-dev/relor/internal/storage"
	"github.com/gemlab-dev/relor/internal/storage/kv"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Logger interface {
	InfoContext(ctx context.Context, msg string, args ...any)
}

type Server struct {
	pb.UnimplementedWorkflowServiceServer

	logger Logger
	store  *kv.WorkflowStorage
}

func New(l Logger, s *kv.WorkflowStorage) *Server {
	return &Server{
		logger: l,
		store:  s,
	}
}

func (s *Server) Run(ctx context.Context, in *pb.RunRequest) (*pb.RunResponse, error) {
	g := model.Graph{}
	if err := g.FromProto(in.Graph); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse graph: %v", err)
	}
	if in.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}
	w := model.NewWorkflow(id, &g, time.Now())
	s.logger.InfoContext(ctx, "Running workflow", "w", w)
	if err := s.store.CreateWorkflow(ctx, *w); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create workflow: %v", err)
	}
	return &pb.RunResponse{}, nil
}

func (s *Server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	if in.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}
	w, err := s.store.GetWorkflow(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get workflow: %v", err)
	}

	return &pb.GetResponse{
		State: &pb.WorkflowState{
			Status:      string(w.Status),
			CurrentNode: w.CurrentNode(),
		},
		TransitionId: w.LastTransitionId(),
	}, nil
}

func (s *Server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	if in.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}
	if in.ResultLabel == "" {
		return nil, status.Errorf(codes.InvalidArgument, "result label is empty")
	}
	na := storage.NextAction{
		ID:                  id,
		Label:               in.ResultLabel,
		LastKnownTransition: in.TransitionId,
	}
	if err := s.store.UpdateNextAction(ctx, na); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update next action: %v", err)
	}
	return &pb.UpdateResponse{}, nil
}

func (s *Server) History(ctx context.Context, in *pb.GetRequest) (*pb.HistoryResponse, error) {
	if in.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is empty")
	}
	wid, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse id: %v", err)
	}
	w, err := s.store.GetWorkflow(ctx, wid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get workflow: %v", err)
	}
	th, err := s.store.GetHistory(ctx, wid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get transition history: %v", err)
	}
	var ts []*pb.Transition
	ith := th
	for ith != nil {
		t := &pb.Transition{
			Label:    ith.Label(),
			Walltime: durationpb.New(ith.Walltime()),
		}
		t.From, t.To = ith.FromTo()
		ts = append(ts, t)
		ith = ith.Next()
	}

	gv, err := graphviz.Dot(*w, th)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate graphviz: %v", err)
	}

	return &pb.HistoryResponse{
		Transitions: ts,
		Graphviz:    gv,
	}, nil
}
