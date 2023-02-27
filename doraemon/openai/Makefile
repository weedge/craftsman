projectName ?= openai

.PHONY: help
help:
	@echo "help:"
	@echo "use [ make init ] init project"
	@echo "use [ make test-text-model ] test text models"
	@echo "use [ make run ] run text completion chat"

.PHONY: init 
init:
	@go mod init
	@go install github.com/spf13/cobra-cli@latest
	@GOWORK=off cobra-cli init .
	@GOWORK=off cobra-cli add text-completion
	

.PHONY: test-text-model
test-text-model:
	@go run main.go textCompletion -p "what gpt-3" -M "text-davinci-003"
	@go run main.go textCompletion -p "what gpt-3" -M "text-davinci-002"
	@go run main.go textCompletion -p "what gpt-3" -M "text-davinci-001"
	@go run main.go textCompletion -p "what gpt-3" -M "curie-instruct-beta"
	@go run main.go textCompletion -p "what gpt-3" -M "curie"
	@go run main.go textCompletion -p "what gpt-3" -M "ada"
	@go run main.go textCompletion -p "what gpt-3" -M "text-curie-001"
	@go run main.go textCompletion -p "what gpt-3" -M "text-babbage-001"
	@go run main.go textCompletion -p "what gpt-3" -M "text-ada-001"
	@go run main.go textCompletion -p "what gpt-3" -M "davinci-instruct-beta"
	@go run main.go textCompletion -p "what gpt-3" -M "davinci"
	@go run main.go textCompletion -p "what gpt-3" -M "babbage"

.PHONY: run
run:
	@go run main.go textCompletion 