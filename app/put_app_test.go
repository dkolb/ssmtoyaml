package app_test

import (
	"fmt"
	"testing"

	"gitlab.com/dkub/ssmparams/app"
	"gopkg.in/yaml.v3"
)

func TestImportApp(t *testing.T) {
	a := app.PutApp{
		File:        "",
		Interactive: false,
		RetryLimit:  0,
		Region:      "us-east-1",
	}

	a.Exec()
}

func TestImportUnmarshalData(t *testing.T) {
	a := app.PutApp{
		File:        "",
		Interactive: false,
		RetryLimit:  0,
		Region:      "us-east-1",
	}
	paramTree, err := a.UnmarshalData([]byte(data))
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	data, derr := yaml.Marshal(paramTree)
	if derr != nil {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(string(data))
}

func TestYamlUnmarshalScratch(t *testing.T) {
	data := `---
test: data
more: datatest
other: datatest`
	var node yaml.Node

	err := yaml.Unmarshal([]byte(data), &node)

	fmt.Println(node)
	fmt.Println(err)
}

var data string = `AnotherApplication:
  dev:
      SomeSetting:
          _type: String
          _value: AnotherSetting
Application:
  dev:
      GithubPassword:
          _type: SecureString
          _value: Stuff
          _key: alias/basic-data-symmetric
          _tags:
              ATag: A value
              AnotherTag: Another value
      GithubUsername:
          _type: String
          _value: Stuff
  prod:
      GithubPassword:
          _type: SecureString
          _value: Stuff
      GithubUsername:
          _type: String
          _value: Stuff`
