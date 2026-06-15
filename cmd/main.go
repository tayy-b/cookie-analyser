package main

import (
	"cookie-analyser/internal"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cookie-analyser",
	Short: "a simple cli tool to analyse cookie usage",
	Long:  "cookie-analyser is a cli tool that currently returns the most active cookie for a given day, parsed from a csv file.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLogger()
	},
	RunE: internal.RunCookieFinder,
}

func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logrus.SetLevel(logrus.InfoLevel)
}

func init() {
	rootCmd.Flags().StringP("file", "f", "", "csv file path")
	rootCmd.Flags().StringP("date", "d", "", "date in UTC format YYYY-MM-DD")
	rootCmd.MarkFlagRequired("file")
	rootCmd.MarkFlagRequired("date")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("%s", err.Error())
		os.Exit(1)
	}
}
