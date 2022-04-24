package internal

// No-op functions that are currently used from parse.y
func sqlite3ExprDelete(db *sqlite3, p *Expr) {}

func sqlite3ExprListDelete(db *sqlite3, p *ExprList) {}

func sqlite3SelectDelete(db *sqlite3, p *Select) {}

func sqlite3SrcListDelete(db *sqlite3, p *SrcList) {}

func sqlite3WithDelete(db *sqlite3, p *With) {}

func sqlite3WindowListDelete(db *sqlite3, p *Window) {}

func sqlite3WindowDelete(db *sqlite3, p *Window) {}

func sqlite3DeleteTriggerStep(db *sqlite3, p *TriggerStep) {}

func sqlite3IdListDelete(db *sqlite3, p *IdList) {}

func sqlite3ErrorMsg(p *Parse, fmt string, args ...interface{}) {}

func sqlite3FinishCoding(*Parse) {}

func sqlite3BeginTransaction(*Parse, uint16) {}

func sqlite3EndTransaction(*Parse, uint16) {}

func sqlite3Savepoint(*Parse, int, *Token) {}

func sqlite3StartTable(*Parse, *Token, *Token, int, int, int, int) {}

func sqlite3EndTable(*Parse, *Token, *Token, uint32, *Select) {}

func sqlite3_strnicmp(textA, textB string, length int) int {
	return 0
}

func sqlite3AddColumn(*Parse, Token, Token) {}

func sqlite3AddDefaultValue(*Parse, *Expr, string, string) {}

// ./parse.go:2628:3: undefined: sqlite3ErrorMsg
