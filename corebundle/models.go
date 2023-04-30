package corebundle

type SystemConfiguration struct {
	Database DatabaseConfiguration `json:"database"`
}

type DatabaseConfiguration struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	User     string `json:"user"`
	Password string `json:"pwassword"`
	Name     string `json:"name"`
}
