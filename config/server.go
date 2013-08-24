package config

var ServerConfig Server

type Server struct {
	ListenAddress string
	NotifyEmail   string
	FromEmail     string
	SMTPAddress   string
}
