.PHONY: release clean

release: .release
.release: *.go static/* static/templates/* static/css/kinkdiff.css
	go test ./...
	commit=$$(git rev-parse HEAD); \
	GOOS=linux \
	GOARCH=arm64 \
	go build \
		-ldflags="-X 'main.CommitHash=$$commit' -X 'main.BuildTimestamp=$$(date -Iseconds --utc)'" \
		-o bootstrap \
		.
	rm -f kinkdiff.zip
	zip kinkdiff.zip bootstrap
	aws s3 cp kinkdiff.zip s3://kinkdiff-symme-link/kinkdiff.zip
	aws lambda update-function-code --function-name kinkdiff --s3-bucket kinkdiff-symme-link --s3-key kinkdiff.zip
	touch .release

static/css/kinkdiff.css: nvm bootstrap-5.3.3 kinkdiff.scss
	export NVM_DIR="$$(pwd)/nvm"; \
	. nvm/nvm.sh; \
	nvm install 20; \
	nvm use 20; \
	npm install -g sass; \
	sass kinkdiff.scss static/css/kinkdiff.css
	touch static/css/kinkdiff.css

nvm:
	git clone --depth=1 --branch=v0.39.7 https://github.com/nvm-sh/nvm.git
	touch nvm

bootstrap-5.3.3: bootstrap-5.3.3.zip
	unzip bootstrap-5.3.3.zip
	touch bootstrap-5.3.3

bootstrap-5.3.3.zip:
	curl -L -o bootstrap-5.3.3.zip https://github.com/twbs/bootstrap/archive/v5.3.3.zip

.git/config: gitconfig
	cp gitconfig .git/config

clean:
	rm -r -f \
		 .release \
		 bootstrap \
		 bootstrap-5.3.3 \
		 bootstrap-5.3.3.zip \
		 kinkdiff.zip \
		 nvm \
		 package.json.lock \
		 static/css
