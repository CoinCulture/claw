package main

import (
	"github.com/russross/blackfriday"
	"github.com/spf13/viper"
)

// append to list if not in list already
func appendNew(list []string, name string) []string {
	found := false
	for _, el := range list {
		if el == name {
			found = true
		}
	}
	if !found {
		list = append(list, name)
	}
	return list
}

// load params config

func loadConfig() (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigName("params")
	config.SetConfigType("toml")
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func markdown2html(b []byte) []byte {
	return blackfriday.MarkdownBasic(b)
}
