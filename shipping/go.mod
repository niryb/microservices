module github.com/niryb/microservices/shipping

go 1.24.0

replace github.com/niryb/microservices-proto/golang/shipping => ../../microservices-proto/golang/shipping

require google.golang.org/grpc v1.78.0

require github.com/niryb/microservices-proto/golang/shipping v0.0.0-20260126115315-ddf2abd81ac5

require (
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260122232226-8e98ce8d340d // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
