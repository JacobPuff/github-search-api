package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.
	searchv1 "github.com/JacobPuff/github-search-api/gen/proto/go/githubsearch/v1"
	"google.golang.org/grpc"
)

const RESULTS_PER_PAGE = "30"

type GithubResults struct {
	TotalCount        int  `json:"total_count"`
	IncompleteResults bool `json:"incomplete_results"`
	Items             []struct {
		FileBlob       string `json:"html_url"`
		RepositoryData struct {
			URL string `json:"html_url"`
		} `json:"repository"`
	} `json:"items"`
}

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
	log.Printf("Got a request to search for '%t' filtered to user '%t'", search_term, user)
	githubRequest, err := http.NewRequest("GET", "https://api.github.com/search/code", nil)
	if err != nil {
		log.Printf("An error occured when creating request")
	}

	githubRequest.Header.Add("Accept", "application/vnd.github.v3+json")
	queryParams := githubRequest.URL.Query()
	if user != "" {
		log.Println("USER SET")
		search_term = fmt.Sprintf("%s user:%s", search_term, user)
	}

	queryParams.Add("q", search_term)
	queryParams.Add("per_page", RESULTS_PER_PAGE)
	githubRequest.URL.RawQuery = queryParams.Encode()

	response, err := http.DefaultClient.Do(githubRequest)
	if err != nil {
		log.Printf("An error occured when sending request to github api %s", err)
	}
	defer response.Body.Close()
	log.Printf("URL: %s", githubRequest.URL)
	log.Printf("STATUS: %s", response.Status)

	githubResults := &GithubResults{}
	err = json.NewDecoder(response.Body).Decode(githubResults)
	if err != nil {
		log.Printf("An error occured while decoding json %s", err)
	}
	grpcResults := []*searchv1.Result{}
	for _, item := range githubResults.Items {
		result := &searchv1.Result{
			FileUrl: item.FileBlob,
			Repo:    item.RepositoryData.URL + ".git",
		}
		grpcResults = append(grpcResults, result)
	}
	log.Println("\n")
	return &searchv1.SearchResponse{Results: grpcResults}, nil
}
