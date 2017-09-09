#!/usr/bin/env bats

# TODO: Auto generate this? Do so in topological order to pinpoint breaking libraries.

@test "Install gmp" {
	gogurt gmp
}

@test "Install mpfr" {
	gogurt mpfr
}

@test "Install mpc" {
	gogurt mpc
}

@test "Install gcc" {
	gogurt gcc
}
