## Current Status

```
go build
# github.com/kyleconroy/golite
./parse.go:3735:109: invalid operation: operator - not defined on yypParser.yystack[yypParser.yytos + 0].minor.yy0.z (variable of type []byte)
./parse.go:3938:109: cannot use 0 (untyped int constant) as []byte value in assignment
./parse.go:3953:44: cannot use yypParser.yystack[yypParser.yytos + -2].minor.yy614 (variable of type *ExprList) as type int in argument to sqlite3CreateIndex
./parse.go:3953:95: cannot use yypParser.yystack[yypParser.yytos + 0].minor.yy394 (variable of type int) as type *Token in argument to sqlite3CreateIndex
./parse.go:4079:54: cannot use yypParser.yystack[yypParser.yytos + 0].major (variable of type uint16) as type int in assignment
./parse.go:4090: cannot use yypParser.yystack[yypParser.yytos + -7].minor.yy394 (variable of type int) as type uint32 in argument to sqlite3SelectNew
./parse.go:4097: cannot use yypParser.yystack[yypParser.yytos + -8].minor.yy394 (variable of type int) as type uint32 in argument to sqlite3SelectNew
./parse.go:4116:20: assignment mismatch: 2 variables but 1 value
./parse.go:4153:54: cannot use 0 (untyped int constant) as *ExprList value in assignment
./parse.go:4170:49: cannot use 0 (untyped int constant) as []byte value in argument to sqlite3Expr
./parse.go:4170:49: too many errors
make: *** [build] Error 2
```

- File src/parse.y artifact b86d56b4 on branch trunk
- File src/tokenize.c artifact a38f5205 on branch trunk
- File src/sqliteInt.h artifact 36b5d1cc on branch trunk

