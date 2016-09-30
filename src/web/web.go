package main

import (
	"config"
	"fmt"
)

type NameStruct struct {
	First string `json:"first"`
	Second string `json:"second"`
}

type DataConfig struct {
	Name NameStruct `json:"name"`
}

func main() {
	dataConfig := &DataConfig{}
	configor := config.NewConfig()
	configor.Load(dataConfig, "/Users/henry-sun/data/www/my-go-config/src/configure/config.json")
}