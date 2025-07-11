# OPP - Obfuscated Pre-Processor

This page describes OPP, a free open-source preprocessor that you can add to all languages you want to use it in. It is used in the new editions of my languages Sorted!, Smith#, and Java2K.

## Quick Start

### Installation

```bash
go install github.com/p-nand-q/opp/cmd/opp@latest
```

### Basic Usage

```bash
# Process a file
opp input.opp -o output.go

# Define variables from command line
opp -D DEBUG=1 -D VERSION=2 input.opp

# Output to stdout
opp input.opp
```

### Example

```go
##~DEBUG|~DEBUG
// This line only appears if DEBUG is undefined
##.
##~(~DEBUG|~DEBUG)|~(~DEBUG|~DEBUG)
// This line only appears if DEBUG is defined
##.
```

## Go Implementation

This is a Go implementation of OPP. Build from source:

```bash
git clone https://github.com/p-nand-q/opp
cd opp
go build -o opp ./cmd/opp
```

## Conditional Compilation

What would a preprocessor be without conditional compilation? In OPP, you have the power of NAND to do all conditional compilation. Use

```
##~a|~b
```

(read: not a or not b) to compile the following lines if a is undefined or b is undefined. If you want to check if only a is undefined, use

```
##~a|~a
```

If you want to check if a is defined, use

```
##~(~a|~a)|~(~a|~a)
```

The statement

```
##.
```

ends a conditional compilation block. There is no #else in OPP, but you can use

```
##@
```

to do an if-else. For example, if you have a section A to be compiled if B is defined, and a section C to be compiled otherwise, you'd write something like

```
##~A|~A
..B..
##@~(~A|~A)|~(~A|~A)
..C..
##.
```

Note that you cannot use whitespaces in the expressions or between ## and ~. Who needs code indentation, anyway? All of this makes the (recursive) syntax so simple its a joke:

```
SYNTAX = '~' OBJECT '|~' OBJECT.
OBJECT = VARIABLE | '(' SYNTAX ')'.
```

Variables can be specified either as environment variables or explicitly as macros (see below).

**Implementation Note**: The original C++ implementation has several bugs and differences:
- Uses `.` as separator instead of `|` (so `##~a.~b` instead of `##~a|~b`)
- Implements AND logic instead of NAND due to a bug
- Only checks environment variables, not defined macros in conditionals
- Our implementation follows the specification as documented

## Includes

You can use OPP to include other sourcefiles in your code, by using the following statement

```
##<<Filename>.
```

**Important**: The `.` at the end is mandatory - it terminates the filename. This resolves the ambiguity between file inclusion (`##<filename.`) and other directives that might start with `##<`.

Because filenames with dots are frequent, you can use the escape sequence `\.` to specify a single dot. Because filenames with `\` are frequent on certain OSs (OSsi?) you can use `..` to specify a single `\`. Because filenames with `..` are frequent on certain OSs, you can use `\\` to specify a single `..`. Because filenames with `\\` are not so frequent on the OS in question, but still possible, you can use `//` to specify a single `\\`, and if your code editor has syntax coloring, you get the rest of the line in comment color, without paying extra. Here is a sample filename in C (on Windows NT) and its equivalent on OPP. Judge for yourself which is more phon - er - fun.

```c
#include "\\server\users\opp\sample.h"
```
```
##<//server..users..opp..sample\.h.
```

Of course, you must specify absolute paths, because OPP does not support proprietary environment variables such as "INCLUDE" or "PATH".

**Implementation Note**: The original adds a newline before included content.

## Defining Macros

OPP supports defining function macros, by utilizing the following syntax.

```
##:<name of macro, followed by a single blank> <macro definition>
```

Everything after the single blank is treated as part of the macro body. Macro arguments are implicitly defined by using #0 to refer to the first argument, #1 to the second and so on. Macronames can be virtually anything, including esoteric characters like ? or ä. Here is a macro, that evaluates the max of two arguments as a function "§".

