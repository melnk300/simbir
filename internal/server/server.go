package server

import (
	"net/http"
)

type Config struct {
	Host string `json:"host"`
	Port string `json:"post"`
}

func InitServer(cfg Config) error {
	r := initRoute()
	return http.ListenAndServe(
		cfg.Host+":"+cfg.Port,
		r,
	)
}
