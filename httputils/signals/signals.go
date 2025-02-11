package signals

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

func ListenAndServeWithSignals(srv *http.Server, sc chan os.Signal, timeout int, signals ...os.Signal) error {

	tov := (timeout * 1000) / int(2)
	// ensure the timeout is always > 0
	if tov <= 0 {
		tov = 500
	}
	timeoutDuration := time.Duration(tov) * time.Millisecond

	signal.Notify(sc, signals...)
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Unexpected error during service shutown: %s", err)
		}
	}()

	log.Printf("Server started on %s.", srv.Addr)
	sig := <-sc
	log.Printf("Received signal %s.", sig)

	log.Printf("Graceful termination window starting: %s", timeoutDuration)
	time.Sleep(timeoutDuration)

	log.Printf("Server shutdown: %s", timeoutDuration)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	return srv.Shutdown(ctx)
}
