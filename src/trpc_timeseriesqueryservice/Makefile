setup:
	go install trpc.group/trpc-go/trpc-cmdline/trpc@latest && trpc setup
pb:
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
