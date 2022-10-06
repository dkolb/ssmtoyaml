docs: build
	GEN_DOCS=makeitso GEN_DOCS_DIR=./docs ./ssmtoyaml
	cat docs/ssmtoyaml.md docs/ssmtoyaml_*.md > README.md

build: ssmtoyaml

ssmtoyaml:
	go build

clean:
	rm -f ./ssmtoyaml docs/* ./README.md