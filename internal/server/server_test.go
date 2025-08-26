package server

import (
	"context"
	"net/http"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/gemlab-dev/relor/internal/job"
	"github.com/gemlab-dev/relor/internal/workflow"
)

type loggerMock struct {
	lastInfoMsg, lastErrMsg string
	mu                      sync.Mutex
}

func (l *loggerMock) InfoContext(ctx context.Context, msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lastInfoMsg = msg
}

func (l *loggerMock) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lastErrMsg = msg
}

type mockGraphvizHandler struct{}

func (m *mockGraphvizHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestShutdownWithCtxCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	logger := &loggerMock{}
	srv := New(8080, logger, &workflow.Server{}, &job.Server{}, &mockGraphvizHandler{})
	go srv.Serve(ctx)

	cancel()

	time.Sleep(10 * time.Millisecond)

	logger.mu.Lock()
	defer logger.mu.Unlock()

	want := "Stopping server"
	if logger.lastInfoMsg != want {
		t.Errorf("Unexpected log message: %q; want %q", logger.lastInfoMsg, want)
	}

	if logger.lastErrMsg != "" {
		t.Errorf("Unexpected log message: %q; want empty message", logger.lastErrMsg)
	}
}

func TestShutdownWithSignal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := &loggerMock{}
	srv := New(8080, logger, &workflow.Server{}, &job.Server{}, &mockGraphvizHandler{})
	srv.notify = func(c chan<- os.Signal, sig ...os.Signal) {
		c <- syscall.SIGTERM
	}

	go srv.Serve(ctx)

	time.Sleep(10 * time.Millisecond)

	logger.mu.Lock()
	defer logger.mu.Unlock()

	want := "Stopping server"
	if logger.lastInfoMsg != want {
		t.Errorf("Unexpected log message: %q; want %q", logger.lastInfoMsg, want)
	}

	if logger.lastErrMsg != "" {
		t.Errorf("Unexpected log message: %q; want empty message", logger.lastErrMsg)
	}
}

func TestShutdownWithError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := &loggerMock{}
	srv := New(-8080, logger, &workflow.Server{}, &job.Server{}, &mockGraphvizHandler{}) // invalid port

	go srv.Serve(ctx)

	time.Sleep(10 * time.Millisecond)

	logger.mu.Lock()
	defer logger.mu.Unlock()

	want := "Starting server"
	if logger.lastInfoMsg != want {
		t.Errorf("Unexpected log message: %q; want %q", logger.lastInfoMsg, want)
	}

	if logger.lastErrMsg == "" {
		t.Error("Expected error message", logger.lastErrMsg)
	}
}
