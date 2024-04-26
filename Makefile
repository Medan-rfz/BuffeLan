
gen-syso-file:
	rsrc -ico ./resources/app_icon.ico -o ./cmd/rsrc.syso

build-all: gen-syso-file
	go build -ldflags -H=windowsgui -o ./bin/BuffeLan.exe ./cmd

