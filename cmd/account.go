package main

import (
	"os"

	"kaixiao7/account/internal/account"
)

func main() {
	command := account.NewAccountServerCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
