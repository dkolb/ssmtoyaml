package app_test

import (
	"gitlab.com/dkub/ssmparams/app"
	"testing"
)

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
