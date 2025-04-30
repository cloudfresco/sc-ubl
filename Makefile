
SHELL := /bin/bash

# The name of the executable (default is current directory name)
#TARGET := $(shell echo $${PWD\#\#*/})
#.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
#VERSION := 1.0.0
#VERSION          := $(shell git describe --tags --always --dirty="-dev")
#DATE             := $(shell date -u '+%Y-%m-%d-%H%M UTC')
#VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'
#BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
#LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory

SRC = $(shell find . -type f -name '*.go' -not -path "./internal/protogen/*")
SRCPROTO = $(shell find . -type f -name '*.proto'")
MFILE = cmd/main.go
EXEC = cmd/sc-ubl
PKGS = $(go list ./... | grep -v /proto/ | grep -v /proto-gen/)

.PHONY: all chk lint proto build buildp test clean fmt gocritic staticcheck errcheck revive golangcilint protofmt protolint tidy pkgupd run runp doc

all: chk protoc buildp

chk: goimports fmt gocritic staticcheck errcheck protofmt protolint

rev: revive

lint: golangcilint
 
#protoc:
#  protoc \
#    --go_out=. \
#    --go_opt=module=github.com/cloudfresco/sc-dcsa \
#    $(SRCPROTO)


build: 
	@echo "Building sc-ubl"	
	@go build -o $(EXEC) $(MFILE)

buildp:
	@echo "Building sc-ubl"	
	@go build -o $(EXEC) $(MFILE)

test:
	@mysql -uroot -p$(SC_UBL_DBPASSROOT) -e 'DROP DATABASE IF EXISTS  $(SC_UBL_DBNAME_TEST);'
	@mysql -uroot -p$(SC_UBL_DBPASSROOT) -e 'CREATE DATABASE $(SC_UBL_DBNAME_TEST);'
	@mysql -uroot -p$(SC_UBL_DBPASSROOT) -e "GRANT ALL ON *.* TO '$(SC_UBL_DBUSER_TEST)'@'$(SC_UBL_DBHOST)';"
	@mysql -uroot -p$(SC_UBL_DBPASSROOT) -e 'FLUSH PRIVILEGES;'
	@mysql -u$(SC_UBL_DBUSER_TEST) -p$(SC_UBL_DBPASS_TEST)  $(SC_UBL_DBNAME_TEST) < sql/mysql/sc_ubl_mysql_schema.sql

	@echo "Starting tests"
	go test -v github.com/cloudfresco/sc-ubl/internal/controllers/taxcontrollers
	#@for pkg in $$(go list ./...); do echo "Testing" $$pkg && go test -v $$pkg; done		

clean:
	@rm -f $(EXEC)

goimports:
	@echo "Running goimports"		
	@goimports -l -w $(SRC)

fmt:
	@echo "Running gofumpt"
	@gofumpt -l -w .
	@echo "Running gofmt"		
	@gofmt -s -l -w $(SRC)

gocritic:
	@echo "Running gocritic"
	@gocritic check $(SRC)

staticcheck:
	@echo "Running staticcheck"
	@staticcheck ./...

errcheck:
	@echo "Running errcheck"
	@errcheck ./...

revive:
	@echo "Running revive"
	@revive $(SRC)

golangcilint:
	@echo "Running golangci-lint"
	@golangci-lint run

protofmt:
	@echo "Running protofmt"
	cd internal/proto && buf format -w

protolint:
	@echo "Running protolint"
	@buf lint

tidy:
	go mod tidy -v -e

pkgupd:
	go get -u ./...
	go mod tidy -v -e

run: build
	@echo "Starting sc-ubl"	
	@./$(EXEC) --dev 

runp: buildp	
	@echo "Starting sc-ubl"	
	@./$(EXEC) 

doc: 

