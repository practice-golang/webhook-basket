package model

type Request struct {
	Repository  Repository     `json:"repository"`  // Repository name and clone URL
	Pusher      Pusher         `json:"pusher"`      // Pusher is the user who pushed the commit
	Ftp         FtpHostSetting `json:"ftp"`         // ftp, sftp server info
	Destination string         `json:"destination"` // Deployment Root on ftp, sftp server
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
