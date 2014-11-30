build:
	go build R.go

test: build
	./test.sh
