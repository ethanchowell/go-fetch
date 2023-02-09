package main

import (
	"github.com/ethanchowell/artifact-manager/pkg/cmd"
	"log"
)

func main() {
	rootCmd := cmd.New()
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
