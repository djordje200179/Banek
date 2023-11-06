# Banek

The programming language of the future. It is a dynamically typed language
with a syntax mixing the best of Rust, Go, JS, Python and C.

The language is still in development and is not ready for production use.
It was made as a hobby project, and it will be updated whenever I have time
and motivation to add new features.

## Toolchain

Every part of the toolchain is written in Go. So my language depends on
Go runtime and its standard library.

Lexing and parsing is done completely by hand. I did not use any external
libraries for this, because I wanted to learn how it works. Both the lexer
and the parser run on their own goroutines, so they can work in parallel.
They work in a streaming fashion, producing tokens and sending them to the
consumers via buffered channels.

Parsing is done using a recursive descent Pratt parser. It is a top-down
parser that uses operator precedence to parse expressions. It is very
simple to implement, and it is also very fast.

### Interpreter
The interpreter is used for running Banek programs without compiling them
or for running the REPL. 
Executing a program with the interpreter is slower than executing the bytecode 
because it recursively evaluates the AST and parses nodes on every step.

### Compiler
The compiler is used for compiling Banek programs to bytecode. The bytecode
can then be run by the emulator. It does not perform any optimizations yet.

### Emulator
The emulator is used for running Banek bytecode. It is pretty fast, but
the bottleneck is the Go runtime.

### Disassembler
The disassembler is used for disassembling Banek bytecode. It is useful
for inspecting the generated bytecode.


## Features

### Variables

#### Types
Banek is a dynamically typed language, meaning that the type of variable
is not known until runtime. The following types are supported:
- booleans
- integers
- strings
- arrays
- functions
- `undefined` constant 

```banek
let num = 1;
let s = "Hello, world!";
let arr = [1, 2, 3];
let fun = || -> 1;
```

#### Mutability
Variables are by default immutable. This means that they cannot be
reassigned after they are declared. To make a variable mutable, the
keyword `mut` must be used before the variable name.

```banek
let a = 1;
let mut b = 2;

a = 3; // Error
b = 4; // OK
```

#### Visibility
Variables have block scope. This means that variables declared inside a block
are not accessible outside of that block.

```banek
let a = 1;
{
    let b = 2;
}
// b is not accessible here
```

### Functions

#### Built-in functions
Few functions are built-in to the language. These are needed for
interacting with the outside world and for more complex operations.

- `print` - Prints given values to the standard output
- `println` - Prints given values to the standard output and adds a newline
- `read` - Reads a value from the standard input and returns it as a string
- `readln` - Reads a line from the standard input and returns it as a string
- `len` - Returns the length of the given array
- `str` - Converts the given value to a string
- `int` - Converts the given value to an integer

```banek
print("Hello, world!");
print(len([1, 2, 3])); // Prints 3
```

#### Visibility

Functions are declared either by statements or as lambda expressions.
When declared as a statement, keyword `function` and the name of the function
is used. When declared as a lambda expression, the keyword `fn` is used.

```banek
func foo() {
    return 1;
}

let doubler = |num| -> num * 2;
```

#### Closure

Functions can access and capture variables from the outer scope. 

```banek
func adder(increment) {
    return |num| -> num + increment;
}

let addOne = adder(1);
print(addOne(1)); // Prints 2
```


#### Values

Functions are first-class citizens, meaning that they can be passed as
any other value.

```banek
func getter() {
    return 1;
}

func foo(getter) {
    print(getter());
}
```

#### Arguments

Functions can take any number of arguments. Arguments are declared inside
parentheses after the function name.

If a function is called with too few arguments, the missing arguments are
set to `undefined`. If a function is called with too many arguments, the extra
arguments are ignored.

```banek
func foo(a, b) {
    return a + b;
}
```

### Operators

#### Arithmetic operators
All the basic arithmetic operators are supported. They not only work 
with numbers, but some of them also work with strings and arrays.

- `+` - Addition or concatenation
- `-` - Subtraction
- `*` - Multiplication
- `/` - Division
- `%` - Remainder
- `^` - Exponentiation

```banek
print("Hello, " + "world!");    // Prints "Hello, world!"
print([1, 2] + [3, 4]);         // Prints [1, 2, 3, 4]
print(5 ^ 2);                   // Prints 25
```

#### Comparison operators
Comparison operators are used to compare two values. They return a boolean
value depending on the result of the comparison.

Every value can be compared to every other value for equality. 
However, only integers and strings can be compared for ordering.

```banek
print(1 == 2);      // Prints false
print(1 != "1");    // Prints true
print(1 < 2);       // Prints true
```

### Control flow

#### If statements
If statements are used to execute a block of code conditionally, like
in most other languages. 

Brackets are not needed around the condition. Bodies of consequent and
alternative branches do not need to be blocks if they consist of only one statement.

```banek
if 1 == 2 then
    print("1 is equal to 2");
else
    print("1 is not equal to 2");
```

Also, these statements can be used as expressions. In this case, the
alternative branch is mandatory, and the value of the expression is
the value of the executed branch.

```banek
let a = if 1 == 2 then "same" else "different";
```

#### While loops
While loops are used to execute a block of code repeatedly until a
condition is met. 

Like if statements, brackets are not needed around the condition and the
body does not need to be a block if it consists of only one statement.

```banek
let mut i = 0;
while i < 10 do {
    print(i);
    i = i + 1;
}
```

