package main

import (
	"github.com/joho/godotenv"
	"github.com/purwandi/shift"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func main() {
	_ = godotenv.Load()

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)

	shift.Initialize(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
