package main

import (
	"fmt"
	"log"
	"net/http"

	httpapi "github.com/iamyblitz/pr-reviewer-service/internal/http"
	"github.com/iamyblitz/pr-reviewer-service/internal/repo"
	"github.com/iamyblitz/pr-reviewer-service/internal/service"
)

func main() {

	r := repo.NewMemoryRepo()
	svc := service.NewService(r)

	router := httpapi.NewRouter(svc)

	fmt.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}

}
