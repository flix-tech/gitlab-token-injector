install:
	glide install

build-ci: install
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.__VERSION__=${TRAVIS_TAG} -s -w -extldflags '-static'" -o ./dist/gitlab-token-injector
