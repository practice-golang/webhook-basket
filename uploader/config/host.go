package config

type Host struct {
	Type     string
	Hostname string
	Port     string
	Username string
	Password string
	SrcBase  string
	DstBase  string
	// QueSheets []QueSheet
}
