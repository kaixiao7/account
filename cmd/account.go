package main

import (
	"kaixiao7/account/internal/account"
	"os"
)

func main() {
	command := account.NewAccountServerCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
