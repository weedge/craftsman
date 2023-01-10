

# mysql
MYSQL_VER ?= $(shell mysql --version | awk -F '[ |.]' '{print $$4}')
MYSQL_USER ?= root
MYSQL_PASSWD ?= 123
MYSQL_HOST ?= 127.0.0.1
MYSQL_PORT ?= 3306
MYSQL_DB ?= pay

# refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
MYSQL_DSN ?= $(MYSQL_USER):$(MYSQL_PASSWD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DB)?charset=utf8mb4&parseTime=True&loc=Local


.PHONY: help
help:
	@echo "help:"
	@echo "use [ make gorm-gen ] cmd to generate da dao,model"
	@echo "use [ make wire-gen ] cmd to generate server by using injector"
	@echo "use [ make mysql-dump-table] cmd to dump create table sql"
	@echo "use [ make domain-mockgen or domain-mockery ] cmd to generate interface mock impl for test"
	@echo "use [ make build] cmd to exec go build on local"

.PHONY: mysql-dump-table
mysql-dump-table:
	@mysqldump -u$(MYSQL_USER) -p$(MYSQL_PASSWD) -h$(MYSQL_HOST) -P$(MYSQL_PORT) -d $(MYSQL_DB) > ./data/$(MYSQL_DB)-mysql$(MYSQL_VER)-innodb.sql

# guides https://gorm.io/gen/index.html
# refer https://gorm.io/gen/gen_tool.html 
GORM_GEN_DBTABLES ?= user_asset,user_asset_record
GORM_GEN_OUTFILE ?= pay_gen.go
GORM_GEN_OUTPATH ?= ./internal/da/dao
.PHONY: gorm-gen
gorm-gen:
	@find internal/da/dao -type f -name "*_test*" | xargs rm
	@go install gorm.io/gen/tools/gentool@latest
	@gentool -dsn "$(MYSQL_DSN)" 			\
		-tables "$(GORM_GEN_DBTABLES)" 		\
		-outFile "$(GORM_GEN_OUTFILE)" 		\
		-outPath "$(GORM_GEN_OUTPATH)" 		\
		-withUnitTest 						\
		-fieldWithTypeTag true 				\


# guides https://github.com/google/wire/blob/main/docs/guide.md
.PHONY: wire-gen
wire-gen:
	@go install github.com/google/wire/cmd/wire@latest
	@cd ./internal/da && wire && cd -
	@cd ./internal/station && wire && cd -
	@cd ./internal/gw && wire && cd -


MOCK_SOURCE_DIR ?= ./internal/da/domain
MOCK_DESTINATION_DIR ?= ./internal/da/domain/mocks
.PHONY: domain-mockgen
domain-mockgen:
	@mkdir -p $(MOCK_DESTINATION_DIR)
	@find $(MOCK_DESTINATION_DIR) -type f -name "*_mock.go" | xargs rm -f
	@go install github.com/golang/mock/mockgen@v1.6.0
	@find $(MOCK_SOURCE_DIR) -type f -maxdepth 1 \
	| sed 's#.*/##' \
	| awk -F . '{print $$1}' \
	| xargs -I {} mockgen \
		-source=$(MOCK_SOURCE_DIR)/{}.go \
		-destination=$(MOCK_DESTINATION_DIR)/{}_mock.go \
		-package=mocks


.PHONY: domain-mockery
domain-mockery:
	@find $(MOCK_SOURCE_DIR) -type f -name "i_*.go" | xargs rm -f
	@go install github.com/vektra/mockery/v2@latest
	@mockery --all --case underscore --dir $(MOCK_SOURCE_DIR) --output $(MOCK_DESTINATION_DIR)


# change from caicloud
# Module name.
NAME := payment
# Project main package location.
CMD_DIR := ./cmd
# Project output directory.
OUTPUT_DIR := ./bin
# Build directory.
BUILD_DIR := ./build
# Current version of the project.
GOCOMMON     := $(shell if [ ! -f go.mod ]; then echo $(ROOT)/vendor/; fi)github.com/weedge/craftsman/cloudwego/payment
#VERSION      ?= $(shell git describe --tags --always --dirty)
BRANCH       ?= $(shell git branch | grep \* | cut -d ' ' -f2)
GITCOMMIT    ?= $(shell git rev-parse HEAD)
GITTREESTATE ?= $(if $(shell git status --porcelain),dirty,clean)
BUILDDATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION      ?= latest
appVersion   ?= $(VERSION)

.PHONY: build build-local
build: build-local
build-local:
	@go build -v -o $(OUTPUT_DIR)/$(NAME)                                  \
	  -ldflags "-s -w -X $(GOCOMMON)/version.module=$(NAME)                \
	    -X $(GOCOMMON)/version.version=$(VERSION)                          \
	    -X $(GOCOMMON)/version.branch=$(BRANCH)                            \
	    -X $(GOCOMMON)/version.gitCommit=$(GITCOMMIT)                      \
	    -X $(GOCOMMON)/version.gitTreeState=$(GITTREESTATE)                \
	    -X $(GOCOMMON)/version.buildDate=$(BUILDDATE)"                     \
	  $(CMD_DIR) && chmod +x $(OUTPUT_DIR)/$(NAME);