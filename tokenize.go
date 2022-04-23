/*
** 2001 September 15
**
** The author disclaims copyright to this source code.  In place of
** a legal notice, here is a blessing:
**
**    May you do good and not evil.
**    May you find forgiveness for yourself and forgive others.
**    May you share freely, never taking more than you give.
**
*************************************************************************
** An tokenizer for SQL
**
** This file contains C code that splits an SQL input string up into
** individual tokens and sends those tokens one-by-one over to the
** parser for analysis.
 */
package internal

/* Character classes for tokenizing
**
** In the sqlite3GetToken() function, a switch() on aiClass[c] is implemented
** using a lookup table, whereas a switch() directly on c uses a binary search.
** The lookup table is much faster.  To maximize speed, and to ensure that
** a lookup table is used, all of the classes need to be small integers and
** all of them need to be used within the switch.
 */
const (
	CC_X        = 0  /* The letter 'x', or start of BLOB literal */
	CC_KYWD0    = 1  /* First letter of a keyword */
	CC_KYWD     = 2  /* Alphabetics or '_'.  Usable in a keyword */
	CC_DIGIT    = 3  /* Digits */
	CC_DOLLAR   = 4  /* '$' */
	CC_VARALPHA = 5  /* '@', '#', ':'.  Alphabetic SQL variables */
	CC_VARNUM   = 6  /* '?'.  Numeric SQL variables */
	CC_SPACE    = 7  /* Space characters */
	CC_QUOTE    = 8  /* '"', '\'', or '`'.  String literals, quoted ids */
	CC_QUOTE2   = 9  /* '['.   [...] style quoted ids */
	CC_PIPE     = 10 /* '|'.   Bitwise OR or concatenate */
	CC_MINUS    = 11 /* '-'.  Minus or SQL-style comment */
	CC_LT       = 12 /* '<'.  Part of < or <= or <> */
	CC_GT       = 13 /* '>'.  Part of > or >= */
	CC_EQ       = 14 /* '='.  Part of = or == */
	CC_BANG     = 15 /* '!'.  Part of != */
	CC_SLASH    = 16 /* '/'.  / or c-style comment */
	CC_LP       = 17 /* '(' */
	CC_RP       = 18 /* ')' */
	CC_SEMI     = 19 /* ';' */
	CC_PLUS     = 20 /* '+' */
	CC_STAR     = 21 /* '*' */
	CC_PERCENT  = 22 /* '%' */
	CC_COMMA    = 23 /* ',' */
	CC_AND      = 24 /* '&' */
	CC_TILDA    = 25 /* '~' */
	CC_DOT      = 26 /* '.' */
	CC_ID       = 27 /* unicode characters usable in IDs */
	CC_ILLEGAL  = 28 /* Illegal character */
	CC_NUL      = 29 /* 0x00 */
	CC_BOM      = 30 /* First byte of UTF8 BOM:  0xEF 0xBB 0xBF */
)

// SQLITE_ASCII
var aiClass = []rune{
	/*         x0  x1  x2  x3  x4  x5  x6  x7  x8  x9  xa  xb  xc  xd  xe  xf */
	/* 0x */ 29, 28, 28, 28, 28, 28, 28, 28, 28, 7, 7, 28, 7, 7, 28, 28,
	/* 1x */ 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	/* 2x */ 7, 15, 8, 5, 4, 22, 24, 8, 17, 18, 21, 20, 23, 11, 26, 16,
	/* 3x */ 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 19, 12, 14, 13, 6,
	/* 4x */ 5, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	/* 5x */ 1, 1, 1, 1, 1, 1, 1, 1, 0, 2, 2, 9, 28, 28, 28, 2,
	/* 6x */ 8, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	/* 7x */ 1, 1, 1, 1, 1, 1, 1, 1, 0, 2, 2, 28, 10, 28, 25, 28,
	/* 8x */ 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	/* 9x */ 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	/* Ax */ 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	/* Bx */ 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	/* Cx */ 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	/* Dx */ 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	/* Ex */ 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 30,
	/* Fx */ 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
}

