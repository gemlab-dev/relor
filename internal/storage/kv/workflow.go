package kv

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"time"

	pb "github.com/gemlab-dev/relor/gen/pb/db"
	"github.com/gemlab-dev/relor/internal/model"
	"github.com/gemlab-dev/relor/internal/storage"
	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/encoding/protojson"
)

type WorkflowStorage struct {
	db         *bolt.DB
	bucketName string
}

func NewWorkflowStorage(path, bucketName string) (*WorkflowStorage, error) {
	opts := &bolt.Options{Timeout: 1 * time.Second}
	db, err := bolt.Open(path, 0600, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open KV store: %v", err)
	}
	return &WorkflowStorage{db: db, bucketName: bucketName}, nil
}

func (w *WorkflowStorage) Close() error {
	return w.db.Close()
}

func (s *WorkflowStorage) CreateWorkflow(ctx context.Context, w model.Workflow) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(s.bucketName))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}

		return s.saveNewWorkflow(b, w)
	})
}

func (s *WorkflowStorage) GetWorkflow(ctx context.Context, id uuid.UUID) (*model.Workflow, error) {
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

func (s *WorkflowStorage) UpdateNextAction(ctx context.Context, na storage.NextAction) error {
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

		if wf.LastTransitionId != na.LastKnownTransition {
			return fmt.Errorf("current transition ID does not match workflow last transition ID")
		}

		// Validate current action and get next node.
		nextNode, err := wf.Graph.NextNodeID(wf.CurrentNode, na.Label)
		if err != nil {
			return fmt.Errorf("failed to get next node: %w", err)
		}

		newTnPb := pb.WorkflowTransition{
			Id:        wf.LastTransitionId + 1,
			FromNode:  wf.CurrentNode,
			ToNode:    nextNode,
			Label:     na.Label,
			Timestamp: time.Now().Unix(),
		}
		// Write transition
		tnKey, err := transitionKey(na.ID, newTnPb.Id)
		if err != nil {
			return fmt.Errorf("failed to create transition key: %w", err)
		}
		tnData, err := protojson.Marshal(&newTnPb)
		if err != nil {
			return fmt.Errorf("failed to marshal transition proto: %v", err)
		}
		if err := b.Put(tnKey, tnData); err != nil {
			return fmt.Errorf("failed to save transition: %v", err)
		}

		nextLabels, err := wf.Graph.OutLabels(nextNode)
		if err != nil {
			return fmt.Errorf("failed to get out labels: %w", err)
		}
		if len(nextLabels) == 0 {
			wf.Status = model.WorkflowStatusCompleted
		}

		return s.saveWorkflow(b, *wf)
	})
}

func pseudoUUID(id uint64) uuid.UUID {
	paddedBytes := binary.BigEndian.AppendUint64(make([]byte, 8), id)
	return uuid.Must(uuid.FromBytes(paddedBytes))
}

