# Tiny Lang

Tiny lang is a small language created for recreational purposes. It is fully implemented in Go.

## Specification

You can find the language formal definition in the `spec.bnf` file in ebnf syntax.

## Capabilities

Currently, the language supports the following types:

- Float64
- Strings
- Booleans
- Arrays

Additionally, the language supports the following constructs:

- Print
- Arithmetic Operations
- Boolean Operations
- If statements
- While statements

## Using the language

You can use the language by cloning the repository and building it using `go build .`. If you have go installed, you can install the binary doing `go install github.com/printchard/tiny-lang`.

Once you have the binary, you can enter REPL mode by running the binary, or interpret a file if you provide the filename as a CLI argument. The default extension for the language is `.tiny`.
