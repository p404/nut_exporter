
all: clean build

clean:
	rm -rf nut_expoter

deps:
	dep ensure

build: deps
	go build
