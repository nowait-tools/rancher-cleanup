.PHONY: push build

PWD := `pwd`

build:
	docker run -it \
        -v $(PWD):/go/src/app \
        -e "GOOS=linux" \
        -e "GOARCH=amd64" \
        -w /go/src/app golang:1.7.0-alpine \
        go build -o rancher-cleanup
	docker build -t nowait/rancher-cleanup:1.0 .

push:
	docker push \
	nowait/rancher-cleanup:1.0

