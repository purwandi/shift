package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/purwandi/shift"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func main() {
	_ = godotenv.Load()

	shift.Initialize(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
