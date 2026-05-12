package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
	core_logger "github.com/pkpal-uhobp/fulfillment-app/internal/core/logger"
	core_http_middleware "github.com/pkpal-uhobp/fulfillment-app/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux         *http.ServeMux
	config      Config
	log         *core_logger.Logger
	middlewares []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middlewares ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:         http.NewServeMux(),
		config:      config,
		log:         log,
		middlewares: middlewares,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.APIVersion())

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	handler := http.Handler(s.mux)

	if len(s.middlewares) > 0 {
		handler = core_http_middleware.ChainMiddlewares(
			handler,
			s.middlewares...,
		)
	}

	server := &http.Server{
		Addr:         s.config.Addr,
		Handler:      handler,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	errCh := make(chan error, 1)

	go func() {
		s.log.Info(
			"starting HTTP server",
			zap.String("addr", s.config.Addr),
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf(
				"%w: listen and serve HTTP: %v",
				core_errors.ErrInternal,
				err,
			)
			return
		}

		errCh <- nil
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}

		return nil

	case <-ctx.Done():
		s.log.Info("shutting down HTTP server")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf(
				"%w: shutdown HTTP server: %v",
				core_errors.ErrInternal,
				err,
			)
		}

		s.log.Info("HTTP server stopped")

		return nil
	}
}
