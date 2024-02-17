release: .release
.release: main.go
	GOOS=linux GOARCH=arm64 go build -o bootstrap main.go
	zip kinkdiff.zip bootstrap
	aws s3 cp kinkdiff.zip s3://kinkdiff-symme-link/kinkdiff.zip
	touch .release

.git/config: gitconfig
	cp gitconfig .git/config
