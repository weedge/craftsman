
# mysql
MYSQL_VER ?= $(shell mysql --version | awk -F '[ |.]' '{print $$4}')
MYSQL_USER ?= root
MYSQL_PASSWD ?= 123
MYSQL_HOST ?= 127.0.0.1
MYSQL_PORT ?= 3306
MYSQL_DB ?= pay

# refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
MYSQL_DSN ?= $(MYSQL_USER):$(MYSQL_PASSWD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DB)?charset=utf8mb4&parseTime=True&loc=Local

# guides https://gorm.io/gen/index.html
# refer https://gorm.io/gen/gen_tool.html 
GORM_GEN_DBTABLES ?= user_assert,user_assert_record
GORM_GEN_OUTFILE ?= pay_gen.go
GORM_GEN_OUTPATH ?= ./internal/da/dao

.PHONY: help
help:
	@echo "help:\n"
	@echo "use [ make gorm-gen ] cmd to generate da dao,model\n"
	@echo "use [ make mysql-dump-table] cmd to dump create table sql\n"

.PHONY: mysql-dump-table
mysql-dump-table:
	@mysqldump -u$(MYSQL_USER) -p$(MYSQL_PASSWD) -h$(MYSQL_HOST) -P$(MYSQL_PORT) -d $(MYSQL_DB) > ./data/$(MYSQL_DB)-mysql$(MYSQL_VER)-innodb.sql

.PHONY: gorm-gen
gorm-gen:
	@go install gorm.io/gen/tools/gentool@latest
	@gentool -dsn "$(MYSQL_DSN)" -tables "$(MYSQL_DBTABLES)" -outFile "$(GORM_GEN_OUTFILE)" -outPath "$(GORM_GEN_OUTPATH)" -fieldWithTypeTag true
