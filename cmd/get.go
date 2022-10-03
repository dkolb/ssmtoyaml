/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/dkub/ssmparams/app"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieves an entire tree of your SSM param store as a YAML document.",
	Long: `This command will retrieve a tree or subtree beginning at --ssm_root 
into a well structured YAML document for ease of editing or copying between
environments.`,
	Run: func(cmd *cobra.Command, args []string) {
		exportApp := &app.ExportApp{
			SsmPathRoot:    PathRoot,
			ExportFile:     FilePath,
			Decrypt:        Decrypt,
			ForceOverwrite: ForceOverwrite,
			Region:         Region,
		}

		err := exportApp.Exec()
		if err != nil {
			panic(err)
		} else {
			os.Exit(0)
		}
	},
}

var (
	FilePath       string
	Decrypt        bool
	PathRoot       string
	ForceOverwrite bool
)

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	// Note that these Flags are from github.com/spf13/pflag and not the builtin
	// flag pkg.
	getCmd.Flags().StringVarP(
		&FilePath,
		"file",
		"f",
		"./ssmparams_out.yaml",
		"The file to write YAML commands out to.",
	)

	getCmd.Flags().BoolVarP(
		&Decrypt,
		"decrypt",
		"d",
		false,
		"Set to decrypt SecureString values.",
	)

	getCmd.Flags().StringVarP(
		&PathRoot,
		"ssm-root",
		"r",
		"/",
		"A path root to retrieve from.",
	)

	getCmd.MarkFlagFilename("file")
}
