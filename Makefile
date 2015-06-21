test: lint go-test

lint:
	./lint.sh

go-test:
	go test
bench:
	go test -bench . -benchmem

coverage:
	out=`mktemp -t "forge-coverage"`; go test -coverprofile=$$out && go tool cover -html=$$out

.PHONY: lint test
