package main

import (
	_ "github.com/go-playground/validator/v10"
	_ "github.com/spf13/viper/remote"

	"order-service/cmd"
)

func main() {
	cmd.Run()
}
