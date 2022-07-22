package model

type Request struct {
	Repository  Repository `json:"repository"`
	Pusher      Pusher     `json:"pusher"`
	Destination string     `json:"destination"`
}

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	CloneURL string `json:"clone_url"`
}

type Pusher struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
