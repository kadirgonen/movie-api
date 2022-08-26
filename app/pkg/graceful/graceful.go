package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// ShutdownGin shuts down the given HTTP server gracefully when receiving an os.Interrupt or syscall.SIGTERM signal.
// It will wait for the specified timeout to stop hanging HTTP handlers.
func ShutdownGin(srv *http.Server, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Debug("base.graceful.shutdown")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("base.graceful.shutdown Server forced to shutdown: ", zap.Error(err))
	}
	zap.L().Debug("base.graceful.shutdown server closing")
}
