# golang-grpc-poc

## Introduction
A poc UDP bi direction boardcast communcation 

## Prerequisite
- Install golang 1.19+

## How to Play
```
// navigate to dir ./src

// Start the listener server
go run main.go

// Start the sender 
// Open another terminal

go run sender/main.go  -message {your message} -k {repeat counter}

```
