package utils

import (
	"context"
	"main/internal/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// GracefullShutdown is waiting for SIGINT or SIGTERM signals from http.Server to perform shutdown and log this action.
func GracefullShutdown(server *http.Server) {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-termChan

	if err := server.Shutdown(context.Background()); err != nil {
		log.Logger.Fatalln("server.Shutdown:", err)
	}

	log.Logger.Infoln("Gracefully shutdown after signal:", sig)

}
