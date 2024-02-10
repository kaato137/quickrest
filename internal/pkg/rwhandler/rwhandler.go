package rwhandler

import (
	"net/http"
	"sync"
)

type RWHandler struct {
	handler http.Handler
	mutex   sync.RWMutex
}

func New(h http.Handler) *RWHandler {
	return &RWHandler{
		handler: h,
	}
}

func (h *RWHandler) SetHandler(nh http.Handler) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.handler = nh
}

func (h *RWHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	h.handler.ServeHTTP(rw, r)
}
