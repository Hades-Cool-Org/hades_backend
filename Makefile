# ...

ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

gen_docs_md:
  go run .\app\. -docs=markdown
