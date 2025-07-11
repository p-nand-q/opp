# Proposed Specification for Stringize and Charize Operators

## Overview

The stringize (`#"`) and charize (`#'`) operators convert macro arguments into string or character literals during macro expansion.

## Syntax

### Stringize Operator (`#"`)
- **Syntax**: `#"#n` where `n` is the argument number (0-based)
- **Effect**: Wraps the expanded argument in double quotes
- **Whitespace**: No space allowed between `#"` and `#n`

### Charize Operator (`#'`) 
- **Syntax**: `#'#n` where `n` is the argument number (0-based)
- **Effect**: Wraps the expanded argument in single quotes
- **Whitespace**: No space allowed between `#'` and `#n`

## Behavior Rules

1. **Immediate Application**: The operator must immediately precede the argument reference with no intervening characters.
   - ✅ Valid: `#"#0`, `#'#1`
   - ❌ Invalid: `#" #0`, `#' #1`

2. **Single Argument Only**: Each operator applies to exactly one argument reference.
   - `#"#0#1` → `"expanded_arg0"#1` (only #0 is stringized)

3. **Escape Handling**: 
   - For stringize: Internal `"` characters are escaped as `\"`
   - For charize: Internal `'` characters are escaped as `\'`
   - Backslashes are escaped as `\\`

4. **Order of Operations**:
   1. Argument substitution happens first
   2. Stringize/charize wrapping happens second
   3. Further macro expansion happens third

5. **No Nesting**: These operators cannot be nested or combined.
   - ❌ Invalid: `#"#'#0` 

## Examples

### Basic Stringize
```
##:STR(x) #"#0
STR(hello) → "hello"
STR(hello world) → "hello world"
STR(test"quote) → "test\"quote"
```

### Basic Charize
```
##:CHR(x) #'#0
CHR(a) → 'a'
CHR(hello) → 'hello'
CHR(it's) → 'it\'s'
```

### Mixed Usage
```
##:LOG(func,msg) printf("%s: %s\n", #"#0, #"#1)
LOG(main, Starting program) → printf("%s: %s\n", "main", "Starting program")
```

### With Expressions
```
##:ASSERT(expr) if (!(#0)) error(#"#0)
ASSERT(x > 0) → if (!(x > 0)) error("x > 0")
```

### Multiple Arguments
```
##:PAIR(a,b) {#"#0, #'#1}
PAIR(name, c) → {"name", 'c'}
```

## Edge Cases

### Empty Arguments
```
##:STR(x) #"#0
STR() → ""
```

### Numeric Arguments
```
##:STR(x) #"#0
STR(123) → "123"
STR(3.14) → "3.14"
```

### Arguments with Operators
```
##:STR(x) #"#0
STR(2+2) → "2+2"
```

### Varargs Interaction
```
##:ALL(args) #"##0..n
ALL(a,b,c) → "a,b,c"  // Entire varargs list is stringized as one unit
```

## What These Operators Do NOT Do

1. **No recursive stringizing**: If the argument contains `#"`, it's treated as literal text
2. **No macro expansion before stringizing**: The argument is stringized as-is
3. **No concatenation**: `#"#0#"#1` does not produce `"arg0arg1"`

## Implementation Notes

- These operators are only valid within macro definitions
- They are processed during macro expansion, not as standalone directives
- Using them outside macro definitions is a syntax error
- The `##,#` escape sequence can be used to output literal `#"` or `#'` in nested macros

## Rationale

This specification:
1. Removes all ambiguity about placement and spacing
2. Provides clear examples for every use case
3. Defines edge case behavior explicitly
4. Maintains simplicity while being useful
5. Avoids conflicts with other OPP features