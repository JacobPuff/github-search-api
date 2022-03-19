# github-search-api
This was a project for an interview I did. Basically just a limited wrapper of Githubs search api.
## Design decisions
### API
I decided to go with GRPC for the api because it's really nice to work with in go, and I lack experience in the GRPC arena. It's good for keeping an api contract, and transfers nicely into go structs.