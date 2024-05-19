.PHONY: release

release: .release
.release: *.go static/* static/templates/*
	go test ./...
	commit=$$(git rev-parse HEAD); \
	GOOS=linux \
	GOARCH=arm64 \
	go build \
		-ldflags="-X 'main.CommitHash=$$commit' -X 'main.BuildTimestamp=$$(date -Iseconds --utc)'" \
		-o bootstrap \
		.
	zip kinkdiff.zip bootstrap
	aws s3 cp kinkdiff.zip s3://kinkdiff-symme-link/kinkdiff.zip
	aws lambda update-function-code --function-name kinkdiff --s3-bucket kinkdiff-symme-link --s3-key kinkdiff.zip
	touch .release

.git/config: gitconfig
	cp gitconfig .git/config
