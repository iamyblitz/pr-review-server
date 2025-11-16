package httpapi

import (
	"net/http"

	"github.com/iamyblitz/pr-reviewer-service/internal/service"
)

func NewRouter(svc *service.Service) http.Handler {
	h := NewHandler(svc)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/team/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.AddTeam(w, r)
	})

	mux.HandleFunc("/team/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.GetTeam(w, r)
	})

	mux.HandleFunc("/users/setIsActive", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.SetUserIsActive(w, r)
	})

	mux.HandleFunc("/pullRequest/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.CreatePullRequest(w, r)
	})

	mux.HandleFunc("/pullRequest/merge", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.MergePullRequest(w, r)
	})

	mux.HandleFunc("/users/getReview", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.GetUserReviews(w, r)
	})

	mux.HandleFunc("/pullRequest/reassign", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.ReassignReviewer(w, r)
	})

	return mux
}
