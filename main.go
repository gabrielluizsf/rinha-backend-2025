package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gabrielluizsf/rinha-backend-2005/db"
	"github.com/gabrielluizsf/rinha-backend-2005/routes"
	"github.com/gabrielluizsf/rinha-backend-2005/worker"
)

func main() {
	db.InitRedis()

	routerManager := http.NewServeMux()
	routes.InitRoutes(routerManager)

	worker.StartLeaderElection()
	worker.StartWorker()

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", routerManager))
}
