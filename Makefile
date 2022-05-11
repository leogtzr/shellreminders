.DEFAULT_GOAL := build
INSTALL_SCRIPT=./install.sh
BIN_FILE=./remindthis

build:
	go build -o "${BIN_FILE}"

install:
	go build -o "${BIN_FILE}"
	${INSTALL_SCRIPT}
	cp ${BIN_FILE} ~/bin

clean:
	go clean
	rm --force cp.out

check:
	go test -v ./shellreminders/...	

cover:
	go test -v ./shellreminders/... -coverprofile ./cp.out
	go tool cover -html=cp.out

test:
	go test -v ./shellreminders/...	
