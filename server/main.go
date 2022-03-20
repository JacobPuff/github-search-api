package main

import (
	"context"
	"fmt"
	"log"
	"net"

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
	listenOn := "0.0.0.0:8080"
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", listenOn, err)
	}

	server := grpc.NewServer()
	searchv1.RegisterGithubSearchServiceServer(server, &GithubSearchServiceServer{})
	log.Println("Listening on", listenOn)
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

type GithubSearchServiceServer struct {
	searchv1.UnimplementedGithubSearchServiceServer
}

func (s *GithubSearchServiceServer) Search(ctx context.Context, req *searchv1.SearchRequest) (*searchv1.SearchResponse, error) {
	search_term := req.GetSearchTerm()
	user := req.GetUser()
	log.Printf("Got a request to search for '%s' filtered to user '%s'", search_term, user)

	return &searchv1.SearchResponse{Results: []*searchv1.Result{{FileUrl: "yeet", Repo: "haw"}}}, nil
}
