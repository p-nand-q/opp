# Archive Directory

This directory contains the original C++ implementation of OPP (Obfuscated Pre-Processor) for historical reference.

## Contents

- `opp/` - Original C++ source code from p-nand-q.com
  - `OPP.cpp`, `OPP.h` - Main preprocessor implementation
  - `main.cpp` - Command-line interface
  - `gtools.cpp`, `gtools.h` - Utility functions
  - `precomp.cpp`, `precomp.h` - Precompiled headers
  - `OPP.dsp`, `OPP.dsw` - Visual Studio project files
  - `OPP.exe` - Original Windows executable
- `DIFFERENCES.md` - Detailed comparison between specification and original implementation

## Notes

The original implementation contains several bugs that differ from the specification:
- Uses AND logic instead of NAND
- Uses `.` as separator instead of `|`
- Both `{` and `}` increment the open brace counter
- Missing implementation of the undefine (`-`) directive

The Go implementation in this repository follows the specification as documented, not the buggy behavior of this original code.

## Historical Value

These files are preserved to:
1. Document the evolution of the language
2. Allow comparison between original and new implementation
3. Provide context for design decisions in the Go version