// SQLITE_EBCDIC
var ebcdis_aiClass = []rune{
	/*         x0  x1  x2  x3  x4  x5  x6  x7  x8  x9  xa  xb  xc  xd  xe  xf */
	/* 0x */ 29, 28, 28, 28, 28, 7, 28, 28, 28, 28, 28, 28, 7, 7, 28, 28,
	/* 1x */ 28, 28, 28, 28, 28, 7, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	/* 2x */ 28, 28, 28, 28, 28, 7, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	/* 3x */ 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	/* 4x */ 7, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 26, 12, 17, 20, 10,
	/* 5x */ 24, 28, 28, 28, 28, 28, 28, 28, 28, 28, 15, 4, 21, 18, 19, 28,
	/* 6x */ 11, 16, 28, 28, 28, 28, 28, 28, 28, 28, 28, 23, 22, 2, 13, 6,
	/* 7x */ 28, 28, 28, 28, 28, 28, 28, 28, 28, 8, 5, 5, 5, 8, 14, 8,
	/* 8x */ 28, 1, 1, 1, 1, 1, 1, 1, 1, 1, 28, 28, 28, 28, 28, 28,
	/* 9x */ 28, 1, 1, 1, 1, 1, 1, 1, 1, 1, 28, 28, 28, 28, 28, 28,
	/* Ax */ 28, 25, 1, 1, 1, 1, 1, 0, 2, 2, 28, 28, 28, 28, 28, 28,
	/* Bx */ 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 9, 28, 28, 28, 28, 28,
	/* Cx */ 28, 1, 1, 1, 1, 1, 1, 1, 1, 1, 28, 28, 28, 28, 28, 28,
	/* Dx */ 28, 1, 1, 1, 1, 1, 1, 1, 1, 1, 28, 28, 28, 28, 28, 28,
	/* Ex */ 28, 28, 1, 1, 1, 1, 1, 0, 2, 2, 28, 28, 28, 28, 28, 28,
	/* Fx */ 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 28, 28, 28, 28, 28, 28,
}

/*
** The charMap() macro maps alphabetic characters (only) into their
** lower-case ASCII equivalent.  On ASCII machines, this is just
** an upper-to-lower case map.  On EBCDIC machines we also need
** to adjust the encoding.  The mapping is only valid for alphabetics
** which are the only characters for which this feature is used.
**
** Used by keywordhash.h
 */
//TODO: #ifdef SQLITE_ASCII
//TODO: # define charMap(X) sqlite3UpperToLower[(unsigned char)X]
//TODO: #endif
//TODO: #ifdef SQLITE_EBCDIC
//TODO: # define charMap(X) ebcdicToAscii[(unsigned char)X]
var ebcdicToAscii = []rune{
	/* 0   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F */
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, /* 0x */
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, /* 1x */
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, /* 2x */
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, /* 3x */
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, /* 4x */
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, /* 5x */
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 95, 0, 0, /* 6x */
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, /* 7x */
	0, 97, 98, 99, 100, 101, 102, 103, 104, 105, 0, 0, 0, 0, 0, 0, /* 8x */
	0, 106, 107, 108, 109, 110, 111, 112, 113, 114, 0, 0, 0, 0, 0, 0, /* 9x */
	0, 0, 115, 116, 117, 118, 119, 120, 121, 122, 0, 0, 0, 0, 0, 0, /* Ax */
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, /* Bx */
	0, 97, 98, 99, 100, 101, 102, 103, 104, 105, 0, 0, 0, 0, 0, 0, /* Cx */
	0, 106, 107, 108, 109, 110, 111, 112, 113, 114, 0, 0, 0, 0, 0, 0, /* Dx */
	0, 0, 115, 116, 117, 118, 119, 120, 121, 122, 0, 0, 0, 0, 0, 0, /* Ex */
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, /* Fx */
}

/*
** The sqlite3KeywordCode function looks up an identifier to determine if
** it is a keyword.  If it is a keyword, the token code of that keyword is
** returned.  If the input is not a keyword, TK_ID is returned.
**
** The implementation of this routine was generated by a program,
** mkkeywordhash.c, located in the tool subdirectory of the distribution.
** The output of the mkkeywordhash.c program is written into a file
** named keywordhash.h and then included into this source file by
** the #include below.
 */

