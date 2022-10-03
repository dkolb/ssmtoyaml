package app

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"gopkg.in/yaml.v3"

	"gitlab.com/dkub/ssmparams/types"
	"gitlab.com/dkub/ssmparams/utils"
)

type ImportApp struct {
	File                string
	InteractiveDisabled bool
	RetryLimit          int
	Region              string
	client              *ssm.Client
}

func (im *ImportApp) Exec() error {
	var err error
	im.client, err = utils.InitializeSsmClient(&im.Region)
	if err != nil {
		log.Println("failed to init ssm client", err)
		return err
	}

	data, err := os.ReadFile(im.File)
	if err != nil {
		log.Println("failed to init ssm client", err)
		return err
	}

	var paramTree *types.ParameterTree
	paramTree, err = im.UnmarshalData(data)

	fmt.Println(paramTree)

	return nil
}

func (im *ImportApp) UnmarshalData(data []byte) (*types.ParameterTree, error) {
	paramTree := types.NewParameterTree()
	err := yaml.Unmarshal(data, paramTree)
	fmt.Println(err)
	return paramTree, err
}
