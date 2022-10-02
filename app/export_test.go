package app_test

import (
	"fmt"
	"os"
	"testing"

	"gitlab.com/dkub/ssmparams/app"
)

var testYaml string = `AnotherApplication:
    dev:
        SomeSetting:
            value: AnotherSetting
            type: String
Application:
    dev:
        GithubPassword:
            value: Stuff
            type: SecureString
            key: alias/basic-data-symmetric
						tags:
								ATag: A Value
								AnotherTag: Another value
        GithubUsername:
            value: Stuff
            type: String
    prod:
        GithubPassword:
            value: Stuff
            type: SecureString
        GithubUsername:
            value: Stuff
            type: String
`

func TestExportApp(t *testing.T) {
	env := os.Environ()
	for _, e := range env {
		fmt.Println(e)
	}
	a := &app.ExportApp{
		SsmPathRoot:    "/Application",
		ExportFile:     "test.yaml",
		Decrypt:        true,
		ForceOverwrite: true,
		Region:         "us-east-1",
	}
	a.Exec()
}
