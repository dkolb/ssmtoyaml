docs: build
	GEN_DOCS=makeitso GEN_DOCS_DIR=./docs ./ssmparams
	cat docs/ssmparams.md docs/ssmparams_*.md > README.md

build:
	go build