func (s *WorkflowStorage) GetHistory(ctx context.Context, workflowID uuid.UUID) (*model.Transition, error) {
	var transitions []model.RawTransition
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		if b == nil {
			return fmt.Errorf("bucket not found: %s", s.bucketName)
		}

		//TODO: Maybe check if workflowID exists.

		prefix, err := transitionPrefix(workflowID)
		if err != nil {
			return fmt.Errorf("failed to create transition prefix: %w", err)
		}
		c := b.Cursor()
		lastID := uint64(0)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var tpb pb.WorkflowTransition
			err = protojson.Unmarshal(v, &tpb)
			if err != nil {
				return fmt.Errorf("failed to unmarshal transition proto: %v", err)
			}

			//TODO: This is a temp fix for Transition model relying on unsorted UUIDs.
			// The storage layer now offers sorted sequential IDs for transitions.
			// The model will be changed and this code will be removed.
			if lastID == 0 {
				lastID = tpb.Id
			}

			prevUUID := uuid.NullUUID{}
			if tpb.Id > 1 {
				prevUUID = uuid.NullUUID{
					UUID:  pseudoUUID(tpb.Id - 1),
					Valid: true,
				}
			}
			nextUUID := uuid.NullUUID{}
			if tpb.Id < lastID {
				nextUUID = uuid.NullUUID{
					UUID:  pseudoUUID(tpb.Id + 1),
					Valid: true,
				}
			}

			transitions = append(transitions, model.RawTransition{
				ID:      pseudoUUID(tpb.Id),
				WID:     workflowID,
				Prev:    nextUUID,
				Next:    prevUUID,
				From:    tpb.FromNode,
				To:      tpb.ToNode,
				Label:   tpb.Label,
				Created: time.Unix(tpb.Timestamp, 0),
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow history: %w", err)
	}
	return model.NewTransitionHistory(time.Unix(0, 0), transitions)
}

func (s *WorkflowStorage) UpdateTimeout(ctx context.Context, id uuid.UUID, timeout time.Duration) error {
	return fmt.Errorf("not implemented")
}

func (s *WorkflowStorage) GetNextWorkflows(ctx context.Context) ([]model.Workflow, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *WorkflowStorage) loadWorkflow(b *bolt.Bucket, id uuid.UUID) (*model.Workflow, error) {
	key, err := workflowKey(id)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow key: %w", err)
	}
	data := b.Get(key)
	if data == nil {
		return nil, fmt.Errorf("workflow not found")
	}
	var wpb pb.Workflow
	err = protojson.Unmarshal(data, &wpb)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow proto: %v", err)
	}

	currentNode := ""
	lastTransactionID := uint64(0)
	{
		tnPrefix, err := transitionPrefix(id)
		if err != nil {
			return nil, fmt.Errorf("failed to create transition prefix: %w", err)
		}
		tnkey, tnVal := b.Cursor().Seek(tnPrefix)
		if tnkey != nil && bytes.HasPrefix(tnkey, tnPrefix) {
			if tnVal == nil {
				return nil, fmt.Errorf("transition value is missing %s", tnkey)
			}
			var tnpb pb.WorkflowTransition
			err = protojson.Unmarshal(tnVal, &tnpb)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal transition proto: %v", err)
			}
			currentNode = tnpb.ToNode
			lastTransactionID = tnpb.Id
		}
	}

	var w model.Workflow
	if err := w.FromProto(&wpb, currentNode, lastTransactionID); err != nil {
		return nil, fmt.Errorf("failed to convert workflow from proto: %w", err)
	}

	return &w, nil
}

func (s *WorkflowStorage) saveNewWorkflow(b *bolt.Bucket, w model.Workflow) error {
	wpb, err := w.ToProto()
	if err != nil {
		return fmt.Errorf("failed to convert workflow to proto: %w", err)
	}
	data, err := protojson.Marshal(wpb)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow proto to json: %w", err)
	}

	key, err := workflowKey(w.ID)
	if err != nil {
		return fmt.Errorf("failed to create workflow key: %w", err)
	}

	if d := b.Get(key); d != nil {
		return fmt.Errorf("workflow already exists: %s", w.ID)
	}

	return b.Put(key, data)
}

func (s *WorkflowStorage) saveWorkflow(b *bolt.Bucket, w model.Workflow) error {
	wpb, err := w.ToProto()
	if err != nil {
		return fmt.Errorf("failed to convert workflow to proto: %w", err)
	}
	data, err := protojson.Marshal(wpb)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow proto to json: %w", err)
	}

	key, err := workflowKey(w.ID)
	if err != nil {
		return fmt.Errorf("failed to create workflow key: %w", err)
	}

	return b.Put(key, data)
}

func workflowKey(id uuid.UUID) ([]byte, error) {
	return WorkflowKey([]byte{}, id)
}

func transitionPrefix(id uuid.UUID) ([]byte, error) {
	b, err := workflowKey(id)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow key: %w", err)
	}
	return append(b, byte('/')), nil
}

func transitionKey(id uuid.UUID, sequenceID uint64) ([]byte, error) {
	b, err := transitionPrefix(id)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow key: %w", err)
	}
	return TransitionKey(b, sequenceID), nil
}
