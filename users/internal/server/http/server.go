package http

import (
	"context"
	"ecom-users/internal/application"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer struct {
	srv http.Server
}

func New(app *application.Application) HttpServer {
	return HttpServer{
		srv: http.Server{
			Addr: fmt.Sprintf(":%d", app.Config.PORT),
			ReadTimeout: 10* time.Second,
			WriteTimeout: 30* time.Second,
			IdleTimeout: time.Minute,
		},
	}
}

func (s *HttpServer) Routes(routes http.Handler) {
	s.srv.Handler = routes
}

func (s *HttpServer) Serve(app *application.Application) error {
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <- quit

		app.Logger.PrintInfo("caught signal", map[string]string{
			"signal": sig.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.WG.Wait()
		shutdownError <- nil
	}()

	err := s.srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	app.Logger.PrintInfo("starting server", map[string]string{
		"addr": s.srv.Addr,
	})

	err = <- shutdownError
	if err != nil {
		return err
	}

	app.Logger.PrintInfo("stopped server", map[string]string{
		"addr": fmt.Sprintf("%s", app.Config.PORT),
	})

	return nil
}