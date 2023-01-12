SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm
go build -trimpath -ldflags "-w -s" -tags="androidgodns"
