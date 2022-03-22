# github-search-api
This was a project for an interview I did. Basically just a limited wrapper of Githubs search api.
## Design decisions
### GRPC
I decided to go with GRPC for the api because it's really nice to work with in go, and I lack experience in the GRPC arena. It's good for keeping an api contract, and compiles nicely into go structs.
### Using Docker
Really simple to do, and easy to use. Makes the server simple to deploy.
### Github search endpoint
There are multiple endpoints for specific things in github's search api. I decided to use `/search/code` because it returns files, and in the protos the result has a "file_url".
I think they match up pretty nicely, and I didn't want to overwhelm the response with too many results per endpoint I added. I cut down this endpoint from 100 results to 30 already.
### What is a repo?
For the repo field in the Result, I could have put a few things like `git@github.com/[url stuff]` but I decided to go with the https `.git` file e.g.
`https://github.com/JacobPuff/github-search-api.git`. What I return would depend on the use case, and this covers two for one in.
You can clone the repo, and github will redirect you if you visit in a browser to view the repo.
## Limitations
Limitations are pretty much the same as github's. You have to narrow down the search because this code doesn't use any authentication.
You can find out how to narrow it down here [https://docs.github.com/en/search-github/searching-on-github/searching-code](https://docs.github.com/en/search-github/searching-on-github/searching-code)
## Running
### Docker-compose
Easiest method, run `docker-compose up` in the root directory and it will bring up the server and then the client that'll run a search.
### Docker
For the server with docker you'll want to build it with
```
docker build . -t gh-search-project-server
```
And if you want the client too
```
docker build . -f client.Dockerfile -t gh-search-project-client
```

To run just the server use
```
docker run -it -p 8080:8080 gh-search-project-server
```
Afterwards you can kill the program with `Ctrl+C` or equivalent. You know how docker works.

To run the server _and_ the client you'll need to make a docker network
```
docker network create ghsp-network
```
then run the server
```
docker run -it --rm --name ghsp-server --network ghsp-network gh-search-project-server
```
then you can run the client using this, the client doesn't need a name but it can make looking at the containers running easier
```
docker run -it --rm --name ghsp-client --network ghsp-network gh-search-project-client
```
With this method youll want to clean up the docker network
```
docker network rm ghsp-network
```
### Just golang
You can run this project without docker, it just requires a few more initial steps. This is how I run the client for manual testing for now.
1. Install buf
2. Install go grpc dependencies
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
3. Compile protobufs.
```
buf generate
```
4. Run the server
```
go run server/main.go
```
5. Run the client. The client has two constants at the top of the file if you wanna change their values for the search.
```
go run client/client.go
```