.POSIX:

GO=go

check:
	$(GO) vet ./...
	gocyclo --over 10 .
	staticcheck --checks=all ./cache/ ./errors/ ./log/ ./netutil/ ./stringutil/ ./timeutil/
	$(GO) test --count 1 --race ./...
