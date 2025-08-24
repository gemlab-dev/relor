package graphviz

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-graphviz"
	"github.com/google/uuid"

	"github.com/gemlab-dev/relor/internal/model"
)

const prefix = "/"

type Logger interface {
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type WorkflowStorage interface {
	GetWorkflow(ctx context.Context, id uuid.UUID) (*model.Workflow, error)
	GetHistory(ctx context.Context, id uuid.UUID) (*model.Transition, error)
}

type Handler struct {
	wfStore WorkflowStorage
	logger  Logger
}

func NewHandler(logger Logger, wfStore WorkflowStorage) *Handler {
	return &Handler{
		wfStore: wfStore,
		logger:  logger,
	}
}

// ServeHTTP handles GET requests to generate and serve workflow graphs as SVG
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path (e.g., /{id})
	path := r.URL.Path
	if len(path) <= len(prefix) {
		http.Error(w, "Workflow ID is required", http.StatusBadRequest)
		return
	}

	idStr := path[len(prefix):]

	// Parse workflow ID
	wid, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Failed to parse workflow ID", "err", err, "id", idStr)
		http.Error(w, "Invalid workflow ID format", http.StatusBadRequest)
		return
	}

	// Get workflow from storage
	workflow, err := h.wfStore.GetWorkflow(r.Context(), wid)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Failed to get workflow", "err", err, "id", wid)
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	if workflow == nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	// Get transition history
	th, err := h.wfStore.GetHistory(r.Context(), wid)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Failed to get transition history", "err", err, "id", wid)
		http.Error(w, "Failed to get workflow history", http.StatusInternalServerError)
		return
	}

	// Generate Graphviz DOT representation
	dotContent, err := Dot(*workflow, th)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Failed to generate graphviz", "err", err, "id", wid)
		http.Error(w, "Failed to generate graph", http.StatusInternalServerError)
		return
	}

	// Convert DOT to SVG
	svg, err := h.dotToSVG(r.Context(), dotContent)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Failed to convert DOT to SVG", "err", err)
		http.Error(w, "Failed to generate SVG", http.StatusInternalServerError)
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache")

	// Write SVG response
	_, err = w.Write([]byte(svg))
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Failed to write SVG response", "err", err)
	}
}

// dotToSVG converts DOT format to SVG using goccy/go-graphviz library
func (h *Handler) dotToSVG(ctx context.Context, dotContent string) (string, error) {
	g, err := graphviz.New(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create graphviz instance: %w", err)
	}
	defer g.Close()

	graph, err := graphviz.ParseBytes([]byte(dotContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse DOT: %w", err)
	}
	defer graph.Close()

	var buf bytes.Buffer
	if err := g.Render(ctx, graph, graphviz.SVG, &buf); err != nil {
		return "", fmt.Errorf("failed to render SVG: %w", err)
	}

	return buf.String(), nil
}
