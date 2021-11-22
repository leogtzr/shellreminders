INSTALL_SCRIPT=./install.sh
BIN_FILE=./remindthis

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
	go test -coverprofile ./cp.out ./shellreminders/...
	go tool cover -html=cp.out

test:
	go test
