package model

type Policy struct {
	From        string
	To          string
	HostRewrite string `mapstructure:"host_rewrite"`
}