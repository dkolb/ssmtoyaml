## ssmparams

A program for managing your SSM params as YAML files.

### Synopsis

This program allows you to export your SSM parameters into a
YAML file that represents their Path-Like naming structure and manage their
values and some attributes.

This is a rewrite of a ruby gem I also authored.

### Options

```
  -h, --help            help for ssmparams
      --region string   AWS Region to run against. (default "us-east-1")
```

### SEE ALSO

* [ssmparams get](docs/ssmparams_get.md)	 - Retrieves an entire tree of your SSM param store as a YAML document.
* [ssmparams put](docs/ssmparams_put.md)	 - Put YAML file of parameters into SSM Parameter store.

## ssmparams get

Retrieves an entire tree of your SSM param store as a YAML document.

### Synopsis

This command will retrieve a tree or subtree beginning at --ssm_root 
into a well structured YAML document for ease of editing or copying between
environments.

```
ssmparams get [flags]
```

### Options

```
  -d, --decrypt           Set to decrypt SecureString values.
  -f, --force-overwrite   Overwrite the --out-file if it exists.
  -h, --help              help for get
      --ignore-tags       Do not write _tags keys to the output file.
  -o, --out-file string   The file to write YAML commands out to. (default "./ssmparams_out.yaml")
  -r, --ssm-root string   A path root to retrieve from. (default "/")
```

### Options inherited from parent commands

```
      --region string   AWS Region to run against. (default "us-east-1")
```

### SEE ALSO

* [ssmparams](docs/ssmparams.md)	 - A program for managing your SSM params as YAML files.

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

