app=crayon
server=timeseriesquery
SERVICE=timeseriesquery
USER := $(shell git config user.name) # 部署用户名
USERNAME ?=
PASSWORD ?=
IMAGE_URL ?= mirrors.tencent.com/crayon/timeseriesquery
IMAGE_VERSION ?= 1.0.4

setup:
	go install trpc.group/trpc-go/trpc-cmdline/trpc@latest && trpc setup
pb:
	trpc-pub create -p timeseriesquery/time_series_query.proto -o timeseriesquery --rpconly --nogomod --mock=false
serve:
	go run main.go -conf conf/trpc_go.yaml

bstream:
	go run client/bstream/client.go -conf client/trpc_go.yaml

cstream:
	go run client/cstream/client.go -conf client/trpc_go.yaml

sstream:
	go run client/sstream/client.go -conf client/trpc_go.yaml

ustream:
	go run client/ustream/client.go -conf client/trpc_go.yaml

local:
	@echo "download dependencies"
	@go mod tidy
	@go mod download
	@go mod vendor
	@echo "========= local build ========="
	sh ./stke_build/build.sh -v enable

# make docker-image USERNAME= PASSWORD= IMAGE_URL=csighub.tencentyun.com/xxx IMAGE_VERSION=xxx
docker-image:
	@echo "download dependencies"
	@go mod tidy
	@go mod download
	@go mod vendor
	@echo "========= linux build ========="
	sh ./stke_build/build.sh -v enable -a linux
#	@echo "docker login"
#	@echo $(PASSWORD) | docker login --username $(USERNAME) --password-stdin csighub.tencentyun.com
	@echo "build image"
	docker build --no-cache --pull -f ./stke_build/image/Dockerfile . -t $(SERVICE):latest
	@echo "tags image"
	docker tag $(SERVICE):latest $(IMAGE_URL):$(IMAGE_VERSION)
	@echo "push image to remote"
	docker push $(IMAGE_URL):$(IMAGE_VERSION)
