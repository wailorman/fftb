setup:
	go install github.com/golang/mock/mockgen
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	go get github.com/twitchtv/twirp-ruby/protoc-gen-twirp_ruby
	go get github.com/twitchtv/twirp

proto:
	protoc --go_out=. --go_opt=paths=source_relative --twirp_opt=paths=source_relative --twirp_out=. --ruby_out=dealer/lib/pb --twirp_ruby_out=dealer/lib/pb pkg/distributed/remote/pb/fftb.proto

build:
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=`git describe`" -o dist/windows_amd64/fftb.exe cmd/main/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=`git describe`" -o dist/linux_amd64/fftb cmd/main/main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-X main.version=`git describe`" -o dist/linux_arm64/fftb cmd/main/main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.version=`git describe`" -o dist/darwin_amd64/fftb cmd/main/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.version=`git describe`" -o dist/darwin_arm64/fftb cmd/main/main.go
