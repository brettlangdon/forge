test: lint go-test

lint:
	./lint.sh

go-test:
	go test

.PHONY: lint test
