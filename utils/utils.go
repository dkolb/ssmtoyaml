package utils

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func AwsLoadConfig(region *string) (*aws.Config, error) {
	var (
		conf aws.Config
		err  error
	)

	if region == nil {
		conf, err = config.LoadDefaultConfig(
			context.TODO(),
		)
	} else {
		conf, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(*region),
		)
	}
	if err != nil {
		log.Println("error loading default config:", err)
	}
	return &conf, err
}

func EnvWithDefault(key string, defaultValue string) string {
	v, hasV := os.LookupEnv(key)
	if hasV {
		return v
	} else {
		return defaultValue
	}
}
