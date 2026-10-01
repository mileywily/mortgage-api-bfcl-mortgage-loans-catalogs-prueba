package server_fif

import (
	"context"
	"errors"
	loggerFif "github.com/falabella-regulado/go-lib-logger-fif"
	"github.com/falabella-regulado/go-lib-server-fif/configuration"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type ServerOptions struct {
	Logger  loggerFif.Logger
	Handler http.Handler
}

func StartServer(options ServerOptions) error {

	config := configuration.GetConfiguration()
	port := config.Port

	if len(port) > 0 && port[0] != ':' {
		port = ":" + port
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      options.Handler,
		WriteTimeout: config.WriteTolerance,
		ReadTimeout:  config.ReadTolerance,
		IdleTimeout:  config.IdleTolerance,
	}

	go func() error {
		// Servicio en ejecución
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			options.Logger.Error("listen: %s\n", "cause: ", err)
			return err
		}
		return nil
	}()

	// Manejo de señales de interrupción para apagar el servidor correctamente
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	options.Logger.Info("Stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), config.WaitTolerance)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			options.Logger.Info("executing forced shutdown")
			if err := srv.Close(); err != nil {
				options.Logger.Error("server close failed error")
				return err
			}
			options.Logger.Info("forced shutdown completed")
			return nil
		}
	}
	options.Logger.Info("server close completed")

	return nil
}
