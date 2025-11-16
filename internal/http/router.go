package httpapi

import (
	"net/http"

	"github.com/iamyblitz/pr-reviewer-service/internal/service"
)

func NewRouter(svc *service.Service) http.Handler {
	h := NewHandler(svc)

	mux := http.NewServeMux()

	// health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// team/add
	mux.HandleFunc("/team/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.AddTeam(w, r)
	})

	return mux
}
