package main

import (
	"github.com/purwandi/shift"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{}
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "database migration utility",
	}

	shift.SetConfig(shift.Config{
		Driver:   "postgres",
		Username: "user",
		Password: "password",
		Database: "database",
		Port:     5432,
		Options:  "sslmode=disable",
	})
	shift.Initialize(migrateCmd)

	rootCmd.AddCommand(migrateCmd)
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
