package app_test

import (
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
	a := &app.ExportApp{
		SsmPathRoot:    "/PingFederate/e2e",
		ExportFile:     "test.yaml",
		Decrypt:        true,
		ForceOverwrite: true,
		Region:         "us-east-1",
	}
	a.Exec()
}
