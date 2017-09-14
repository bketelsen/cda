APPNAME=cda
VERSION_TAG=`git describe 2>/dev/null | cut -f 1 -d '-' 2>/dev/null`
COMMIT_HASH=`git rev-parse --short=8 HEAD 2>/dev/null`
BUILD_TIME=`date +%FT%T%z`
LDFLAGS=-ldflags "-s -w \
    -X github.com/bketelsen/cda/cmd.CommitHash=${COMMIT_HASH} \
    -X github.com/bketelsen/cda/cmd.BuildTime=${BUILD_TIME} \
    -X github.com/bketelsen/cda/cmd.Tag=${VERSION_TAG}"

all: fast

clean:
	go clean
	rm ./${APPNAME} || true
	rm -rf ./target || true

build: clean linux darwin windows

fast:
	go build -o ${APPNAME} ${LDFLAGS}

linux:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ./target/linux_amd64/${APPNAME}

darwin:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ./target/darwin_amd64/${APPNAME}

windows:
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o ./target/windows_386/${APPNAME}.exe
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ./target/windows_amd64/${APPNAME}.exe
