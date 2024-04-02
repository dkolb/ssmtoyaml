module github.com/dkolb/ssmtoyaml

go 1.22

retract (
	v0.1.4 //Bad module paths
	v0.1.5 //Bad dependency combination in aws-sdk-go-v2
)

require (
	github.com/aquasecurity/table v1.8.0
	github.com/aws/aws-sdk-go-v2 v1.26.1
	github.com/aws/aws-sdk-go-v2/config v1.27.10
	github.com/aws/aws-sdk-go-v2/service/ssm v1.49.5
	github.com/hashicorp/go-multierror v1.1.1
	github.com/liamg/tml v0.6.1
	github.com/manifoldco/promptui v0.9.0
	github.com/spf13/cobra v1.8.0
)

require (
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.2 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/term v0.0.0-20220526004731-065cf7ba2467 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.17.10 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/iam v1.31.4
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.20.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.23.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.28.6
	github.com/aws/smithy-go v1.20.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1
)
