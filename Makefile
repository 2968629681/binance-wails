run: 
	@wails dev

build-darwin-arm64:
	@wails build -clean -platform darwin/arm64

build-linux-arm64:
	@wails build -clean -platform linux/arm64

build-linux-arm:
	@wails build -clean -platform linux/arm

build-linux-amd64:
	@wails build -clean -platform linux/amd64

build-windows-arm64:
	@wails build -clean -platform windows/arm64

build-windows-amd64:
	@wails build -clean -platform windows/amd64

build-windows-386:
	@wails build -clean -platform windows/386

clean:
	@rm -rf build/bin/
