
linux:
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -trimpath -ldflags="-s -w" -installsuffix cgo -o tshell-linux main.go
mac:
	GOOS=darwin CGO_ENABLED=0 GOARCH=amd64 go build -trimpath -ldflags="-s -w" -installsuffix cgo -o tshell-mac main.go
windows:
	GOOS=windows CGO_ENABLED=0 GOARCH=amd64 go build -trimpath -ldflags="-s -w" -installsuffix cgo -o tshell.exe main.go
