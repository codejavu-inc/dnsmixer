go mod init dnsmixer
go mod tidy
go get golang.org/x/net/publicsuffix@latest
go build -o dnsmixer
