SERVICE_NAME=cvt
PROJECT_ROOT=$(shell git rev-parse --show-toplevel)

build:
	go build -o $(PROJECT_ROOT)/bin/$(SERVICE_NAME)

exec:
	$(PROJECT_ROOT)/bin/$(SERVICE_NAME) -i=internal/cvt/walktest -o=pkg/img -be=jpg -ae=png
