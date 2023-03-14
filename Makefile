linux:
	env GOOS=linux go build --ldflags="-s -w" -o ./out/linux/ivy cmd/cli/main.go
darwin:
	env GOOS=darwin go build -ldflags="-s -w" -o ./out/darwin/ivy cmd/cli/main.go
windows:
	env GOOS=windows go build -ldflags="-s -w" -o ./out/windows/ivy cmd/cli/main.go

build: linux darwin windows

server.darwin:
	env GOOS=darwin go build -o ./out/darwin/ivy_server cmd/server/main.go

run.server: server.darwin
	bash ./install_js.sh
	out/darwin/ivy_server

release.cli: build
	zip -r out/ivy_windows.zip out/windows/
	zip -r out/ivy_linux.zip out/linux/
	zip -r out/ivy_darwin.zip out/darwin/
