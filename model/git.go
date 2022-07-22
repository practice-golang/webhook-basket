package model

import "github.com/go-git/go-git/v5/plumbing/transport/http"

var AuthInfo = &http.BasicAuth{
	Username: "username",
	Password: "password",
}
