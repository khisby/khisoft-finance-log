# Khisoft Finance Log
Financial Record can store your income and outcome finance data with daily or monthly reports.  
using WhatsappBot and Spreadsheet API

## Env Var & Credentials
GSHEET_ID -> from url spreadsheet after `https://docs.google.com/spreadsheets/d/` before `/edit`  
WHITELIST_USER array of string -> whatsapp number using 62(Indonesia ID) separated by comma  

You need to download credential.json from your service account to write spreadsheet

## How to Build
Linux `GOOS=linux GOARCH=amd64 go build -o bin/finance-log main.go`  
Mac `GOOS=darwin GOARCH=arm64 go build -o bin/finance-log main.go`  
Windows `GOOS=windows GOARCH=amd64 go build -o bin/finance-log main.go`  

## How to run 
Build `./finance-log`  
Source Code `go run main.go`

## Screenshot
![image](https://user-images.githubusercontent.com/24775167/233133633-b61423d2-3645-4a80-97c4-64a98d3d8cd5.png)
![image](https://user-images.githubusercontent.com/24775167/233133464-92591fb6-b797-4760-ad2e-d91eb03ec76f.png)
![image](https://user-images.githubusercontent.com/24775167/233134050-97f9fd58-88b1-4ecc-97bf-2518bdd7a142.png)
![image](https://user-images.githubusercontent.com/24775167/233133192-b42ca26c-ca0b-47ca-817b-55eaca09bb96.png)
![image](https://user-images.githubusercontent.com/24775167/233133355-afe15f0d-ad2b-4752-b387-987067cd938b.png)

