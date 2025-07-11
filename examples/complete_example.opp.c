// Complete OPP Example demonstrating all features
#include <stdio.h>

##:DEBUG_LOG(msg) printf("[DEBUG %d] " msg "\n", ##_)
##:MAX(a,b) ((a)>(b)?(a):(b))
##:§ MAX

##~(~DEBUG|~DEBUG)|~(~DEBUG|~DEBUG)
// Running in DEBUG mode
DEBUG_LOG("Program started");
##.

int main() {
    ##~WINDOWS|~WINDOWS
    printf("Running on non-Windows system\n");
    ##@~(~WINDOWS|~WINDOWS)|~(~WINDOWS|~WINDOWS)  
    printf("Running on Windows\n");
    ##.
    
    // Using unicode macro name
    int x = §(10, 20);
    printf("Max of 10 and 20 is %d\n", x);
    
    // Predefined macros
    printf("Line number minus 5: %d\n", ##_);
    printf("Random number: %d\n", ##$);
    
    { // Brace counting
        printf("Opening braces so far: %d\n", ##{);
        {
            printf("More braces: %d\n", ##{);
        }
    }
    
    printf("Closing braces mod 5: %d\n", ##});
    
    return 0;
}

##~(~INCLUDE_FOOTER|~INCLUDE_FOOTER)|~(~INCLUDE_FOOTER|~INCLUDE_FOOTER)
##<footer\.h.
##.