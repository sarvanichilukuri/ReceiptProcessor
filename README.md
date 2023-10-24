# ReceiptProcessor
Take home Test for ReceiptProcessor for Fetch

Expectations:

1) Go has to installed
2) My project/folder name is ReceiptProcessingServer
3) A module has to be created to manage the dependencies , I have created it using the command: go mod init example/ReceiptProcessingServer
4) Please install gin if not present using the command : go get github.com/gin-gonic/gin

To execute:
The application could be run using the command : go run main.go

My port is listening on localhost:8080

The application supports the below URLs:

POST request - localhost:8080/receipts/process

GET request - localhost:8080/receipts/:id/points
