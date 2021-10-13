.POSIX:

GO=go

check:
	$(GO) vet ./...
	gocyclo --over 10 .
	errcheck ./...
	staticcheck --checks=all ./...
	$(MAKE) test

test:
# TODO(a.garipov): Add shuffle in Go 1.17.
	$(GO) test --count 1 --coverprofile='./coverage.txt' --race ./...
