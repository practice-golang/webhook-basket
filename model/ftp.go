package model

type FtpHostSetting struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var FtpServerInfo = FtpHostSetting{
	Host:     "ftp.example.com",
	Port:     "21",
	Username: "username",
	Password: "password",
}
