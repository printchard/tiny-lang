# Tiny Lang

Tiny lang is a small, dynamically-typed programming language created for recreational purposes. It is implemented in Go and includes a lexer, recursive-descent parser, interpreter, and source formatter.

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

- **Variables** - Declare variables with `let name := value` and assign with `=`
- **Functions** - Define named functions with `func name: arg1, arg2 { }`, or create function values with `func: arg1, arg2 { }`
- **Function calls** - Call named or function-valued expressions with `name(args)`
- **Closures** - Functions capture their defining environment
- **Control flow** - `if`/`else if`/`else` statements and `while` loops
- **Return statements** - Return from functions with `return` or `return value`
- **Arrays** - Create arrays with `[value1, value2]`, access elements with `values[index]`, and assign indexed values with `values[index] = value`
- **Operators**:
  - Arithmetic: `+`, `-`, `*`, `/`
  - Comparison: `==`, `!=`, `<`, `<=`, `>`, `>=`
  - Logical: `&&`, `||`, `!`
  - Unary negation: `-value`
  - Array indexing: `values[index]`
- **Operator precedence** - Arithmetic, comparison, and logical operators are parsed and formatted according to their precedence and associativity
- **Runtime errors** - Lexer, parser, and runtime errors include source location information
- **Built-in Functions**:
  - `print(value, ...)` - Print values to stdout
  - `len(value)` - Return the length of an array or string
  - `push(array, value, ...)` - Return an array containing the original elements followed by the supplied values
  - `map(array, function)` - Return an array containing the function result for each element

### Source Formatter

The `formatter` package provides a canonical formatter for Tiny Lang source:

```go
formatted, err := formatter.Format(source)
```

The formatter normalizes indentation, spaces, operator formatting, arrays, function calls, blocks, and top-level function spacing. It preserves expression meaning by adding parentheses where required by operator precedence and associativity. Formatting is idempotent:

```text
Format(Format(source)) == Format(source)
```

### Testing

The project has tests at several levels:

- Lexer tokenization and lexer errors
- Parser acceptance and syntax errors
- Expression evaluation and interpreter behavior
- Golden end-to-end programs
- Formatter precedence and idempotency
- Semantic preservation between original and formatted programs

Run the complete test suite with:

```bash
go test ./...
```

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

Function values and built-ins can be used together:

```tiny
let numbers := [1, 2, 3]
let doubled := map(numbers, func: value {
  return value * 2
})

print(len(doubled))
print(push(doubled, 8))
```

## Using the language

You can use the language by downloading one of the binaries in the releases on GitHub. If you have Go installed, you can install the binary with `go install github.com/printchard/tiny-lang` or clone the repository and build from source:

```bash
go build -o tiny-lang .
```

Once you have the binary, you can enter REPL mode by running the binary, or interpret a file if you provide the filename as a CLI argument. The default extension for the language is `.tiny`.

The formatter is also available through the CLI. To write formatted source to standard output:

```bash
tiny-lang fmt program.tiny
```

To format a file in place, use the `-w` flag:

```bash
tiny-lang fmt -w program.tiny
```

### Building for Multiple Platforms

Run the included build script to create binaries for all supported platforms:

```bash
./release.sh
```

This will generate binaries in the `build/` directory for:

- macOS (amd64, arm64)
- Linux (amd64, arm64)
- Windows (amd64, arm64)
