.PHONY: assets
assets:
	mkdir -p 'assets/web'
	cp 'resources/icon.png' 'assets/web/favicon.png'
	cp 'resources/web/main.css' 'assets/web/main.css'
	cp 'resources/web/main.js' 'assets/web/main.js'
	cp 'resources/web/index.html' 'assets/index.html'
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" 'assets/web'
	go run ./cmd/rallypack/

.PHONY: play
play:
	go run ./cmd/rallymka/

.PHONY: debug
debug:
	go run -tags debug ./cmd/rallymka/

.PHONY: wasm
wasm:
	GOOS=js GOARCH=wasm go build -o './assets/web/main.wasm' './cmd/rallywasm/'

.PHONY: web
web:
	go run github.com/mokiat/httpserv@master -dir './assets' -host 127.0.0.1
