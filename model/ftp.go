package model

type FtpHostSetting struct {
	Type       string `json:"type"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	SshKeyPath string `json:"ssh-key-path"`
	UseSshKey  bool   `json:"use-ssh-key"`
	Passive    bool   `json:"passive"`
}

var FtpServerInfo = FtpHostSetting{
	Type:       "ftp",
	Host:       "ftp.example.com",
	Port:       "21",
	Username:   "username",
	Password:   "password",
	SshKeyPath: "",
	UseSshKey:  false,
	Passive:    true,
}
