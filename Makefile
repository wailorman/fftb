generate-distributed-client:
	oapi-codegen --config pkg/distributed/schema/generator_config.yml pkg/distributed/schema/swagger.yml

generate-distributed-mocks:
	mockgen -source=pkg/distributed/models/interfaces.go -destination=pkg/distributed/models/mocks/mocks.gen.go IDealer

distribute:
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=`git describe`" -o ~/Resilio\ Sync/wailorman_cmder/fftbdev.exe cmd/main/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=`git describe`" -o ~/Resilio\ Sync/kuzmech_cmder/fftbdev.exe cmd/main/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=`git describe`" -o ~/Resilio\ Sync/polardeer_cmder/fftbdev.exe cmd/main/main.go
	go build -ldflags="-X main.version=`git describe`" -o ~/bin/fftbdev cmd/main/main.go
