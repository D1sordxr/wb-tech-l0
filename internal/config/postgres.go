package config

import "fmt"

type Postgres struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Database   string `yaml:"database"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Migrations bool   `yaml:"migrations"`
}

func (p *Postgres) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		p.Host, p.User, p.Password, p.Database, p.Port,
	)
}
