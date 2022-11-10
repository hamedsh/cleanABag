GOCMD := CGO_ENABLED=0 go
BINARY := cleanABag
BINDIR := ./bin
VERSION := 0.1.0

GOLDFLAGS := -s -w -X main.Version=$(VERSION)

.PHONY: build
build:
	${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}

.PHONY: buildAll
buildAll:
	GOOS=linux GOARCH=amd64 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_linux_amd64
	GOOS=linux GOARCH=arm ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_linux_arm
	GOOS=linux GOARCH=arm64 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_linux_arm64
	GOOS=linux GOARCH=386 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_linux_386

.PHONY: clean
clean:
	rm -f ${BINDIR}/${BINARY}

fmt:
	go fmt ./...

.PHONY: release
release:
	echo "Tagging version ${VERSION}"
	git tag -a v${VERSION} -m "New released tag: v${VERSION}"
	GOOS=linux GOARCH=amd64 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_arm
	GOOS=linux GOARCH=arm64 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_arm64
	GOOS=linux GOARCH=386 ${GOCMD} build -ldflags "$(GOLDFLAGS)" -o ${BINDIR}/${BINARY}_${VERSION}_linux_386

.PHONY: dependencies
dependencies:
	${GOCMD} get github.com/Strubbl/wallabago/v7
	${GOCMD} get github.com/spf13/cobra
