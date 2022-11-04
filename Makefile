setup:
	go install github.com/golang/mock/mockgen
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	go get github.com/twitchtv/twirp-ruby/protoc-gen-twirp_ruby
	go get github.com/twitchtv/twirp

distributed-mocks:
	mockgen -source=pkg/distributed/models/interfaces.go -destination=pkg/distributed/models/mocks/mocks.gen.go IDealer

distributed-grpc:
	protoc --go_out=. --go_opt=paths=source_relative --twirp_opt=paths=source_relative --twirp_out=. --ruby_out=dealer/lib/pb --twirp_ruby_out=dealer/lib/pb pkg/distributed/remote/pb/fftb.proto

deploy:
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=`git describe`" -o ~/Resilio\ Sync/wailorman_cmder/fftbdev.exe cmd/main/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=`git describe`" -o ~/Resilio\ Sync/kuzmech_cmder/fftbdev.exe cmd/main/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=`git describe`" -o ~/Resilio\ Sync/polardeer_cmder/fftbdev.exe cmd/main/main.go
	go build -ldflags="-X main.version=`git describe`" -o ~/bin/fftbdev cmd/main/main.go
