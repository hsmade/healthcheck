// Implements health checks for a service, exposed through http
// First create a `Health`, and then start adding health checks by calling `Add("name")`
// This will return a `HealthCheck` that you can `Update(true)`
// The `Serve(":8080") will start an http server that responds with `200` or `500` on any path.
package healthcheck


import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Health struct {
	health map[string] bool
}

type HealthCheck struct {
	Name   string
	Status bool
	parent *Health
}

func New() Health {
	return Health{
		make(map[string]bool, 1),
	}
}

func (h *Health) Add(name string) *HealthCheck{
	c := HealthCheck{
		name,
		false,
		h,
	}
	h.Update(name, false)
	return &c
}

func (h *Health) Update(child string, status bool) {
	h.health[child] = status
}

func (h *Health) Status() bool {
	for _, status := range h.health {
		if !status {
			return false
		}
	}
	return true
}

func (h *Health) handler(w http.ResponseWriter, r *http.Request) {
	if h.Status() {
		w.WriteHeader(200)
	}
	w.WriteHeader(500)
}

// Start a webserver on `addr` (":8080"). If you specify `delayStop`, the it will set the health to false on sigterm
// and keep serving the health check status `500` for that time. This is to make sure that in flight requests get
// handled while the load balancer finds out we're stopping.
func (h *Health) Serve(addr string, delayStop time.Duration) error {
	if delayStop != 0 {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			for name := range h.health {
				h.Update(name, false)
			}
			time.Sleep(delayStop)
			os.Exit(1)
		}()
	}

	http.HandleFunc("/health", h.handler)
	return http.ListenAndServe(addr, nil)
}

func (h *HealthCheck) Update(status bool) {
	h.parent.Update(h.Name, status)
	h.Status = status
}

