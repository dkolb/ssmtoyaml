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

