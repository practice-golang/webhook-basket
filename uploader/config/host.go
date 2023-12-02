package config

type Host struct {
	Type       string
	Hostname   string
	Port       string
	Username   string
	Password   string
	SshKeyPath string
	SshKeyData string
	SrcBase    string
	DstBase    string
	SrcName    string // Source repository name
	DstName    string // Rename the repository on ftp, sftp server
	UseSshKey  bool
	Passive    bool
	// QueSheets []QueSheet
}
