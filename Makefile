INSTALL_SCRIPT=./install.sh
BIN_FILE=./shellreminders

install:
	go build -o shellreminders
	${INSTALL_SCRIPT}
	cp ${BIN_FILE} ~/bin

clean:
	go clean

check:
	go test

cover:
	go test -coverprofile cp.out
	go tool cover -html=cp.out

test:
	go test
