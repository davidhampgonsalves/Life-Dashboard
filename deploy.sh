#!/bin/sh
terraform init && 
GOOS=linux GOARCH=amd64 go build -o ldb cmd/main/handler.go &&
rm -f ldb.zip &&
zip -r ldb.zip ldb static creds &&
terraform apply