# Contains commands that are used during development

# run tests and ignore output that indicates that directories don't have a folder
test:
	go test ./... 2>&1 | grep -v "\[no test files\]" | grep -v "no tests to run"

# run all benchmarks and only display the benchmark results
benchmark:
	go test -bench=. ./... | grep -E "\bBenchmark"

# run user facing ecs benchmarks
benchmark-ecs:
	go test -bench=. ./benchmark/ecs/ | grep -E "\bBenchmark"
