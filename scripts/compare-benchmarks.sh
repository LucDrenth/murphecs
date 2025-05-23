#!/bin/bash

print_help() {
    echo ""
    echo "Options:"
    echo "  -commit           Target commit hash (defaults to latest commit of 'main' branch)"
    echo "  -count            Number of benchmark runs (defaults to 1)"
    echo "  -dir              Path to the benchmark directory (defaults to current '/benchmark')"
    echo "  -func             Benchmark function to run (defaults to '.' to run all functions)"
    echo "  -help             Show this help message"
    echo ""
    exit 0
}

# process args
target_commit_hash=$(git rev-parse --verify main)
number_of_runs=1
benchmark_dir=./benchmark
benchmark_func='.'

while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -commit)
            if [[ -n "$2" ]] && [[ ! "$2" =~ ^- ]]; then
                target_commit_hash="$2"
                shift # past argument
                shift # past value
            fi
            ;;
        -count)
            if [[ -n "$2" ]] && [[ ! "$2" =~ ^- ]]; then
                number_of_runs="$2"
                shift # past argument
                shift # past value
            fi
            ;;
        -dir)
            if [[ -n "$2" ]] && [[ ! "$2" =~ ^- ]]; then
                benchmark_dir="$2"
                shift # past argument
                shift # past value
            fi
            ;;
        -func)
            if [[ -n "$2" ]] && [[ ! "$2" =~ ^- ]]; then
                benchmark_func="$2"
                shift # past argument
                shift # past value
            fi
            ;;
        -help)
            print_help
            ;;
        *)    # unknown option
            echo "Unknown option: $1" >&2
            print_help
            exit 1
            ;;
    esac
done

# for detecting manually quitting the script
quit=n
trap 'quit=y' INT


# temporarily stash local changes
has_local_changes=0
if [[ `git status --porcelain` ]]; then
    has_local_changes=1
    echo "stashing local changes"
    git stash --quiet
fi


# prepare output files
mkdir -p tmp
rm -f tmp/current.txt
rm -f tmp/target.txt


# run target benchmarks
git checkout ${target_commit_hash} --quiet
echo "running benchmarks for target ..."
go test -count=${number_of_runs} -bench=${benchmark_func} ${benchmark_dir} -timeout 1h > tmp/target.txt
if [ $? -ne 0 ]; then
    if [ x"$quit" != xy ]; then
        # did not manually stop the script
        echo "ERROR: benchmarks failed. See /tmp/target.txt for more details"
    fi
    
    git switch - --quiet
    if [ $has_local_changes -eq 1 ]; then
        git stash pop --quiet
    fi
    
    exit 1
fi


# revert temporary changes
git switch - --quiet
if [ $has_local_changes -eq 1 ]; then
    git stash pop --quiet
fi


# run current benchmarks
echo "running benchmarks for current ..."
go test -count=${number_of_runs} -bench=${benchmark_func} ${benchmark_dir} -timeout 1h > tmp/current.txt
if [ $? -ne 0 ]; then
    if [ x"$quit" != xy ]; then
        # did not manually stop the script
        echo "ERROR: benchmarks failed. See /tmp/current.txt for more details"
    fi
    
    exit 1
fi


# compare benchmarks
echo ""
go tool benchstat tmp/target.txt tmp/current.txt
