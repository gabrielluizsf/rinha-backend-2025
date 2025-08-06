package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gabrielluizsf/rinha-backend-2005/db"
	"github.com/gabrielluizsf/rinha-backend-2005/routes"
	"github.com/gabrielluizsf/rinha-backend-2005/worker"
)

const PORT = "8080"

func main() {
	db.InitRedis()

	routerManager := http.NewServeMux()
	routes.InitRoutes(routerManager)

	worker.StartLeaderElection()
	worker.StartWorker()

	fmt.Printf("Server running on :%s\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), routerManager))
}
