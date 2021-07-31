Code for https://www.udemy.com/course/grpc-golang

Uses updated version of protobuf from the course which slightly changes things

Generate greet protobuf using

```
protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     greet/greetpb/greet.proto
```