/*
** If X is a character that can be used in an identifier then
** IdChar(X) will be true.  Otherwise it is false.
**
** For ASCII, any character with the high-order bit set is
** allowed in an identifier.  For 7-bit characters,
** sqlite3IsIdChar[X] must be 1.
**
** For EBCDIC, the rules are more complex but have the same
** end result.
**
** Ticket #1066.  the SQL standard does not allow '$' in the
** middle of identifiers.  But many SQL implementations do.
** SQLite will allow '$' in identifiers for compatibility.
** But the feature is undocumented.
 */
//TODO: #ifdef SQLITE_ASCII
//TODO: #define IdChar(C)  ((sqlite3CtypeMap[(unsigned char)C]&0x46)!=0)
//TODO: #endif
//TODO: #ifdef SQLITE_EBCDIC
var sqlite3IsEbcdicIdChar = []rune{
	/* x0 x1 x2 x3 x4 x5 x6 x7 x8 x9 xA xB xC xD xE xF */
	0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, /* 4x */
	0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 0, /* 5x */
	0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 0, 0, /* 6x */
	0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, /* 7x */
	0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, /* 8x */
	0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 1, 0, /* 9x */
	1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, /* Ax */
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, /* Bx */
	0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, /* Cx */
	0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, /* Dx */
	0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, /* Ex */
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 0, /* Fx */
}

//TODO: #define IdChar(C)  (((c=C)>=0x42 && sqlite3IsEbcdicIdChar[c-0x40]))
//TODO: Placeholder function to compile
func IdChar(c uint8) int {
	return 0
}

/* Make the IdChar function accessible from ctime.c and alter.c */
func sqlite3IsIdChar(c uint8) int { return IdChar(c) }

// TODO:
// #ifndef SQLITE_OMIT_WINDOWFUNC
// /*
// ** Return the id of the next token in string (*pz). Before returning, set
// ** (*pz) to point to the byte following the parsed token.
// */
// static int getToken(const unsigned char **pz){
//   const unsigned char *z = *pz;
//   int t;                          /* Token type to return */
//   do {
//     z += sqlite3GetToken(z, &t);
//   }while( t==TK_SPACE );
//   if( t==TK_ID
//    || t==TK_STRING
//    || t==TK_JOIN_KW
//    || t==TK_WINDOW
//    || t==TK_OVER
//    || sqlite3ParserFallback(t)==TK_ID
//   ){
//     t = TK_ID;
//   }
//   *pz = z;
//   return t;
// }
//
// /*
// ** The following three functions are called immediately after the tokenizer
// ** reads the keywords WINDOW, OVER and FILTER, respectively, to determine
// ** whether the token should be treated as a keyword or an SQL identifier.
// ** This cannot be handled by the usual lemon %fallback method, due to
// ** the ambiguity in some constructions. e.g.
// **
// **   SELECT sum(x) OVER ...
// **
// ** In the above, "OVER" might be a keyword, or it might be an alias for the
// ** sum(x) expression. If a "%fallback ID OVER" directive were added to
// ** grammar, then SQLite would always treat "OVER" as an alias, making it
// ** impossible to call a window-function without a FILTER clause.
// **
// ** WINDOW is treated as a keyword if:
// **
// **   * the following token is an identifier, or a keyword that can fallback
// **     to being an identifier, and
// **   * the token after than one is TK_AS.
// **
// ** OVER is a keyword if:
// **
// **   * the previous token was TK_RP, and
// **   * the next token is either TK_LP or an identifier.
// **
// ** FILTER is a keyword if:
// **
// **   * the previous token was TK_RP, and
// **   * the next token is TK_LP.
// */
// static int analyzeWindowKeyword(const unsigned char *z){
//   int t;
//   t = getToken(&z);
//   if( t!=TK_ID ) return TK_ID;
//   t = getToken(&z);
//   if( t!=TK_AS ) return TK_ID;
//   return TK_WINDOW;
// }
// static int analyzeOverKeyword(const unsigned char *z, int lastToken){
//   if( lastToken==TK_RP ){
//     int t = getToken(&z);
//     if( t==TK_LP || t==TK_ID ) return TK_OVER;
//   }
//   return TK_ID;
// }
// static int analyzeFilterKeyword(const unsigned char *z, int lastToken){
//   if( lastToken==TK_RP && getToken(&z)==TK_LP ){
//     return TK_FILTER;
//   }
//   return TK_ID;
// }
// #endif /* SQLITE_OMIT_WINDOWFUNC */
