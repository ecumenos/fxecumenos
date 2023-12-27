GOPRIVATE=github.com/ecumenos
SHELL=/bin/sh

.PHONY: all
all: tidy check fmt lint test mock tidy

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test-short
test-short: ## Run tests, skipping slower integration tests
	go test -test.short ./...

.PHONY: test-interop
test-interop: ## Run tests, including local interop (requires services running)
	go clean -testcache && go test -tags=localinterop ./...

.PHONY: coverage-html
coverage-html: ## Generate test coverage report and open in browser
	go test ./... -coverpkg=./... -coverprofile=test-coverage.out
	go tool cover -html=test-coverage.out

.PHONY: lint
lint: ## Verify code style and run static checks
	go vet -asmdecl -assign -atomic -bools -buildtag -cgocall -copylocks -httpresponse -loopclosure -lostcancel -nilfunc -printf -shift -stdmethods -structtag -tests -unmarshal -unreachable -unsafeptr -unusedresult ./...
	test -z $(gofmt -l ./...)

.PHONY: fmt
fmt: ## Run syntax re-formatting (modify in place)
	go fmt ./...

.PHONY: check
check: ## Compile everything, checking syntax (does not output binaries)
	go build ./...

.PHONY: mock
mock: mock_clean
	go generate ./...

.PHONY: mock_clean
mock_clean:
	find . -name "*.go" -path "**/mocks/*" | while read file; do rm $$file; done;
