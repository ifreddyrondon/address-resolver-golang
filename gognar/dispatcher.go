package gognar

import (
	"net/http"
)

type DispatchHandler struct {
	stack []func(next http.HandlerFunc) http.HandlerFunc
	final http.HandlerFunc
}

func (h *DispatchHandler) Attach(middleware ...func(http.HandlerFunc) http.HandlerFunc) {
	h.stack = append(h.stack, middleware...)
}

func (h *DispatchHandler) Finalize(final http.HandlerFunc) {
	h.final = final
	for i := len(h.stack); i > 0; i-- {
		h.final = h.stack[i-1](h.final)
	}
}

func (h DispatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.final(w, r)
}
