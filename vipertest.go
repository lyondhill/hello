package main

import "github.com/spf13/viper"

func main() {
	go setDefaults()
	setDefaults()
}

func setDefaults() {
	for {
		viper.SetDefault("thing", "valule")
	}
}
