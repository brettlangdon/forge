test: lint go-test

lint:
	./lint.sh

go-test:
	go test
bench:
	go test -bench .

.PHONY: lint test
