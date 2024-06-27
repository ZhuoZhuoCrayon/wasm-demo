APP    := collector
TARGET := tars

GOBUILD      := go build
DOCKER_BUILD := docker build

REPO         ?= mirrors.tencent.com/crayon
VERSION  ?= $(shell date "+%Y%m%d%H%M%S")
# VERSION  ?= mac

LOWCASE_TARGET := $(shell echo $(TARGET) | tr '[:upper:]' '[:lower:]')
IMG_REPO       := $(REPO)/$(LOWCASE_TARGET)
