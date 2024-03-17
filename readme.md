  

# GoLang gRPC Project with SSL

  

Welcome to my GoLang repo showcasing gRPC implementation with SSL encryption. This project demonstrates how to build secure communication channels between client and server using gRPC in GoLang.

  
## Features
Basically this is a my first and practise project for goland and grpc in HTTP/2 API. I have implemented following concepts using grpc and protobuf3:

 1. gRPC service definition in .proto files
 2. Server & Client Code in Golang using the protoc gRPC Plugin
 3. Unary, Server Streaming, Client Streaming & Bi-Directional Streaming API
 4. Advanced concepts such as  SSL Security & Encryption 
 5. CRUD API on top of MongoDB
 6. File Streaming using grpc 
## Installation

To use this project, you need to have Go installed on your system. If you haven't installed Go yet, you have to download it first.

Once Go is installed, you can clone this repository:

`git clone git@github.com:bilalhaider789/golang-grpc.git` 

After cloning the repository, navigate to the project directory:

`cd golang-grpc` 

Then, install the dependencies:

`go mod tidy`

SSL certificates maynot work for you. For this disable the tls or generate your own ssl certificates. I have provided the instructions and commands to generate it in ssh folder.
