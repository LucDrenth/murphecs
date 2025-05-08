package ecs

type Query interface {
	Exec(world *World)

	// Prepare extracts the query options and puts it in CombinedQueryOptions. This should be called
	// once, after which the query is ready to be used (e.g. Exec can be called).
	//
	// This step is necessary because the way that queries are created is optimized for users, and not
	// for the computer. This method closes that gap.
	Prepare() error

	// Validate checks if there are any unexpected or unoptimized things. It returns an error if there
	// is something that can be optimized. The error should be treated as a warning, and does not mean
	// that the query can not be executed.
	//
	// Prepare must be called without returning any errors before calling this method.
	Validate() error

	// Clear the query results that got filled when last running Exec.
	ClearResults()

	// IsLazy returns wether this query should be treated as lazy or not. Being lazy means that it should
	// not get executed automatically, and should be done by the user.
	IsLazy() bool
}
