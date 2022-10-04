## ssmparams put

Put YAML file of parameters into SSM Parameter store.

### Synopsis

This command will import a YAML file into SSM Parameter store. The file
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


```
ssmparams put [flags]
```

### Options

```
  -h, --help              help for put
  -i, --in-file string    File to import into SSM.
      --no-interact       Use to disable Y/N check.
      --retry-limit int   Limit on retries for failed parameter updates. (default 3)
  -t, --tags string       A file containing a YAML map of TagName: TagValue pairs to add to all parameters
```

### Options inherited from parent commands

```
      --region string   AWS Region to run against. (default "us-east-1")
```

### SEE ALSO

* [ssmparams](docs/ssmparams.md)	 - A program for managing your SSM params as YAML files.

