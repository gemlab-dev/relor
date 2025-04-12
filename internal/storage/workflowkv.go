package storage

import (
	"context"
	"fmt"
	"time"

	pb "github.com/gemlab-dev/relor/gen/pb/db"
	"github.com/gemlab-dev/relor/internal/model"
	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/encoding/protojson"
)

const defaultIndex = 0
const defaultNamespace = 0

type WorkflowKVStorage struct {
	db         *bolt.DB
	bucketName string
}

func NewWorkflowKVStorage(path, bucketName string) (*WorkflowKVStorage, error) {
	opts := &bolt.Options{Timeout: 1 * time.Second}
	db, err := bolt.Open(path, 0600, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open KV store: %v", err)
	}

	return &WorkflowKVStorage{db: db, bucketName: bucketName}, nil
}

func (w *WorkflowKVStorage) Close() error {
	return w.db.Close()
}

func (s *WorkflowKVStorage) CreateWorkflow(ctx context.Context, w model.Workflow) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(s.bucketName))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}

		return s.saveNewWorkflow(b, w)
	})
}

func (s *WorkflowKVStorage) GetWorkflow(ctx context.Context, id uuid.UUID) (*model.Workflow, error) {
	var result *model.Workflow
	f := func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		var err error
		result, err = s.loadWorkflow(b, id)
		if err != nil {
			return err
		}
		return nil
	}

	if err := s.db.View(f); err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	return result, nil
}

func (s *WorkflowKVStorage) UpdateNextAction(ctx context.Context, na NextAction) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		if b == nil {
			return fmt.Errorf("bucket not found: %s", s.bucketName)
		}

		wf, err := s.loadWorkflow(b, na.ID)
		if err != nil {
			return fmt.Errorf("failed to load workflow: %w", err)
		}
		if wf == nil {
			return fmt.Errorf("workflow not found: %s", na.ID)
		}

		// Validate current action and get next node.
		nextNode, err := wf.Graph.NextNodeID(wf.CurrentNode, na.Label)
		if err != nil {
			return fmt.Errorf("failed to get next node: %w", err)
		}

		newTnPb := pb.WorkflowTransition{
			WorkflowId: wf.ID.String(),
			FromNode:   wf.CurrentNode,
			ToNode:     nextNode,
			Label:      na.Label,
			Timestamp:  time.Now().Unix(),
		}
		if wf.LastTransitionId != nil {
			tnKey := transitionKey(defaultNamespace, na.ID, *wf.LastTransitionId)
			tnData := b.Get(tnKey)
			if tnData == nil {
				return fmt.Errorf("transition not found; key: %s", tnKey)
			}

			var tnpb pb.WorkflowTransition
			err := protojson.Unmarshal(tnData, &tnpb)
			if err != nil {
				return fmt.Errorf("failed to unmarshal transition proto: %v", err)
			}
			newTnPb.Id = tnpb.Id + 1
		}
		// Write transition
		tnKey := transitionKey(defaultNamespace, na.ID, newTnPb.Id)
		tnData, err := protojson.Marshal(&newTnPb)
		if err != nil {
			return fmt.Errorf("failed to marshal transition proto: %v", err)
		}
		if err := b.Put(tnKey, tnData); err != nil {
			return fmt.Errorf("failed to save transition: %v", err)
		}

		// Update workflow
		wf.LastTransitionId = &newTnPb.Id
		wf.CurrentNode = nextNode

		nextLabels, err := wf.Graph.OutLabels(nextNode)
		if err != nil {
			return fmt.Errorf("failed to get out labels: %w", err)
		}
		if len(nextLabels) == 0 {
			wf.Status = model.WorkflowStatusCompleted
		}

		return s.saveNewWorkflow(b, *wf)
	})
}

func (s *WorkflowKVStorage) GetHistory(ctx context.Context, workflowID uuid.UUID) (*model.Transition, error) {
	return nil, nil
}

func (s *WorkflowKVStorage) UpdateTimeout(ctx context.Context, id uuid.UUID, timeout time.Duration) error {
	return nil
}

func (s *WorkflowKVStorage) GetNextWorkflows(ctx context.Context) ([]model.Workflow, error) {
	var workflows []model.Workflow
	return workflows, nil
}

func (s *WorkflowKVStorage) loadWorkflow(b *bolt.Bucket, id uuid.UUID) (*model.Workflow, error) {
	key := workflowKey(defaultNamespace, id)
	data := b.Get(key)
	if data == nil {
		return nil, fmt.Errorf("workflow not found")
	}

	var wpb pb.Workflow
	err := protojson.Unmarshal(data, &wpb)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow proto: %v", err)
	}
	var w model.Workflow
	if err := w.FromProto(&wpb); err != nil {
		return nil, fmt.Errorf("failed to convert workflow from proto: %w", err)
	}
	return &w, nil
}

func (s *WorkflowKVStorage) saveNewWorkflow(b *bolt.Bucket, w model.Workflow) error {
	wpb, err := w.ToProto()
	if err != nil {
		return fmt.Errorf("failed to convert workflow to proto: %w", err)
	}
	data, err := protojson.Marshal(wpb)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow proto to json: %w", err)
	}

	key := workflowKey(defaultNamespace, w.ID)

	if d := b.Get(key); d != nil {
		return fmt.Errorf("workflow already exists: %s", w.ID)
	}

	return b.Put(key, data)
}

func workflowKey(namespace int, id uuid.UUID) []byte {
	return fmt.Appendf(nil, "%d/%d/%s", namespace, defaultIndex, id.String())
}

func transitionKey(namespace int, id uuid.UUID, sequenceID int64) []byte {
	return fmt.Appendf(nil, "%d/%d/%s/%d", namespace, defaultIndex, id.String(), sequenceID)
}
