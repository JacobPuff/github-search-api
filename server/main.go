package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"

	appconfig "github.com/JacobPuff/github-search-api/appconfig"
	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.
	searchv1 "github.com/JacobPuff/github-search-api/gen/proto/go/githubsearch/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	appconfig.SetupLogging()
	listenOn := appconfig.SearchServerAddress + ":" + appconfig.SearchServerPort
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		log.Fatalf("failed to listen on %s: %w", listenOn, err)
	}

	server := grpc.NewServer()
	searchv1.RegisterGithubSearchServiceServer(server, &GithubSearchServiceServer{})
	log.Infof("Listening on %s", listenOn)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %w", err)
	}
}

type GithubSearchServiceServer struct {
	searchv1.UnimplementedGithubSearchServiceServer
}

func (s *GithubSearchServiceServer) Search(ctx context.Context, req *searchv1.SearchRequest) (*searchv1.SearchResponse, error) {
	search_term := req.GetSearchTerm()
	user := req.GetUser()
	if search_term == "" && user == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Need at least one of search_term or user")
	}
	log.Infof("Got a request to search for '%s' filtered to user '%s'", search_term, user)
	githubRequest, err := http.NewRequest("GET", "https://api.github.com/search/code", nil)
	if err != nil {
		log.Error("An error occured when creating request")
		return nil, status.Errorf(codes.Internal, "An error occured while processing your request")
	}

	githubRequest.Header.Add("Accept", "application/vnd.github.v3+json")
	queryParams := githubRequest.URL.Query()
	if user != "" {
		log.Debug("USER SET")
		search_term = fmt.Sprintf("%s user:%s", search_term, user)
	}

	queryParams.Add("q", search_term)
	queryParams.Add("per_page", appconfig.ResultsPerPage)
	githubRequest.URL.RawQuery = queryParams.Encode()

	response, err := http.DefaultClient.Do(githubRequest)
	if err != nil {
		log.Errorf("An error occured when sending request to github api %s", err)
		return nil, status.Errorf(codes.Internal, "An error occured while processing your request")
	}
	defer response.Body.Close()
	log.Debugf("URL: %s", githubRequest.URL)
	log.Debugf("STATUS: %s", response.Status)

	if response.StatusCode != http.StatusOK {
		log.Errorf("Received bad request status: %d", response.StatusCode)
		if response.StatusCode == http.StatusForbidden {
			return nil, status.Errorf(codes.Unauthenticated, "You aren't able to access this resource")
		}
		if response.StatusCode == http.StatusUnprocessableEntity {
			return nil, status.Errorf(codes.InvalidArgument, "Search term is invalid, scope needs to be narrowed down")
		}
		if response.StatusCode == http.StatusServiceUnavailable {
			return nil, status.Errorf(codes.Unavailable, "Service Unavailable")
		}
	}

	githubResults := &GithubResults{}
	err = json.NewDecoder(response.Body).Decode(githubResults)
	if err != nil {
		log.Errorf("An error occured while decoding json %s", err)
		return nil, status.Errorf(codes.Internal, "An error occured while processing your request")
	}
	grpcResults := []*searchv1.Result{}
	for _, item := range githubResults.Items {
		result := &searchv1.Result{
			FileUrl: item.FileBlob,
			Repo:    item.RepositoryData.URL + ".git",
		}
		grpcResults = append(grpcResults, result)
	}
	// This makes the output easier to read
	fmt.Println("\n")

	return &searchv1.SearchResponse{Results: grpcResults}, nil
}
