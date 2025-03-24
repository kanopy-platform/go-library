package signals_test

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/kanopy-platform/go-library/httputils/signals"
	"github.com/stretchr/testify/assert"
)

func TestListenAndServeWithSignal(t *testing.T) {
	handler := http.NewServeMux()
	// the delay and the timeout value are being aligned here
	// to have requests that meet the following:
	// 1. They take longer than timeout/2 to finish
	dh := &delayHandler{
		delay: 4,
	}
	timeout := 6
	handler.HandleFunc("/", dh.HandleWithDelay)
	handler.Handle("/404", http.NotFoundHandler())

	srv := &http.Server{
		Addr:    "localhost:11111",
		Handler: handler,
	}
	sc := make(chan os.Signal, 1)
	ec := make(chan error, 1)

	go func() {
		ec <- signals.ListenAndServeWithSignals(srv, sc, timeout, syscall.SIGINT)
	}()

	//ensure the server is running
	time.Sleep(1500 * time.Millisecond)
	_, er := http.Get("http://localhost:11111/404")
	assert.NoError(t, er)

	// queue up requests
	var wg sync.WaitGroup
	var success, failure atomic.Uint32
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, e := http.Get("http://localhost:11111/")

			if e != nil {
				t.Log(e)
				failure.Add(1)
			} else {
				success.Add(1)
			}

		}()
	}

	// let request queue
	time.Sleep(1000 * time.Millisecond)

	sc <- syscall.SIGINT
	dc := make(chan struct{})
	go func() {
		wg.Wait()
		close(dc)
	}()

	var err error
outer:
	for {
		select {
		case err = <-ec:
			break outer
		case <-dc:
			break outer
		default:
			time.Sleep(3 * time.Second)
		}

	}

	assert.NoError(t, err)
	assert.Equal(t, 10, int(success.Load()))
	assert.Equal(t, 0, int(failure.Load()))

}

type delayHandler struct {
	delay int
}

func (d *delayHandler) HandleWithDelay(w http.ResponseWriter, r *http.Request) {
	if d.delay == 0 {
		d.delay = 3
	}
	time.Sleep(time.Duration(d.delay) * time.Second)
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, "OKAY")
}
