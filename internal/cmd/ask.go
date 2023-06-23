package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func askUsername() string {
	fmt.Print("Enter username: ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	if err != nil {
		panic(fmt.Errorf("an error occured while reading username. Error: %s", err))
	}
	username = strings.TrimSuffix(username, "\n")

	return username
}

func askHost() string {
	fmt.Print("Enter grafana host: ")
	reader := bufio.NewReader(os.Stdin)
	host, err := reader.ReadString('\n')
	if err != nil {
		panic(fmt.Errorf("an error occured while reading host. Error: %s", err))
	}
	host = strings.TrimSuffix(host, "\n")
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")

	return host
}

func askPassword() string {
	fmt.Print("Enter password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))

	if err != nil {
		panic(fmt.Errorf("an error occured while reading password. Error: %s", err))
	}
	fmt.Println()

	return string(password)
}
