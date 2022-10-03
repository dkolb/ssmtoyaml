/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/dkub/ssmparams/app"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		importApp := app.ImportApp{
			File:                File,
			InteractiveDisabled: InteractiveDisabled,
			RetryLimit:          RetryLimit,
			Region:              Region,
		}
		err := importApp.Exec()
		if err != nil {
			panic(err)
		} else {
			os.Exit(0)
		}
	},
}

var (
	File                string
	InteractiveDisabled bool
	RetryLimit          int
)

func init() {
	rootCmd.AddCommand(putCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// putCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	// Note that these Flags are from github.com/spf13/pflag and not the builtin
	// flag pkg.
	// putCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	putCmd.Flags().StringVarP(
		&File,
		"file",
		"f",
		"",
		"File to import into SSM.",
	)
	putCmd.Flags().BoolVar(
		&InteractiveDisabled,
		"no-interact",
		false,
		"Use to disable Y/N check.",
	)
	putCmd.Flags().IntVar(
		&RetryLimit,
		"retry-limit",
		3,
		"Limit on retries for failed parameter updates.",
	)

	putCmd.MarkFlagFilename("file")
	putCmd.MarkFlagRequired("file")
}
