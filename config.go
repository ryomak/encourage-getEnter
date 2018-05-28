package main

import "github.com/BurntSushi/toml"

type Config struct {
	Url  string
	Auth string
	Qs   []Query
}

type Query struct {
	Key string
	Val string
}

func GetConfig() Config {
	var config Config
	_, err := toml.DecodeFile(path+"/config.toml", &config)
	if err != nil {
		panic(err)
	}
	return config
}
