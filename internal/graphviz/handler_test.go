package graphviz

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	gpb "github.com/gemlab-dev/relor/gen/pb/graph"
	"github.com/gemlab-dev/relor/internal/model"
	"github.com/google/uuid"
)

// mockLogger implements the Logger interface for testing purposes.
type mockLogger struct{}

func (m *mockLogger) InfoContext(ctx context.Context, msg string, args ...any)  {}
func (m *mockLogger) ErrorContext(ctx context.Context, msg string, args ...any) {}

// mockWorkflowStorage implements the WorkflowStorage interface for testing.
type mockWorkflowStorage struct {
	workflow    *model.Workflow
	history     *model.Transition
	returnError error
}

func (m *mockWorkflowStorage) GetWorkflow(ctx context.Context, id uuid.UUID) (*model.Workflow, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	if m.workflow != nil && m.workflow.ID == id {
		return m.workflow, nil
	}
	return nil, nil // Simulate not found
}

func (m *mockWorkflowStorage) GetHistory(ctx context.Context, id uuid.UUID) (*model.Transition, error) {
	if m.returnError != nil {
		return nil, m.returnError
	}
	if m.history != nil {
		return m.history, nil
	}
	return model.NewTransitionHistory(time.Time{}, nil) // Return empty history
}

func TestGraphvizHandler(t *testing.T) {
	workflowID := uuid.New()

	// A simple workflow and history for successful tests.
	pbGraph := &gpb.Graph{
		Start: "start",
		Nodes: []*gpb.Node{
			{Id: "start"},
			{Id: "end"},
		},
		Edges: []*gpb.Edge{
			{
				FromId: "start",
				ToId:   "end",
				Condition: &gpb.TransitionCondition{
					OperationResult: "next",
				},
			},
		},
	}
	modelGraph := &model.Graph{}
	if err := modelGraph.FromProto(pbGraph); err != nil {
		t.Fatalf("failed to create graph from proto: %v", err)
	}
	successWorkflow := model.NewWorkflow(workflowID, modelGraph, time.Now())

	successHistory, _ := model.NewTransitionHistory(time.Now(), []model.RawTransition{
		{From: "start", To: "end", Label: "next"},
	})

	testCases := []struct {
		name                  string
		method                string
		url                   string
		storage               *mockWorkflowStorage
		expectedStatus        int
		expectedBodyToContain string
	}{
		{
			name:                  "Success",
			method:                http.MethodGet,
			url:                   "/" + workflowID.String(),
			storage:               &mockWorkflowStorage{workflow: successWorkflow, history: successHistory},
			expectedStatus:        http.StatusOK,
			expectedBodyToContain: "<svg",
		},
		{
			name:           "Method Not Allowed",
			method:         http.MethodPost,
			url:            "/" + workflowID.String(),
			storage:        &mockWorkflowStorage{},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:                  "Invalid UUID",
			method:                http.MethodGet,
			url:                   "/not-a-uuid",
			storage:               &mockWorkflowStorage{},
			expectedStatus:        http.StatusBadRequest,
			expectedBodyToContain: "Invalid workflow ID format",
		},
		{
			name:           "Workflow Not Found",
			method:         http.MethodGet,
			url:            "/" + uuid.NewString(), // A different, non-existent UUID
			storage:        &mockWorkflowStorage{workflow: successWorkflow},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Storage Error on GetWorkflow",
			method:         http.MethodGet,
			url:            "/" + workflowID.String(),
			storage:        &mockWorkflowStorage{returnError: errors.New("db is down")},
			expectedStatus: http.StatusNotFound, // Not found is returned for any GetWorkflow error
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := httptest.NewRequest(tc.method, tc.url, nil)
			rr := httptest.NewRecorder()
			handler := NewHandler(&mockLogger{}, tc.storage)

			// Act
			handler.ServeHTTP(rr, req)

			// Assert
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			if tc.expectedBodyToContain != "" {
				if !strings.Contains(rr.Body.String(), tc.expectedBodyToContain) {
					t.Errorf("handler returned unexpected body: got %v want to contain %q", rr.Body.String(), tc.expectedBodyToContain)
				}
			}
		})
	}
}
