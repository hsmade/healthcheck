// Implements health checks for a service, exposed through http
// First create a `Health`, and then start adding health checks by calling `Add("name")`
// This will return a `HealthCheck` that you can `Update(true)`
// The `Serve(":8080") will start an http server that responds with `200` or `500` on any path.
package healthcheck


import (
	"net/http"
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

func (h *Health) Serve(addr string) error {
	http.HandleFunc("/health", h.handler)
	return http.ListenAndServe(addr, nil)
}

func (h *HealthCheck) Update(status bool) {
	h.parent.Update(h.Name, status)
	h.Status = status
}

