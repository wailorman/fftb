distributed-client:
	oapi-codegen --config pkg/distributed/remote/schema/dealer/generator_config.yml pkg/distributed/remote/schema/dealer/swagger.yml
	oapi-codegen --config pkg/distributed/remote/schema/contracter/generator_config.yml pkg/distributed/remote/schema/contracter/swagger.yml

distributed-mocks:
	mockgen -source=pkg/distributed/models/interfaces.go -destination=pkg/distributed/models/mocks/mocks.gen.go IDealer,IContracter

deploy:
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=`git describe`" -o ~/Resilio\ Sync/wailorman_cmder/fftbdev.exe cmd/main/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=`git describe`" -o ~/Resilio\ Sync/kuzmech_cmder/fftbdev.exe cmd/main/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=`git describe`" -o ~/Resilio\ Sync/polardeer_cmder/fftbdev.exe cmd/main/main.go
	go build -ldflags="-X main.version=`git describe`" -o ~/bin/fftbdev cmd/main/main.go
