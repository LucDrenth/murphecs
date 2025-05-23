#!/bin/bash

####################################
##          HOW TO USE            ##
####################################
#
# arg 1: (required) commit hash to compare to
# arg 2: (required) number of runs
# arg 3: (required) benchmarks directory
# arg 4: (optional) benchmarks function
#
# for example: compare-benchmarks.sh d2331b8c11083e15147b4470275c18896d7404e3 10 ./benchmark/
#
#####################################

# process args
if [ -z "$1" ]; then
    echo "ERROR: argument 1 (target_commit_hash) is invalid"
    exit 1
fi
target_commit_hash=$1

if [ -z "$2" ]; then
    echo "ERROR: argument 2 (number_of_runs) is invalid"
    exit 1
fi
number_of_runs=$2

if [ -z "$3" ]; then
    echo "ERROR: argument 3 (benchmark_dir) is invalid"
    exit 1
fi
benchmark_dir=$3

if [ -z "$4" ]; then
    benchmark_func='.'
else
    benchmark_func=$4
fi


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
