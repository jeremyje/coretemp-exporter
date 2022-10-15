GO = go
SOURCE_DIRS=$(shell go list ./... | grep -v '/vendor/')

all: coretemp-exporter.exe

coretemp-exporter.exe: coretemp-exporter.go
	echo GOOS=windows GOARCH=amd64 $(GO) build -o $@ $<

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

.PHONY: test clean
