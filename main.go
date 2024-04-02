/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/dkolb/ssmtoyaml/cmd"
	"github.com/dkolb/ssmtoyaml/utils"
	"github.com/spf13/cobra/doc"
)

func main() {
	if os.Getenv("GEN_DOCS") == "makeitso" {
		generateDocs()
	} else {
		cmd.Execute()
	}
}

func generateDocs() {
	dir := utils.EnvWithDefault("GEN_DOCS_DIR", "./docs")
	linkPrefix := utils.EnvWithDefault("GEN_DOCS_LINK_PREFIX", "docs/")
	ssmtoyaml := cmd.GetRootCmd()
	err := doc.GenMarkdownTreeCustom(
		ssmtoyaml,
		dir,
		func(filename string) string { return "" },
		func(name string) string { return linkPrefix + name },
	)
	if err != nil {
		panic(err)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Doc generation failed with error: %v\n", err)
	}
	fmt.Fprintln(os.Stderr, "Docs generated.")
}
