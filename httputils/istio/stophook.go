package istio

import (
	"net/http"
	"time"
)

// StopHook takes a timeout in milliseconds and sends a request to istio's stop hook.
// It is intended to be used as a defer function
func StopHook(timeout int) {
	c := &http.Client{
		Timeout: time.Duration(timeout) * time.Milisecond,
	}
	r, _ := http.NewRequest("POST", "http://localhost:15000/quitquitquit", nil)
	c.Do(r)
}
