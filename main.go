package main

import (
	"fmt"

	"github.com/gabrielluizsf/rinha-backend-2005/adapter"
	"github.com/gabrielluizsf/rinha-backend-2005/db"
	"github.com/gabrielluizsf/rinha-backend-2005/routes"
	"github.com/gabrielluizsf/rinha-backend-2005/worker"
)

const PORT = "8080"

func main() {
	db.InitRedis()

	routerManager := adapter.Server()
	routes.InitRoutes(routerManager)

	worker.StartLeaderElection()
	worker.StartWorker()

	routerManager.Listen(fmt.Sprintf(":%s", PORT))
}
