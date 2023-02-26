projectName ?= openai

.PHONY: help
help:
	@echo "help:"
	@echo "use [ make init ] init project"


.PHONY: init 
init:
	@go mod init
	@go install github.com/spf13/cobra-cli@latest
	@GOWORK=off cobra-cli init .
	@GOWORK=off cobra-cli add text-completion
	
