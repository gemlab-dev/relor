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

		err = b.Put(key, data)
		if err != nil {
			return fmt.Errorf("failed to put workflow: %v", err)
		}

		return nil
	})
}

func (s *WorkflowKVStorage) UpdateNextAction(ctx context.Context, na NextAction) error {
	return nil
}

func (s *WorkflowKVStorage) GetHistory(ctx context.Context, workflowID uuid.UUID) (*model.Transition, error) {
	return nil, nil
}

func (s *WorkflowKVStorage) UpdateTimeout(ctx context.Context, id uuid.UUID, timeout time.Duration) error {
	return nil
}

func (s *WorkflowKVStorage) GetWorkflow(ctx context.Context, id uuid.UUID) (*model.Workflow, error) {
	var w model.Workflow
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		key := workflowKey(defaultNamespace, id)
		data := b.Get(key)
		if data == nil {
			return fmt.Errorf("workflow not found")
		}

		var wpb pb.Workflow
		err := protojson.Unmarshal(data, &wpb)
		if err != nil {
			return fmt.Errorf("failed to unmarshal workflow proto: %v", err)
		}
		if err := w.FromProto(&wpb); err != nil {
			return fmt.Errorf("failed to convert workflow from proto: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (s *WorkflowKVStorage) GetNextWorkflows(ctx context.Context) ([]model.Workflow, error) {
	var workflows []model.Workflow
	return workflows, nil
}

func workflowKey(namespace int, id uuid.UUID) []byte {
	return fmt.Appendf(nil, "%d/%d/%s", namespace, defaultIndex, id.String())
}
