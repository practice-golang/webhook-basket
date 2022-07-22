package model

type FtpHostSetting struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var FtpServerInfo = FtpHostSetting{
	Type:     "ftp",
	Host:     "ftp.example.com",
	Port:     "21",
	Username: "username",
	Password: "password",
}