```
##:§ ((#0<#1)?#1:#0)
```

Now you can write §(a,b) everywhere in your code without unsuspecting people knowing what the deal is.

You can also include varargs in macros, as in the following example, which is a solution to the age-old problem bothering the C macro language: conditional compilation of printf. The syntax for the vararg sequence is

```
##<from>..n
```

where `<from>` specifies the first argument.

```
##~(~DEBUG|~DEBUG)|~(~DEBUG|~DEBUG)
##:dbg printf(##0..n)
##@~DEBUG|~DEBUG
##:dbg
##.
```

In debug builds, you can use dbg just like normal printf(), in release builds they'll be omitted. To undefine a macro (or an operator, see below), use

```
##-<name of macro or operator>
```

## Predefined macros

The following macros are predefined

- `##i` - the square root of -1 (you need to include complex.h to use this macro)
- `##_` - the current line number minus 5
- `##$` - a pseudo-random number
- `##{` - The number of { in the code up to this point
- `##}` - The number of } in the code up to this point, modulo 5

**Implementation Notes**:
- The original returns `complex(0,1)` for `##i` instead of `1i`
- The original has a bug where both `{` and `}` increment the open brace counter

## Macros inside Macros

In OPP, you can declare macros in other macros. To pass macro arguments or directives literally to nested macro definitions, use the escape sequence `##,#` (a 4-character sequence). This outputs a literal `#` or `##` that won't be processed immediately.

The escape works as follows:
- `##,#0` → `#0` (literal, not expanded)
- `##,##` → `##` (literal, not processed)  
- `##,#1` → `#1` (literal, not expanded)

### Example: Macro that Defines Macros

```
##:DEFINE_GETTER(name) ##:get_##,#0() { return ##,#0; }
DEFINE_GETTER(width)
DEFINE_GETTER(height)
```

This expands to:
```
##:get_width() { return width; }
##:get_height() { return height; }
```

Which then defines two new macros `get_width` and `get_height`.

### Example: Macro with Nested Argument Forwarding

```
##:WRAPPER(fn) ##:safe_##,#0(x) { if (x != NULL) ##,#0(x); }
WRAPPER(free)
WRAPPER(close)
```

Expands to macros that safely call functions only if the argument is non-NULL.

## Preprocessor operators

In macro or operator definitions, macro arguments are specified by #, followed by a zero-based index. You can use the following macro argument operators.

- `#"` - stringize operand. Do not confuse with ##"
- `#'` - charize operand. Do confuse with ##'

Each of these applies to the next argument expanded, and ONLY to the next argument expanded. You can specify them anywhere in a macro definition.

## Working with Lines

You cannot span macros across multiple lines. Use tense code!

## Using OPP in your own programs

You can include the OPP class in your own programs. This Go implementation provides:

```go
import "github.com/p-nand-q/opp"

preprocessor := opp.New()
preprocessor.Define("DEBUG", "1")
output, err := preprocessor.Process(input)
```

Alternatively, you can check out my other programming languages each of which prominently features OPP.

## Known Limitations

The following features from the OPP specification are not yet implemented:
- Function-like macros with arguments (`#0`, `#1`, `##0..n`)
- Stringize (`#"`) and charize (`#'`) operators

See the issue tracker for progress on these features.

## Example Files

The `examples/` directory contains several OPP example files:

### hello.opp.go
A Go program demonstrating basic OPP features including conditional compilation, debug macros, and platform-specific code.

### justif.opp.justif
Shows how to use OPP with the Justif esoteric language, including macro definitions for common operations.

### complete_example.opp.c
A comprehensive C example showcasing all OPP features:
- Unicode macro names (§ for MAX)
- Conditional compilation with NAND logic
- Predefined macros (##_, ##$, ##{, ##})
- Nested conditionals
- Include directives

These examples demonstrate how OPP can obfuscate any language while technically remaining a functional preprocessor.