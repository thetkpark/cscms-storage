protos:
	protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb

protos-js:
	protoc --proto_path=proto proto/*.proto --js_out=import_style=commonjs,binary:client/proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client/proto

.PHONY: protos
