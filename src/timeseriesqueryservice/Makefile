pb:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative timeseriesquery/time_series_query.proto

serve:
	go run server/server.go

bstream:
	go run client/bstream/client.go

cstream:
	go run client/cstream/client.go

sstream:
	go run client/sstream/client.go
