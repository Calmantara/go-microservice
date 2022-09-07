package model

type CommonHost struct {
	Host       string `json:"host"`
	MaxRetries int    `json:"maxretries"`
	Timeout    int    `json:"timeout"`
	TimeToWait int    `json:"timetowait"`
}

type GoWallet struct {
	CommonHost
}

type GoEmitter struct {
	CommonHost
}

type Service struct {
	Name     string `json:"name"`
	GrpcPort string `json:"grpcport"`
}
