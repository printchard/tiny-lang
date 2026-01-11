# Tiny Lang

Tiny lang is a small, dynamically-typed language created for recreational purposes. It is fully implemented in Go.

## Specification

You can find the language formal definition in the `spec.bnf` file in EBNF syntax.

## Capabilities

### Types

- **Number**
- **String**
- **Boolean**
- **Array**
- **Function**
- **Void**

### Language Features

- **Variables** - Declare with `let`, assign with `=`
- **Functions** - Define with `func name: arg1, arg2 { }` syntax, call with `name(args)`
- **Control Flow** - `if`/`else` statements and `while` loops
- **Return Statements** - Early return from functions with `return` or `return value`
- **Operators**:
  - Arithmetic: `+`, `-`, `*`, `/`
  - Comparison: `==`, `!=`, `<`, `<=`, `>`, `>=`
  - Logical: `&&`, `||`, `!`
  - Array indexing: `arr[index]`
- **Built-in Functions**:
  - `print(value, ...)` - Print values to stdout

### Example

```tiny
func greet: name {
  if name {
    print("Hello " + name + "!")
    return true
  }
  print("Hello World!")
  return false
}

let result := greet("Alice")
print(result)

let numbers := [1, 2, 3, 4, 5]
let i := 0
while i < 5 {
  print(numbers[i])
  i = i + 1
}
```

## Using the language

You can use the language by downloading one of the binaries in the releases on GitHub. If you have go installed, you can install the binary doing `go install github.com/printchard/tiny-lang` or cloning the repository and building form source.

Once you have the binary, you can enter REPL mode by running the binary, or interpret a file if you provide the filename as a CLI argument. The default extension for the language is `.tiny`.

### Building for Multiple Platforms

Run the included build script to create binaries for all supported platforms:

```bash
./release.sh
```

This will generate binaries in the `build/` directory for:

- macOS (amd64, arm64)
- Linux (amd64, arm64)
- Windows (amd64, arm64)
