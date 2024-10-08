proto:
	[ -d "./pb" ] && rm -rf pb/* || mkdir pb
	protoc --proto_path=protos --go_out=pb \
	--go_opt=paths=source_relative --go-grpc_out=pb \
	--go-grpc_opt=paths=source_relative \
	protos/*.proto \
