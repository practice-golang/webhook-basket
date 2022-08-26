package model

type Request struct {
	Secret     string         `json:"secret"`      // Secret for authentication
	Repository Repository     `json:"repository"`  // Repository name and clone URL
	Ftp        FtpHostSetting `json:"ftp"`         // ftp, sftp server info
	DeployRoot string         `json:"deploy-root"` // Deployment Root on ftp, sftp server
	DeployName string         `json:"deploy-name"` // Rename the repository on ftp, sftp server
}

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	CloneURL string `json:"clone_url"`
}
