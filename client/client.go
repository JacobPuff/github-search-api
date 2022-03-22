package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.
	searchv1 "github.com/JacobPuff/github-search-api/gen/proto/go/githubsearch/v1"
	"google.golang.org/grpc"
)

func setupLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	setupLogging()
	var SearchTerm = GetEnvOrDefault("SEARCH_TERM", "search repo:jacobpuff/github-search-api")
	var User = GetEnvOrDefault("USER", "")
	connectTo := GetEnvOrDefault("GH_SEARCH_SERVER_ADDRESS", "0.0.0.0")
	connectTo += ":" + GetEnvOrDefault("GH_SEARCH_SERVER_PORT", "9090")

	log.Infof("Attempting connection to %s", connectTo)
	conn, err := grpc.Dial(connectTo, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to GithubSearchService on %s: %s", connectTo, err)
	}
	log.Infof("Connected to %s", connectTo)

	github_search := searchv1.NewGithubSearchServiceClient(conn)
	response, err := github_search.Search(context.Background(), &searchv1.SearchRequest{
		SearchTerm: SearchTerm,
		User:       User,
	})
	if err != nil {
		log.Fatalf("failed to search: %s", err)
	}

	log.Infof("Successfully searched term '%s' filtered to user '%s', got response:", SearchTerm, User)
	for _, item := range response.Results {
		log.Infof("file_url: %s\nrepo: %s\n\n", item.FileUrl, item.Repo)
	}
}

func GetEnvOrDefault(env, defaultValue string) string {
	value := os.Getenv(env)
	if value != "" {
		return value
	}
	return defaultValue
}
