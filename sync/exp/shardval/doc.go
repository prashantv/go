// Package shardval provides an implementation of a sharded value.
//
// This is useful for performance-sensitive code that is running across many cores
// where contention can cause significant slowdowns.
// See benchmarks for more details.
package shardval
