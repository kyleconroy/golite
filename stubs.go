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

func sqlite3_strnicmp(textA, textB []byte, length int) int {
	return 0
}

func sqlite3AddColumn(*Parse, Token, Token) {}

func sqlite3AddDefaultValue(*Parse, *Expr, []byte, []byte) {}

func sqlite3PExpr(*Parse, int, *Expr, *Expr) *Expr {
	return nil
}

func sqlite3ExprIdToTrueFalse(*Expr) int {
	return 0
}

func sqlite3ExprTruthValue(*Expr) bool {
	return false
}

func sqlite3AddNotNull(*Parse, int) {}

func sqlite3AddPrimaryKey(*Parse, *ExprList, int, int, int) {}

func sqlite3CreateIndex(*Parse, *Token, *Token, *SrcList, *ExprList, int, *Token,
	*Expr, int, int, uint8) {
}

func sqlite3AddCheckConstraint(*Parse, *Expr, []byte, []byte) {}

func sqlite3CreateForeignKey(*Parse, *ExprList, *Token, *ExprList, int) {}

func sqlite3DeferForeignKey(*Parse, int)

func sqlite3AddCollateType(*Parse, *Token)

func sqlite3AddGenerated(*Parse, *Expr, *Token)

func sqlite3DropTable(*Parse, *SrcList, int, int)

func sqlite3CreateView(*Parse, *Token, *Token, *Token, *ExprList, *Select, int, int)

func sqlite3Select(*Parse, *Select, *SelectDest) int

func sqlite3SrcListAppendFromTerm(*Parse, *SrcList, *Token, *Token, *Token, *Select, *OnOrUsing) *SrcList {
	return nil
}

func sqlite3SelectNew(*Parse, *ExprList, *SrcList, *Expr, *ExprList,
	*Expr, *ExprList, uint32, *Expr) *Select {
	return nil
}

func sqlite3ExprListAppend(*Parse, *ExprList, *Expr) *ExprList {
	return nil
}

func sqlite3ExprListSetName(*Parse, *ExprList, *Token, int)

func sqlite3ExprListSetSpan(*Parse, *ExprList, []byte, []byte)

func sqlite3Expr(*sqlite3, int, []byte) *Expr {
	return nil
}
