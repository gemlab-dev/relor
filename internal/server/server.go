package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/gemlab-dev/relor/gen/pb/api"
)

type Logger interface {
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type GraphvizHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Notify func(c chan<- os.Signal, sig ...os.Signal)

type Server struct {
	logger  Logger
	port    int
	notify  Notify
	httpMux *http.ServeMux
	wfs     pb.WorkflowServiceServer
	js      pb.JobServiceServer
}

func New(port int, logger Logger, wfs pb.WorkflowServiceServer, js pb.JobServiceServer, gh GraphvizHandler) *Server {
	mux := http.NewServeMux()
	mux.Handle("/graph/", http.StripPrefix("/graph", gh))

	return &Server{
		logger:  logger,
		port:    port,
		notify:  signal.Notify,
		httpMux: mux,
		wfs:     wfs,
		js:      js,
	}
}

func (s Server) Serve(ctx context.Context) error {
	s.logger.InfoContext(ctx, "Starting server", "port", s.port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to listen", "err", err)
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer lis.Close()

	gs := grpc.NewServer()
	// TODO: Add graceful shutdown.
	// defer gs.GracefulStop()
	defer gs.Stop()

	pb.RegisterWorkflowServiceServer(gs, s.wfs)
	pb.RegisterJobServiceServer(gs, s.js)
	reflection.Register(gs)

	stopChan := make(chan os.Signal, 1)
	s.notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	errChan := make(chan error)
	go func() {
		if err := gs.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	s.logger.InfoContext(ctx, "Starting HTTP server", "port", 8080)
	httpSrv := http.Server{
		Addr:    ":8080",
		Handler: s.httpMux,
	}
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("HTTP server failed: %w", err)
		}
	}()
	defer httpSrv.Shutdown(ctx)

	select {
	case err := <-errChan:
		s.logger.ErrorContext(ctx, "Error serving", "err", err)
		return err
	case <-stopChan:
		s.logger.InfoContext(ctx, "Received stop signal")
	case <-ctx.Done():
	}
	s.logger.InfoContext(ctx, "Stopping server")
	return nil
}
