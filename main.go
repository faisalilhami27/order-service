package main

import (
	_ "github.com/spf13/viper/remote"

	"order-service/cmd"
)

func main() {
	cmd.Run()
}
