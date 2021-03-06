package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	graphqldemo "github.com/sunfmin/go-gqlgen-graphql-demo"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	g := handler.GraphQL(graphqldemo.NewExecutableSchema(graphqldemo.Config{Resolvers: &graphqldemo.Resolver{}}))

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", graphqldemo.UserLoaderMiddleware(g))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
