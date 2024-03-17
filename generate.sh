#!/bin/bash
protoc .\blog-mongo\proto\blog.proto --go_out=. --go-grpc_out=. 
