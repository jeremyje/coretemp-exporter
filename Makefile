GO = go
SOURCE_DIRS=$(shell go list ./... | grep -v '/vendor/')

all: coretemp-exporter.exe

coretemp-exporter.exe: coretemp-exporter.go
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GO) build -o $@ $<

run: coretemp-exporter.go
	$(GO) run coretemp-exporter.go -log=cputemps.log

lint:
	$(GO) fmt
	$(GO) vet

test:
	$(GO) test -race ${SOURCE_DIRS} -cover

coverage.txt:
	for sfile in ${SOURCE_DIRS} ; do \
		go test -race "$$sfile" -coverprofile=package.coverage -covermode=atomic; \
		if [ -f package.coverage ]; then \
			cat package.coverage >> coverage.txt; \
			$(RM) package.coverage; \
		fi; \
	done

clean:
	rm -f coretemp-exporter.exe
	rm -f coverage.txt

.PHONY: all run lint test clean
