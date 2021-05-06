package main

import (
	"time"

	"gopkg.in/yaml.v3"
)

type Commit struct {
	Repo    string
	Ref     string
	SHA     string
	Message string `yaml:"-"`
	Date    *time.Time
}

func (s *Commit) YAML() (string, error) {
	b, err := yaml.Marshal(s)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
