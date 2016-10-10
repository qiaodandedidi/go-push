export GOPATH:=$(CURDIR)
all: install

install:
	go install push
