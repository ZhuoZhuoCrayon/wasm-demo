pb:
	# protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative timeseriesquery/time_series_query.proto
	trpc create -p timeseriesquery/time_series_query.proto -o timeseriesquery --rpconly --nogomod --mock=false
serve:
	go run server/server.go -conf server/trpc_go.yaml

bstream:
	go run client/bstream/client.go -conf client/trpc_go.yaml

cstream:
	go run client/cstream/client.go -conf client/trpc_go.yaml

sstream:
	go run client/sstream/client.go -conf client/trpc_go.yaml

ustream:
	go run client/ustream/client.go -conf client/trpc_go.yaml
