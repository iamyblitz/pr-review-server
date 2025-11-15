package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/iamyblitz/pr-reviewer-service/internal/repo"
	"github.com/iamyblitz/pr-reviewer-service/internal/service"
)

func main() {

	r := repo.NewMemoryRepo()
	svc := service.NewService(r)

	_ = svc

	mux := http.NewServeMux()

	mux.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Server alive!"))
	})
	fmt.Println("listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
