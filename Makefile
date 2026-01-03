# Contains commands that are used during development

# run tests and ignore output that indicates that directories don't have a folder
test:
	go test ./... 2>&1 | grep -v "\[no test files\]" | grep -v "no tests to run"

# count-tests print the number of unit tests in this repo, not including tests of dependencies.
count-tests:
	go test -v ./... | grep "\-\-\- PASS: Test" | wc -l

# run all benchmarks and only display the benchmark results
benchmark:
	go test -bench=. ./... | grep -E "\bBenchmark"

# run user facing ecs benchmarks
benchmark-ecs:
	go test -bench=. ./benchmark/ | grep -E "\bBenchmark"

lint: verify-lint-config
	go tool golangci-lint run

fix-lint:
	go tool golangci-lint run --fix

verify-lint-config:
	go tool golangci-lint config verify

# run all verifications that are recommended to run before pushing code
verify: 
	make lint
	make test
	go mod tidy -diff
