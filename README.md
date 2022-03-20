# github-search-api
This was a project for an interview I did. Basically just a limited wrapper of Githubs search api.
## Design decisions
### API
I decided to go with GRPC for the api because it's really nice to work with in go, and I lack experience in the GRPC arena. It's good for keeping an api contract, and transfers nicely into go structs.
### Using Docker
Really simple to do, and easy to use. Makes the server simple to deploy.
## Running
### Docker
For the server with docker you'll want to run 
```
docker build . -t gh-search project
```
followed by
```
docker run -it -p 8080:8080 gh-search-project
```
Afterwards you can kill the program with `Ctrl+C` or equivalent. You know how docker works.
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
5. Run the client. The client has
```
go run client/client.go
```