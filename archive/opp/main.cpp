#include "precomp.h"
#include "opp.h"

int main( int argc, char* argv[] )
{
    if( argc == 3 )
    {
        OPP opp;
        opp.ProcessFile(argv[1],argv[2]);
    }
    else printf( "USAGE: OPP <Source> <Target>" );
	return 0;
} // main()




