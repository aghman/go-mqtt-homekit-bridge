.PHONY: build
build: 
	GOOS=darwin go build -o build/go-mqtt-homekit-bridge_darwin
	GOOS=linux go build -o build/go-mqtt-homekit-bridge_linux
