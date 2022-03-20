package main

import (
	"context"
	"fmt"
	"log"

	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.
	searchv1 "github.com/JacobPuff/github-search-api/gen/proto/go/githubsearch/v1"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func run() error {
	connectTo := "127.0.0.1:8080"
	conn, err := grpc.Dial(connectTo, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to GithubSearchService on %s: %w", connectTo, err)
	}
	log.Println("Connected to", connectTo)

	github_search := searchv1.NewGithubSearchServiceClient(conn)
	search_term := "testin"
	user := "stuff"
	response, err := github_search.Search(context.Background(), &searchv1.SearchRequest{
		SearchTerm: search_term,
		User:       user,
	})
	if err != nil {
		return fmt.Errorf("failed to search: %w", err)
	}

	log.Printf("Successfully searched term '%s' filtered to user '%s', got response %v", search_term, user, response)
	return nil
}
