package main

import (
	"encoding/csv"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

type Config struct {
	Url       string
	Auth      string
	WriteFile string
	WriteType string
	Qs        []Query
}

type Query struct {
	Key string
	Val string
}

func GetConfig() Config {
	var config Config
	p, _ := os.Executable()
	path := filepath.Dir(p)

	_, err := toml.DecodeFile(path+"/config.toml", &config)
	if err != nil {
		panic(err)
	}
	return config
}

func WriteCsv(name string, users []User) {
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Write([]string{"メンターID", "名前", "ふりがな", "担当メンター", "電話番号", "大学/学部", "インターン", "性別", "判定", "理系"})
	for _, v := range users {
		r := "文系"
		if v.Science {
			r = "理系"
		}
		writer.Write([]string{v.ID, v.Name, v.Yomi, v.Mentor, v.Phone, v.Department, v.Intern, v.Gender, v.Eval, r})
	}
	writer.Flush()
}
