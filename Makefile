INSTALL_SCRIPT=./install.sh
BIN_FILE=./shellreminders

install:
	go build -o shellreminders
	${INSTALL_SCRIPT}
	cp ${BIN_FILE} ~/bin

clean:
	go clean