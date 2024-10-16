
build:
	go build -ldflags="-s -w" -trimpath -gcflags="-N -l" -o jv

install:
	mv jv $(GOPATH)/bin
