# Differences Between OPP Specification and Original Implementation

## 1. NAND Logic Separator
- **Spec**: `##~a|~b`
- **Actual**: `##~a.~b`
- The separator is `.` not `|`

## 2. NAND Logic Bug
The original has a logic error - it implements AND instead of NAND:
```cpp
return a && b;  // Wrong! Should be: return a || b;
```

## 3. Brace Counting Bug
Both `{` and `}` increment open brace counter due to copy-paste error:
```cpp
else if( c == '}' )
    m_nCurlyBracketsOpen++;  // Should be m_nCurlyBracketsClose++!
```

## 4. Variables from Environment Only
Conditional compilation only checks environment variables, not defined macros.

## 5. Complex Number Format
`##i` expands to `complex(0,1)` not `1i`

## 6. Additional Features
- `##,#` escape sequence for nested macros (undocumented)
- Include directive adds newline before included content
- No implementation of `-` (undefine) directive

## 7. Recursive Macro Expansion
The implementation re-processes lines after macro expansion (goto restart_processing)

## Our Implementation Choice
Our Go implementation follows the specification as documented, with the following decisions:
- Uses `|` as NAND separator (as specified)
- Implements correct NAND logic
- Fixes the brace counting bug
- Supports both environment variables and defined macros
- Implements the `-` (undefine) directive