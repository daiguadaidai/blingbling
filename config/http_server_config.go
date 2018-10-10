package config

import "fmt"

const (
	HTTP_SERVER_LISTEN_HOST = "0.0.0.0"
	HTTP_SERVER_LISTEN_PORT = 18080
)

type HttpServerConfig struct {
	Host string
	Port int
}

func (this *HttpServerConfig) Address() string {
	return fmt.Sprintf("%v:%v", this.Host, this.Port)
}
