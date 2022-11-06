package main

import (
	"fmt"
	"os"

	"github.com/ryuichi1208/go-check-longtransaction-cnt/cmd"
)

func init() {
	if os.Getenv("MYSQL_PASSWORD") == "" {
		fmt.Println("Environment variable MYSQL_PASSWORD not set")
		os.Exit(1)
	}
}

func main() {
	cmd.Run()
}
