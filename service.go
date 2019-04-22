package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kuritka/break-down.io/common/db"
	"github.com/kuritka/break-down.io/services/portal"
	"go.uber.org/dig"
)


func BuildContainer() *dig.Container {

	dbOptions := db.ClientOptions{
		Collection: "calendar",
		Database: "testing",
		Timeout:5,
		ConnectionString:"mongodb://localhost:27017",
		Provider:db.MongoProvider,
	}


	container := dig.New()
	container.Provide(func() db.ClientOptions {return dbOptions})
	container.Provide(mux.NewRouter)
	container.Provide(portal.NewIDP)
	container.Provide(portal.LoadConfig)
	container.Provide(portal.NewServer)
	return container
}


func main() {

	container := BuildContainer()

	err := container.Invoke(func(server *portal.Server) {
		listenAddr := ":8080"

		envPort := os.Getenv("PORT")
		if len(envPort) > 0 {
			listenAddr = ":" + envPort
		}

		log.Printf("attempting listen on %s", listenAddr)
		log.Fatalln(http.ListenAndServe(listenAddr, server))
	})
	if err != nil {
		panic(err)
	}
}