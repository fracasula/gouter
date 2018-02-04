.PHONY: build run

IMAGENAME=gohttp

default: build run

build:
	docker build -t $(IMAGENAME) .

run:
	docker run --rm -it -p 8080:8080 $(IMAGENAME)
