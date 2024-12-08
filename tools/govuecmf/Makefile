build:
	go build -ldflags="-s -w" govuecmf.go
	$(if $(shell command -v upx || which upx), upx govuecmf)

mac:
	GOOS=darwin go build -ldflags="-s -w" -o govuecmf-darwin govuecmf.go
	$(if $(shell command -v upx || which upx), upx govuecmf-darwin)

win:
	GOOS=windows go build -ldflags="-s -w" -o govuecmf.exe govuecmf.go
	$(if $(shell command -v upx || which upx), upx govuecmf.exe)

linux:
	GOOS=linux go build -ldflags="-s -w" -o govuecmf-linux govuecmf.go
	$(if $(shell command -v upx || which upx), upx govuecmf-linux)
