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
	Short: "Put YAML file of parameters into SSM Parameter store.",
	Long: `This command will import a YAML file into SSM Parameter store. The file
should be of the following format. In this example, a parameter of the name
/Application/Dev/MySetting will be created.

---
Application:
  Dev:
	  MySetting:
		  _type: SecureString
			_value: MySettingValue
			_key: alias/basic-data-symmetric
			_tags:
			  Component: MyApp
				Environment: Dev
				BudgetCode: MYAPP

If provided, a separate YAML file can provide the tags in one place. These tags
will override each _tag value in your input file.  Example tag file:

---
Component: MyApp
Environment: Dev
BudgetCode:> MYAPP
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		importApp := app.PutApp{
			File:        File,
			TagFile:     TagFile,
			Interactive: !InteractiveDisabled,
			RetryLimit:  RetryLimit,
			Region:      Region,
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
	TagFile             string
	InteractiveDisabled bool
	InteractiveEnabled  bool
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
		"in-file",
		"i",
		"",
		"File to import into SSM.",
	)

	putCmd.Flags().StringVarP(
		&TagFile,
		"tags",
		"t",
		"",
		"A file containing a YAML map of TagName: TagValue pairs to add to all parameters",
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
