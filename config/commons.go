package config

type ServerPreference struct {
	AppName     string
	PrepareRoot string // Save path to cloned repos
	ListenAddr  string
	ListenPort  string
}

// Repository - a repo data
type Repository struct {
	Name       string
	CloneURI   string
	SSHos      string
	SSHuri     string
	SSHid      string
	SSHpasswd  string
	DeployRoot string
	RepoRoot   string
}

// Repos - Repositories
var Repos = map[string]Repository{}
var ServerInfo = ServerPreference{}
