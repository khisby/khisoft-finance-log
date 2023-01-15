# Khisoft Finance Log

## Env Var
GSHEET_ID -> from url spreadsheet after `https://docs.google.com/spreadsheets/d/` before `/edit` 
WHITELIST_USER array of string -> whatsapp number using 62(Indonesia ID) separated by comma

## How to Build
Linux `GOOS=linux GOARCH=amd64 go build -o bin/finance-log main.go`
Mac `GOOS=darwin GOARCH=arm64 go build -o bin/finance-log main.go`
Windows `GOOS=windows GOARCH=amd64 go build -o bin/finance-log main.go`

## How to run 
Build `./finance-log`
Source Code `go run main.go`