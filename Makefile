all: format

fmt_flags = -w=true -tabs=false -tabwidth=2

format:
	$(MAKE) formatter name=batch
	$(MAKE) formatter name=print
	$(MAKE) formatter name=thrasher
	$(MAKE) formatter name=server

server: format
	go run server.go batch.go $(args)

print: format
	go run print.go batch.go $(args)

thrasher: format
	go run thrasher.go batch.go $(args)

formatter:
	gofmt $(fmt_flags) $(name).go
