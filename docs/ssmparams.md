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

