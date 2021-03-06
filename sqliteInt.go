package internal

/*
** Allowed values for Table.tabFlags.
**
** TF_OOOHidden applies to tables or view that have hidden columns that are
** followed by non-hidden columns.  Example:  "CREATE VIRTUAL TABLE x USING
** vtab1(a HIDDEN, b);".  Since "b" is a non-hidden column but "a" is hidden,
** the TF_OOOHidden attribute would apply in this case.  Such tables require
** special handling during INSERT processing. The "OOO" means "Out Of Order".
**
** Constraints:
**
**         TF_HasVirtual == COLFLAG_VIRTUAL
**         TF_HasStored  == COLFLAG_STORED
**         TF_HasHidden  == COLFLAG_HIDDEN
 */
const (
	TF_Readonly      = 0x00000001 /* Read-only system table */
	TF_HasHidden     = 0x00000002 /* Has one or more hidden columns */
	TF_HasPrimaryKey = 0x00000004 /* Table has a primary key */
	TF_Autoincrement = 0x00000008 /* Integer primary key is autoincrement */
	TF_HasStat1      = 0x00000010 /* nRowLogEst set from sqlite_stat1 */
	TF_HasVirtual    = 0x00000020 /* Has one or more VIRTUAL columns */
	TF_HasStored     = 0x00000040 /* Has one or more STORED columns */
	TF_HasGenerated  = 0x00000060 /* Combo: HasVirtual + HasStored */
	TF_WithoutRowid  = 0x00000080 /* No rowid.  PRIMARY KEY is the key */
	TF_StatsUsed     = 0x00000100 /* Query planner decisions affected by
	 ** Index.aiRowLogEst[] values */
	TF_NoVisibleRowid = 0x00000200 /* No user-visible "rowid" column */
	TF_OOOHidden      = 0x00000400 /* Out-of-Order hidden columns */
	TF_HasNotNull     = 0x00000800 /* Contains NOT NULL constraints */
	TF_Shadow         = 0x00001000 /* True for a shadow table */
	TF_HasStat4       = 0x00002000 /* STAT4 info available for this table */
	TF_Ephemeral      = 0x00004000 /* An ephemeral table */
	TF_Eponymous      = 0x00008000 /* An eponymous virtual table */
	TF_Strict         = 0x00010000 /* STRICT mode */
)

func ALWAYS(b bool) bool {
	return b
}

func NEVER(b bool) bool {
	return b
}

/*
** Allowed values for Table.eTabType
 */
const (
	TABTYP_NORM = 0 /* Ordinary table */
	TABTYP_VTAB = 1 /* Virtual table */
	TABTYP_VIEW = 2 /* A view */
)

// TODO: I made these up
type Hash uint64
type Pgno uint64

/*
** The bitmask datatype defined below is used for various optimizations.
**
** Changing this from a 64-bit to a 32-bit type limits the number of
** tables in a join to 32 instead of 64.  But it also reduces the size
** of the library by 738 bytes on ix86.
 */
type Bitmask uint64

/*
** Estimated quantities used for query planning are stored as 16-bit
** logarithms.  For quantity X, the value stored is 10*log2(X).  This
** gives a possible range of values of approximately 1.0e986 to 1e-986.
** But the allowed values are "grainy".  Not every value is representable.
** For example, quantities 16 and 17 are both represented by a LogEst
** of 40.  However, since LogEst quantities are suppose to be estimates,
** not exact values, this imprecision is not a problem.
**
** "LogEst" is short for "Logarithmic Estimate".
**
** Examples:
**      1 -> 0              20 -> 43          10000 -> 132
**      2 -> 10             25 -> 46          25000 -> 146
**      3 -> 16            100 -> 66        1000000 -> 199
**      4 -> 20           1000 -> 99        1048576 -> 200
**     10 -> 33           1024 -> 100    4294967296 -> 320
**
** The LogEst can be negative to indicate fractional values.
** Examples:
**
**    0.5 -> -10           0.1 -> -33        0.0625 -> -40
 */
type LogEst int16

/*
** The datatype used to store estimates of the number of rows in a
** table or index.  This is an unsigned integer type.  For 99.9% of
** the world, a 32-bit integer is sufficient.  But a 64-bit integer
** can be used at compile-time if desired.
 */
type tRowcnt uint64 /* 64-bit only if requested at compile-time */

/*
** The datatype ynVar is a signed integer, either 16-bit or 32-bit.
** Usually it is 16-bits.  But if SQLITE_MAX_VARIABLE_NUMBER is greater
** than 32767 we have to make it 32-bit.  16-bit is preferred because
** it uses less memory in the Expr object, which is a big memory user
** in systems with lots of prepared statements.  And few applications
** need more than about 10 or 20 variables.  But some extreme users want
** to have prepared statements with over 32766 variables, and for them
** the option is available (at compile-time).
 */
type ynVar int

/*
** An instance of this structure contains information needed to generate
** code for a SELECT that contains aggregate functions.
**
** If Expr.op==TK_AGG_COLUMN or TK_AGG_FUNCTION then Expr.pAggInfo is a
** pointer to this structure.  The Expr.iAgg field is the index in
** AggInfo.aCol[] or AggInfo.aFunc[] of information needed to generate
** code for that node.
**
** AggInfo.pGroupBy and AggInfo.aFunc.pExpr point to fields within the
** original Select structure that describes the SELECT statement.  These
** fields do not need to be freed when deallocating the AggInfo structure.
 */
type AggInfo struct {
	directMode uint8 /* Direct rendering mode means take data directly
	 ** from source tables rather than from accumulators */
	useSortingIdx uint8 /* In direct mode, reference the sorting index rather
	 ** than the source table */
	sortingIdx     int       /* Cursor number of the sorting index */
	sortingIdxPTab int       /* Cursor number of pseudo-table */
	nSortingColumn int       /* Number of columns in the sorting index */
	mnReg, mxReg   int       /* Range of registers allocated for aCol and aFunc */
	pGroupBy       *ExprList /* The group by clause */
	aCol           *struct { /* For each column used in source tables */
		pTab          *Table /* Source table */
		pCExpr        *Expr  /* The original expression */
		iTable        int    /* Cursor number of the source table */
		iMem          int    /* Memory location that acts as accumulator */
		iColumn       int16  /* Column number within the source table */
		iSorterColumn int16  /* Column number in the sorting index */
	}
	nColumn      int /* Number of used entries in aCol[] */
	nAccumulator int /* Number of columns that show through to the output.
	 ** Additional columns are used only as parameters to
	 ** aggregate functions */
	aFunc *struct { /* For each aggregate function */
		pFExpr    *Expr    /* Expression encoding the function */
		pFunc     *FuncDef /* The aggregate function implementation */
		iMem      int      /* Memory location that acts as accumulator */
		iDistinct int      /* Ephemeral table used to enforce DISTINCT */
		iDistAddr int      /* Address of OP_OpenEphemeral */
	}
	nFunc int    /* Number of entries in aFunc[] */
	selId uint32 /* Select to which this AggInfo belongs */
}

/*
** Information about each column of an SQL table is held in an instance
** of the Column structure, in the Table.aCol[] array.
**
** Definitions:
**
**   "table column index"     This is the index of the column in the
**                            Table.aCol[] array, and also the index of
**                            the column in the original CREATE TABLE stmt.
**
**   "storage column index"   This is the index of the column in the
**                            record BLOB generated by the OP_MakeRecord
**                            opcode.  The storage column index is less than
**                            or equal to the table column index.  It is
**                            equal if and only if there are no VIRTUAL
**                            columns to the left.
**
** Notes on zCnName:
** The zCnName field stores the name of the column, the datatype of the
** column, and the collating sequence for the column, in that order, all in
** a single allocation.  Each string is 0x00 terminated.  The datatype
** is only included if the COLFLAG_HASTYPE bit of colFlags is set and the
** collating sequence name is only included if the COLFLAG_HASCOLL bit is
** set.
 */
type Column struct {
	zCnName []byte /* Name of this column */
	// unsigned notNull :4;  /* An OE_ code for handling a NOT NULL constraint */
	// unsigned eCType :4;   /* One of the standard types */
	affinity rune   /* One of the SQLITE_AFF_... values */
	szEst    uint8  /* Est size of value in this column. sizeof(INT)==1 */
	hName    uint8  /* Column name hash for faster lookup */
	iDflt    uint16 /* 1-based index of DEFAULT.  0 means "none" */
	colFlags uint16 /* Boolean properties.  See COLFLAG_ defines below */
}

/*
** A single common table expression
 */
type Cte struct {
	zName   []byte    /* Name of this CTE */
	pCols   *ExprList /* List of explicit column names, or NULL */
	pSelect *Select   /* The definition of this CTE */
	zCteErr []byte    /* Error message for circular references */
	pUse    *CteUse   /* Usage information for this CTE */
	eM10d   uint8     /* The MATERIALIZED flag */
}

/*
** The Cte object is not guaranteed to persist for the entire duration
** of code generation.  (The query flattener or other parser tree
** edits might delete it.)  The following object records information
** about each Common Table Expression that must be preserved for the
** duration of the parse.
**
** The CteUse objects are freed using sqlite3ParserAddCleanup() rather
** than sqlite3SelectDelete(), which is what enables them to persist
** until the end of code generation.
 */
type CteUse struct {
	nUse    int    /* Number of users of this CTE */
	addrM9e int    /* Start of subroutine to compute materialization */
	regRtn  int    /* Return address register for addrM9e subroutine */
	iCur    int    /* Ephemeral table holding the materialization */
	nRowEst LogEst /* Estimated number of rows in the table */
	eM10d   uint8  /* The MATERIALIZED flag */
}

/*
** Each node of an expression in the parse tree is an instance
** of this structure.
**
** Expr.op is the opcode. The integer parser token codes are reused
** as opcodes here. For example, the parser defines TK_GE to be an integer
** code representing the ">=" operator. This same integer code is reused
** to represent the greater-than-or-equal-to operator in the expression
** tree.
**
** If the expression is an SQL literal (TK_INTEGER, TK_FLOAT, TK_BLOB,
** or TK_STRING), then Expr.u.zToken contains the text of the SQL literal. If
** the expression is a variable (TK_VARIABLE), then Expr.u.zToken contains the
** variable name. Finally, if the expression is an SQL function (TK_FUNCTION),
** then Expr.u.zToken contains the name of the function.
**
** Expr.pRight and Expr.pLeft are the left and right subexpressions of a
** binary operator. Either or both may be NULL.
**
** Expr.x.pList is a list of arguments if the expression is an SQL function,
** a CASE expression or an IN expression of the form "<lhs> IN (<y>, <z>...)".
** Expr.x.pSelect is used if the expression is a sub-select or an expression of
** the form "<lhs> IN (SELECT ...)". If the EP_xIsSelect bit is set in the
** Expr.flags mask, then Expr.x.pSelect is valid. Otherwise, Expr.x.pList is
** valid.
**
** An expression of the form ID or ID.ID refers to a column in a table.
** For such expressions, Expr.op is set to TK_COLUMN and Expr.iTable is
** the integer cursor number of a VDBE cursor pointing to that table and
** Expr.iColumn is the column number for the specific column.  If the
** expression is used as a result in an aggregate SELECT, then the
** value is also stored in the Expr.iAgg column in the aggregate so that
** it can be accessed after all aggregates are computed.
**
** If the expression is an unbound variable marker (a question mark
** character '?' in the original SQL) then the Expr.iTable holds the index
** number for that variable.
**
** If the expression is a subquery then Expr.iColumn holds an integer
** register number containing the result of the subquery.  If the
** subquery gives a constant result, then iTable is -1.  If the subquery
** gives a different answer at different times during statement processing
** then iTable is the address of a subroutine that computes the subquery.
**
** If the Expr is of type OP_Column, and the table it is selecting from
** is a disk table or the "old.*" pseudo-table, then pTab points to the
** corresponding table definition.
**
** ALLOCATION NOTES:
**
** Expr objects can use a lot of memory space in database schema.  To
** help reduce memory requirements, sometimes an Expr object will be
** truncated.  And to reduce the number of memory allocations, sometimes
** two or more Expr objects will be stored in a single memory allocation,
** together with Expr.u.zToken strings.
**
** If the EP_Reduced and EP_TokenOnly flags are set when
** an Expr object is truncated.  When EP_Reduced is set, then all
** the child Expr objects in the Expr.pLeft and Expr.pRight subtrees
** are contained within the same memory allocation.  Note, however, that
** the subtrees in Expr.x.pList or Expr.x.pSelect are always separately
** allocated, regardless of whether or not EP_Reduced is set.
 */
type Expr struct {
	op      uint8 /* Operation performed by this node */
	affExpr rune  /* affinity, or RAISE type */
	op2     uint8 /* TK_REGISTER/TK_TRUTH: original value of Expr.op
	 ** TK_COLUMN: the value of p5 for OP_Column
	 ** TK_AGG_FUNCTION: nesting depth
	 ** TK_FUNCTION: NC_SelfRef flag if needs OP_PureFunc */
	vvaFlags uint8  /* Verification flags. */
	flags    uint32 /* Various flags.  EP_* See below */
	u        struct {
		zToken *rune /* Token value. Zero terminated and dequoted */
		iValue int   /* Non-negative integer value if EP_IntValue */
	}

	/* If the EP_TokenOnly flag is set in the Expr.flags mask, then no
	 ** space is allocated for the fields below this point. An attempt to
	 ** access them will result in a segfault or malfunction.
	 *********************************************************************/

	pLeft  *Expr /* Left subnode */
	pRight *Expr /* Right subnode */
	x      struct {
		pList   *ExprList /* op = IN, EXISTS, SELECT, CASE, FUNCTION, BETWEEN */
		pSelect *Select   /* EP_xIsSelect and op = IN, EXISTS, SELECT */
	}

	/* If the EP_Reduced flag is set in the Expr.flags mask, then no
	 ** space is allocated for the fields below this point. An attempt to
	 ** access them will result in a segfault or malfunction.
	 *********************************************************************/
	nHeight int /* Height of the tree headed by this node */
	iTable  int /* TK_COLUMN: cursor number of table holding column
	 ** TK_REGISTER: register number
	 ** TK_TRIGGER: 1 -> new, 0 -> old
	 ** EP_Unlikely:  134217728 times likelihood
	 ** TK_IN: ephemerial table holding RHS
	 ** TK_SELECT_COLUMN: Number of columns on the LHS
	 ** TK_SELECT: 1st register of result vector */
	iColumn ynVar /* TK_COLUMN: column index.  -1 for rowid.
	 ** TK_VARIABLE: variable number (always >= 1).
	 ** TK_SELECT_COLUMN: column of the result vector */
	iAgg int16 /* Which entry in pAggInfo->aCol[] or ->aFunc[] */
	w    struct {
		iJoin int /* If EP_FromJoin, the right table of the join */
		iOfst int /* else: start of token from start of statement */
	}
	pAggInfo *AggInfo /* Used by TK_AGG_COLUMN and TK_AGG_FUNCTION */
	y        struct {
		pTab *Table /* TK_COLUMN: Table containing column. Can be NULL
		 ** for a column of an index on an expression */
		pWin *Window  /* EP_WinFunc: Window/Filter defn for a function */
		sub  struct { /* TK_IN, TK_SELECT, and TK_EXISTS */
			iAddr     int /* Subroutine entry address */
			regReturn int /* Register used to hold return address */
		}
	}
}

/*
** A list of expressions.  Each expression may optionally have a
** name.  An expr/name combination can be used in several ways, such
** as the list of "expr AS ID" fields following a "SELECT" or in the
** list of "ID = expr" items in an UPDATE.  A list of expressions can
** also be used as the argument to a function, in which case the a.zName
** field is not used.
**
** In order to try to keep memory usage down, the Expr.a.zEName field
** is used for multiple purposes:
**
**     eEName          Usage
**    ----------       -------------------------
**    ENAME_NAME       (1) the AS of result set column
**                     (2) COLUMN= of an UPDATE
**
**    ENAME_TAB        DB.TABLE.NAME used to resolve names
**                     of subqueries
**
**    ENAME_SPAN       Text of the original result set
**                     expression.
 */
type ExprList struct {
	nExpr  int        /* Number of expressions on the list */
	nAlloc int        /* Number of a[] slots allocated */
	a      []struct { /* For each expression in the list */
		pExpr     *Expr  /* The parse tree for this expression */
		zEName    []byte /* Token associated with this expression */
		sortFlags uint8  /* Mask of KEYINFO_ORDER_* flags */
		// unsigned eEName :2;     /* Meaning of zEName */
		// unsigned done :1;       /* A flag to indicate when processing is finished */
		// unsigned reusable :1;   /* Constant expression is reusable */
		// unsigned bSorterRef :1; /* Defer evaluation until after sorting */
		// unsigned bNulls: 1;     /* True if explicit "NULLS FIRST/LAST" */
		// unsigned bUsed: 1;      /* This column used in a SF_NestedFrom subquery */
		u struct {
			x struct { /* Used by any ExprList other than Parse.pConsExpr */
				iOrderByCol uint16 /* For ORDER BY, column number in result set */
				iAlias      uint16 /* Index into Parse.aAlias[] for zName */
			}
			iConstExprReg int /* Register in which Expr value is cached. Used only
			 ** by Parse.pConstExpr */
		}
	} /* One slot for each expression in the list */
}

/*
** Each foreign key constraint is an instance of the following structure.
**
** A foreign key is associated with two tables.  The "from" table is
** the table that contains the REFERENCES clause that creates the foreign
** key.  The "to" table is the table that is named in the REFERENCES clause.
** Consider this example:
**
**     CREATE TABLE ex1(
**       a INTEGER PRIMARY KEY,
**       b INTEGER CONSTRAINT fk1 REFERENCES ex2(x)
**     );
**
** For foreign key "fk1", the from-table is "ex1" and the to-table is "ex2".
** Equivalent names:
**
**     from-table == child-table
**       to-table == parent-table
**
** Each REFERENCES clause generates an instance of the following structure
** which is attached to the from-table.  The to-table need not exist when
** the from-table is created.  The existence of the to-table is not checked.
**
** The list of all parents for child Table X is held at X.pFKey.
**
** A list of all children for a table named Z (which might not even exist)
** is held in Schema.fkeyHash with a hash key of Z.
 */
type FKey struct {
	pFrom     *Table /* Table containing the REFERENCES clause (aka: Child) */
	pNextFrom *FKey  /* Next FKey with the same in pFrom. Next parent of pFrom */
	zTo       []byte /* Name of table that the key points to (aka: Parent) */
	pNextTo   *FKey  /* Next with the same zTo. Next child of zTo. */
	pPrevTo   *FKey  /* Previous with the same zTo */
	nCol      int    /* Number of columns in this key */
	/* EV: R-30323-21917 */
	isDeferred uint8      /* True if constraint checking is deferred till COMMIT */
	aAction    [2]uint8   /* ON DELETE and ON UPDATE actions, respectively */
	apTrigger  [2]Trigger /* Triggers for aAction[] actions */
	aCol       []struct { /* Mapping of columns in pFrom to columns in zTo */
		iFrom int    /* Index of column in pFrom */
		zCol  []byte /* Name of column in zTo.  If NULL use PRIMARY KEY */
	} /* One entry for each of nCol columns */
}

/*
** SQLite supports many different ways to resolve a constraint
** error.  ROLLBACK processing means that a constraint violation
** causes the operation in process to fail and for the current transaction
** to be rolled back.  ABORT processing means the operation in process
** fails and any prior changes from that one operation are backed out,
** but the transaction is not rolled back.  FAIL processing means that
** the operation in progress stops and returns an error code.  But prior
** changes due to the same operation are not backed out and no rollback
** occurs.  IGNORE means that the particular row that caused the constraint
** error is not inserted or updated.  Processing continues and no error
** is returned.  REPLACE means that preexisting database rows that caused
** a UNIQUE constraint violation are removed so that the new insert or
** update can proceed.  Processing continues and no error is reported.
** UPDATE applies to insert operations only and means that the insert
** is omitted and the DO UPDATE clause of an upsert is run instead.
**
** RESTRICT, SETNULL, SETDFLT, and CASCADE actions apply only to foreign keys.
** RESTRICT is the same as ABORT for IMMEDIATE foreign keys and the
** same as ROLLBACK for DEFERRED keys.  SETNULL means that the foreign
** key is set to NULL.  SETDFLT means that the foreign key is set
** to its default value.  CASCADE means that a DELETE or UPDATE of the
** referenced table row is propagated into the row that holds the
** foreign key.
**
** The OE_Default value is a place holder that means to use whatever
** conflict resolution algorthm is required from context.
**
** The following symbolic values are used to record which type
** of conflict resolution action to take.
 */
const (
	OE_None     = 0  /* There is no constraint to check */
	OE_Rollback = 1  /* Fail the operation and rollback the transaction */
	OE_Abort    = 2  /* Back out changes but do no rollback transaction */
	OE_Fail     = 3  /* Stop the operation but leave all prior changes */
	OE_Ignore   = 4  /* Ignore the error. Do not do the INSERT or UPDATE */
	OE_Replace  = 5  /* Delete existing record, then do INSERT or UPDATE */
	OE_Update   = 6  /* Process as a DO UPDATE in an upsert */
	OE_Restrict = 7  /* OE_Abort for IMMEDIATE, OE_Rollback for DEFERRED */
	OE_SetNull  = 8  /* Set the foreign key value to NULL */
	OE_SetDflt  = 9  /* Set the foreign key value to its default */
	OE_Cascade  = 10 /* Cascade the changes */
	OE_Default  = 11 /* Do whatever the default action is */
)

/*
** Each SQL function is defined by an instance of the following
** structure.  For global built-in functions (ex: substr(), max(), count())
** a pointer to this structure is held in the sqlite3BuiltinFunctions object.
** For per-connection application-defined functions, a pointer to this
** structure is held in the db->aHash hash table.
**
** The u.pHash field is used by the global built-ins.  The u.pDestructor
** field is used by per-connection app-def functions.
 */
type FuncDef struct {
	nArg      int8        /* Number of arguments.  -1 means unlimited */
	funcFlags uint32      /* Some combination of SQLITE_FUNC_* */
	pUserData interface{} /* User data parameter */
	pNext     *FuncDef    /* Next function with same name */
	// void (*xSFunc)(sqlite3_context*,int,sqlite3_value**); /* func or agg-step */
	// void (*xFinalize)(sqlite3_context*);                  /* Agg finalizer */
	// void (*xValue)(sqlite3_context*);                     /* Current agg value */
	// void (*xInverse)(sqlite3_context*,int,sqlite3_value**); /* inverse agg-step */
	zName []byte /* SQL name of the function. */
	u     struct {
		pHash       *FuncDef        /* Next with a different name but the same hash */
		pDestructor *FuncDestructor /* Reference counted destructor function */
	} /* pHash if SQLITE_FUNC_BUILTIN, pDestructor otherwise */
}

/*
** This structure encapsulates a user-function destructor callback (as
** configured using create_function_v2()) and a reference counter. When
** create_function_v2() is called to create a function with a destructor,
** a single object of this type is allocated. FuncDestructor.nRef is set to
** the number of FuncDef objects created (either 1 or 3, depending on whether
** or not the specified encoding is SQLITE_ANY). The FuncDef.pDestructor
** member of each of the new FuncDef objects is set to point to the allocated
** FuncDestructor.
**
** Thereafter, when one of the FuncDef objects is deleted, the reference
** count on this object is decremented. When it reaches 0, the destructor
** is invoked and the FuncDestructor structure freed.
 */
type FuncDestructor struct {
	nRef int
	// void (*xDestroy)(void *);
	pUserData interface{}
}

/*
** An instance of this structure can hold a simple list of identifiers,
** such as the list "a,b,c" in the following statements:
**
**      INSERT INTO t(a,b,c) VALUES ...;
**      CREATE INDEX idx ON t(a,b,c);
**      CREATE TRIGGER trig BEFORE UPDATE ON t(a,b,c) ...;
**
** The IdList.a.idx field is used when the IdList represents the list of
** column names after a table name in an INSERT statement.  In the statement
**
**     INSERT INTO t(a,b,c) ...
**
** If "a" is the k-th column of table "t", then IdList.a[0].idx==k.
 */
type IdList struct {
	nId int   /* Number of identifiers on the list */
	eU4 uint8 /* Which element of a.u4 is valid */
	a   []struct {
		zName []byte /* Name of the identifier */
		idx   int    /* Index in some Table.aCol[] of a column named zName */
		pExpr *Expr  /* Expr to implement a USING variable -- NOT USED */
	}
}

/*
** Each SQL index is represented in memory by an
** instance of the following structure.
**
** The columns of the table that are to be indexed are described
** by the aiColumn[] field of this structure.  For example, suppose
** we have the following table and index:
**
**     CREATE TABLE Ex1(c1 int, c2 int, c3 text);
**     CREATE INDEX Ex2 ON Ex1(c3,c1);
**
** In the Table structure describing Ex1, nCol==3 because there are
** three columns in the table.  In the Index structure describing
** Ex2, nColumn==2 since 2 of the 3 columns of Ex1 are indexed.
** The value of aiColumn is {2, 0}.  aiColumn[0]==2 because the
** first column to be indexed (c3) has an index of 2 in Ex1.aCol[].
** The second column to be indexed (c1) has an index of 0 in
** Ex1.aCol[], hence Ex2.aiColumn[1]==0.
**
** The Index.onError field determines whether or not the indexed columns
** must be unique and what to do if they are not.  When Index.onError=OE_None,
** it means this is not a unique index.  Otherwise it is a unique index
** and the value of Index.onError indicate the which conflict resolution
** algorithm to employ whenever an attempt is made to insert a non-unique
** element.
**
** While parsing a CREATE TABLE or CREATE INDEX statement in order to
** generate VDBE code (as opposed to parsing one read from an sqlite_schema
** table as part of parsing an existing database schema), transient instances
** of this structure may be created. In this case the Index.tnum variable is
** used to store the address of a VDBE instruction, not a database page
** number (it cannot - the database page is not allocated until the VDBE
** program is executed). See convertToWithoutRowidTable() for details.
 */
type Index struct {
	zName         []byte    /* Name of this index */
	aiColumn      *int16    /* Which columns are used by this index.  1st is 0 */
	aiRowLogEst   *LogEst   /* From ANALYZE: Est. rows selected by each column */
	pTable        *Table    /* The SQL table being indexed */
	zColAff       []byte    /* String defining the affinity of each column */
	pNext         *Index    /* The next index associated with the same table */
	pSchema       *Schema   /* Schema containing this index */
	aSortOrder    *uint8    /* for each column: True==DESC, False==ASC */
	azColl        []byte    /* Array of collation sequence names for index */
	pPartIdxWhere *Expr     /* WHERE clause for partial indices */
	aColExpr      *ExprList /* Column expressions */
	tnum          Pgno      /* DB Page containing root of this index */
	szIdxRow      LogEst    /* Estimated average row size in bytes */
	nKeyCol       uint16    /* Number of columns forming the key */
	nColumn       uint16    /* Number of columns stored in the index */
	onError       uint8     /* OE_Abort, OE_Ignore, OE_Replace, or OE_None */
	// unsigned idxType:2;      /* 0:Normal 1:UNIQUE, 2:PRIMARY KEY, 3:IPK */
	// unsigned bUnordered:1;   /* Use this index for == or IN queries only */
	// unsigned uniqNotNull:1;  /* True if UNIQUE and NOT NULL for all columns */
	// unsigned isResized:1;    /* True if resizeIndexObject() has been called */
	// unsigned isCovering:1;   /* True if this is a covering index */
	// unsigned noSkipScan:1;   /* Do not try to use skip-scan if true */
	// unsigned hasStat1:1;     /* aiRowLogEst values come from sqlite_stat1 */
	// unsigned bNoQuery:1;     /* Do not use this index to optimize queries */
	// unsigned bAscKeyBug:1;   /* True if the bba7b69f9849b5bf bug applies */
	// unsigned bHasVCol:1;     /* Index references one or more VIRTUAL columns */
	nSample     int          /* Number of elements in aSample[] */
	nSampleCol  int          /* Size of IndexSample.anEq[] and so on */
	aAvgEq      *tRowcnt     /* Average nEq values for keys not in aSample */
	aSample     *IndexSample /* Samples of the left-most key */
	iRowEst     *tRowcnt     /* Non-logarithmic stat1 data for this index */
	nRowEst0    tRowcnt      /* Non-logarithmic number of rows in the index */
	colNotIdxed Bitmask      /* 0 for unindexed columns in pTab */
}

/*
** Allowed values for Index.idxType
 */
const (
	SQLITE_IDXTYPE_APPDEF     = 0 /* Created using CREATE INDEX */
	SQLITE_IDXTYPE_UNIQUE     = 1 /* Implements a UNIQUE constraint */
	SQLITE_IDXTYPE_PRIMARYKEY = 2 /* Is the PRIMARY KEY for the table */
	SQLITE_IDXTYPE_IPK        = 3 /* INTEGER PRIMARY KEY index */
)

/*
** Each sample stored in the sqlite_stat4 table is represented in memory
** using a structure of this type.  See documentation at the top of the
** analyze.c source file for additional information.
 */
type IndexSample struct {
	p     interface{} /* Pointer to sampled record */
	n     int         /* Size of record in bytes */
	anEq  *tRowcnt    /* Est. number of rows where the key equals this sample */
	anLt  *tRowcnt    /* Est. number of rows where key is less than this sample */
	anDLt *tRowcnt    /* Est. number of distinct keys less than this sample */
}

/*
** The OnOrUsing object represents either an ON clause or a USING clause.
** It can never be both at the same time, but it can be neither.
 */
type OnOrUsing struct {
	pOn    *Expr   /* The ON clause of a join */
	pUsing *IdList /* The USING clause of a join */
}

/*
** An SQL parser context.  A copy of this structure is passed through
** the parser and down into all the parser action routine in order to
** carry around information that is global to the entire parse.
**
** The structure is divided into two parts.  When the parser and code
** generate call themselves recursively, the first part of the structure
** is constant but the second part is reset at the beginning and end of
** each recursion.
**
** The nTableLock and aTableLock variables are only used if the shared-cache
** feature is enabled (if sqlite3Tsd()->useSharedData is true). They are
** used to store the set of table-locks required by the statement being
** compiled. Function sqlite3TableLock() is used to add entries to the
** list.
 */
type Parse struct {
	db      *sqlite3 /* The main database structure */
	zErrMsg []byte   /* An error message */
	// Vdbe *pVdbe;         /* An engine for executing database bytecode */
	rc               int   /* Return code from execution */
	colNamesSet      uint8 /* TRUE after OP_ColumnName has been issued to pVdbe */
	checkSchema      uint8 /* Causes schema cookie check after an error */
	nested           uint8 /* Number of nested calls to the parser/code generator */
	nTempReg         uint8 /* Number of temporary registers in aTempReg[] */
	isMultiWrite     uint8 /* True if statement may modify/insert multiple rows */
	mayAbort         uint8 /* True if statement may throw an ABORT exception */
	hasCompound      uint8 /* Need to invoke convertCompoundSelectToSubquery() */
	okConstFactor    uint8 /* OK to factor out constants */
	disableLookaside uint8 /* Number of times lookaside has been disabled */
	disableVtab      uint8 /* Disable all virtual tables for this parse */
	withinRJSubrtn   uint8 /* Nesting level for RIGHT JOIN body subroutines */
	// # if defined(SQLITE_DEBUG) || defined(SQLITE_COVERAGE_TEST)
	earlyCleanup uint8 /* OOM inside sqlite3ParserAddCleanup() */
	// #endif

	nRangeReg int /* Size of the temporary register block */
	iRangeReg int /* First register in temporary register block */
	nErr      int /* Number of errors seen */
	nTab      int /* Number of previously allocated VDBE cursors */
	nMem      int /* Number of memory cells used so far */
	szOpAlloc int /* Bytes of memory space allocated for Vdbe.aOp[] */
	iSelfTab  int /* Table associated with an index on expr, or negative
	 ** of the base register during check-constraint eval */
	nLabel      int  /* The *negative* of the number of labels used */
	nLabelAlloc int  /* Number of slots in aLabel */
	aLabel      *int /* Space to hold the labels */
	//   ExprList *pConstExpr;/* Constant expressions */
	constraintName Token /* Name of the constraint currently being parsed */
	//   yDbMask writeMask;   /* Start a write transaction on these databases */
	//   yDbMask cookieMask;  /* Bitmask of schema verified databases */
	regRowid int /* Register holding rowid of CREATE TABLE entry */
	regRoot  int /* Register holding root page number for new objects */
	nMaxArg  int /* Max args passed to user function by sub-program */
	nSelect  int /* Number of SELECT stmts. Counter for Select.selId */
	// // #ifndef SQLITE_OMIT_SHARED_CACHE
	nTableLock int /* Number of locks in aTableLock */
	//   TableLock *aTableLock; /* Required table locks for shared-cache mode */
	// // #endif
	//   AutoincInfo *pAinc;  /* Information about AUTOINCREMENT counters */
	//   Parse *pToplevel;    /* Parse structure for main program (or NULL) */
	//   Table *pTriggerTab;  /* Table triggers are being coded for */
	//   TriggerPrg *pTriggerPrg;  /* Linked list of coded triggers */
	//   ParseCleanup *pCleanup;   /* List of cleanup operations to run after parse */
	//   union {
	//     int addrCrTab;         /* Address of OP_CreateBtree on CREATE TABLE */
	//     Returning *pReturning; /* The RETURNING clause */
	//   } u1;
	nQueryLoop      uint32 /* Est number of iterations of a query (10*log2(N)) */
	oldmask         uint32 /* Mask of old.* columns referenced */
	newmask         uint32 /* Mask of new.* columns referenced */
	eTriggerOp      uint8  /* TK_UPDATE, TK_INSERT or TK_DELETE */
	bReturning      uint8  /* Coding a RETURNING trigger */
	eOrconf         uint8  /* Default ON CONFLICT policy for trigger steps */
	disableTriggers uint8  /* True to disable triggers */
	//
	//   /**************************************************************************
	//   ** Fields above must be initialized to zero.  The fields that follow,
	//   ** down to the beginning of the recursive section, do not need to be
	//   ** initialized as they will be set before being used.  The boundary is
	//   ** determined by offsetof(Parse,aTempReg).
	//   **************************************************************************/
	//
	//   int aTempReg[8];        /* Holding area for temporary registers */
	//   Parse *pOuterParse;     /* Outer Parse object when nested */
	//   Token sNameToken;       /* Token with unqualified schema object name */
	//
	//   /************************************************************************
	//   ** Above is constant between recursions.  Below is reset before and after
	//   ** each recursion.  The boundary between these two regions is determined
	//   ** using offsetof(Parse,sLastToken) so the sLastToken field must be the
	//   ** first field in the recursive region.
	//   ************************************************************************/
	//
	sLastToken   Token /* The last token parsed */
	nVar         ynVar /* Number of '?' variables seen in the SQL so far */
	iPkSortOrder uint8 /* ASC or DESC for INTEGER PRIMARY KEY */
	explain      uint8 /* True if the EXPLAIN flag is found on the query */
	eParseMode   uint8 /* PARSE_MODE_XXX constant */
	// // #ifndef SQLITE_OMIT_VIRTUALTABLE
	nVtabLock int /* Number of virtual tables to lock */
	// // #endif
	nHeight int /* Expression tree height of current sub-select */
	// // #ifndef SQLITE_OMIT_EXPLAIN
	addrExplain int /* Address of current OP_Explain opcode */
	// // j#endif
	//   VList *pVList;            /* Mapping between variable names and numbers */
	//   Vdbe *pReprepare;         /* VM being reprepared (sqlite3Reprepare()) */
	zTail     []byte /* All SQL text past the last semicolon parsed */
	pNewTable *Table /* A table being constructed by CREATE TABLE */
	pNewIndex *Index /* An index being constructed by CREATE INDEX.
	//                             ** Also used to hold redundant UNIQUE constraints
	//                             ** during a RENAME COLUMN */
	pNewTrigger  *Trigger /* Trigger under construct by a CREATE TRIGGER */
	zAuthContext []byte   /* The 6th parameter to db->xAuth callbacks */
	// #ifndef SQLITE_OMIT_VIRTUALTABLE
	sArg Token /* Complete text of a module argument */
	//   Table **apVtabLock;       /* Pointer to virtual tables needing locking */
	// #endif
	//   With *pWith;              /* Current WITH clause, or NULL */
	// #ifndef SQLITE_OMIT_ALTERTABLE
	//   RenameToken *pRename;     /* Tokens subject to renaming by ALTER TABLE */
	// #endif
}

/*
** All current savepoints are stored in a linked list starting at
** sqlite3.pSavepoint. The first element in the list is the most recently
** opened savepoint. Savepoints are added to the list by the vdbe
** OP_Savepoint instruction.
 */
type Savepoint struct {
	zName            []byte     /* Savepoint name (nul-terminated) */
	nDeferredCons    int64      /* Number of deferred fk violations */
	nDeferredImmCons int64      /* Number of deferred imm fk. */
	pNext            *Savepoint /* Parent savepoint (if any) */
}

/*
** The following are used as the second parameter to sqlite3Savepoint(),
** and as the P1 argument to the OP_Savepoint instruction.
 */
const (
	SAVEPOINT_BEGIN    = 0
	SAVEPOINT_RELEASE  = 1
	SAVEPOINT_ROLLBACK = 2
)

/*
** An instance of the following structure stores a database schema.
**
** Most Schema objects are associated with a Btree.  The exception is
** the Schema for the TEMP databaes (sqlite3.aDb[1]) which is free-standing.
** In shared cache mode, a single Schema object can be shared by multiple
** Btrees that refer to the same underlying BtShared object.
**
** Schema objects are automatically deallocated when the last Btree that
** references them is destroyed.   The TEMP Schema is manually freed by
** sqlite3_close().
*
** A thread must be holding a mutex on the corresponding Btree in order
** to access Schema content.  This implies that the thread must also be
** holding a mutex on the sqlite3 connection pointer that owns the Btree.
** For a TEMP Schema, only the connection mutex is required.
 */
type Schema struct {
	schema_cookie int    /* Database schema version number for this file */
	iGeneration   int    /* Generation counter.  Incremented with each change */
	tblHash       Hash   /* All tables indexed by name */
	idxHash       Hash   /* All (named) indices indexed by name */
	trigHash      Hash   /* All triggers indexed by name */
	fkeyHash      Hash   /* All foreign keys by referenced table name */
	pSeqTab       *Table /* The sqlite_sequence table used by AUTOINCREMENT */
	file_format   uint8  /* Schema format version for this file */
	enc           uint8  /* Text encoding used by this database */
	schemaFlags   uint16 /* Flags associated with this schema */
	cache_size    int    /* Number of pages to use in the cache */
}

/*
** An instance of the following structure contains all information
** needed to generate code for a single SELECT statement.
**
** See the header comment on the computeLimitRegisters() routine for a
** detailed description of the meaning of the iLimit and iOffset fields.
**
** addrOpenEphm[] entries contain the address of OP_OpenEphemeral opcodes.
** These addresses must be stored so that we can go back and fill in
** the P4_KEYINFO and P2 parameters later.  Neither the KeyInfo nor
** the number of columns in P2 can be computed at the same time
** as the OP_OpenEphm instruction is coded because not
** enough information about the compound query is known at that point.
** The KeyInfo for addrOpenTran[0] and [1] contains collating sequences
** for the result set.  The KeyInfo for addrOpenEphm[2] contains collating
** sequences for the ORDER BY clause.
 */
type Select struct {
	op           uint8  /* One of: TK_UNION TK_ALL TK_INTERSECT TK_EXCEPT */
	nSelectRow   LogEst /* Estimated number of result rows */
	selFlags     uint32 /* Various SF_* values */
	iLimit       int
	iOffset      int       /* Memory registers holding LIMIT & OFFSET counters */
	selId        uint32    /* Unique identifier number for this SELECT */
	addrOpenEphm [2]int    /* OP_OpenEphem opcodes related to this select */
	pEList       *ExprList /* The fields of the result */
	pSrc         *SrcList  /* The FROM clause */
	pWhere       *Expr     /* The WHERE clause */
	pGroupBy     *ExprList /* The GROUP BY clause */
	pHaving      *Expr     /* The HAVING clause */
	pOrderBy     *ExprList /* The ORDER BY clause */
	pPrior       *Select   /* Prior select in a compound select statement */
	pNext        *Select   /* Next select to the left in a compound */
	pLimit       *Expr     /* LIMIT expression. NULL means not used. */
	pWith        *With     /* WITH clause attached to this select. Or NULL. */
	pWin         *Window   /* List of window functions */
	pWinDefn     *Window   /* List of named window definitions */
}

/*
** Allowed values for Select.selFlags.  The "SF" prefix stands for
** "Select Flag".
**
** Value constraints (all checked via assert())
**     SF_HasAgg      == NC_HasAgg
**     SF_MinMaxAgg   == NC_MinMaxAgg     == SQLITE_FUNC_MINMAX
**     SF_OrderByReqd == NC_OrderAgg      == SQLITE_FUNC_ANYORDER
**     SF_FixedLimit  == WHERE_USE_LIMIT
 */
const (
	SF_Distinct      = 0x0000001 /* Output should be DISTINCT */
	SF_All           = 0x0000002 /* Includes the ALL keyword */
	SF_Resolved      = 0x0000004 /* Identifiers have been resolved */
	SF_Aggregate     = 0x0000008 /* Contains agg functions or a GROUP BY */
	SF_HasAgg        = 0x0000010 /* Contains aggregate functions */
	SF_UsesEphemeral = 0x0000020 /* Uses the OpenEphemeral opcode */
	SF_Expanded      = 0x0000040 /* sqlite3SelectExpand() called on this */
	SF_HasTypeInfo   = 0x0000080 /* FROM subqueries have Table metadata */
	SF_Compound      = 0x0000100 /* Part of a compound query */
	SF_Values        = 0x0000200 /* Synthesized from VALUES clause */
	SF_MultiValue    = 0x0000400 /* Single VALUES term with multiple rows */
	SF_NestedFrom    = 0x0000800 /* Part of a parenthesized FROM clause */
	SF_MinMaxAgg     = 0x0001000 /* Aggregate containing min() or max() */
	SF_Recursive     = 0x0002000 /* The recursive part of a recursive CTE */
	SF_FixedLimit    = 0x0004000 /* nSelectRow set by a constant LIMIT */
	SF_MaybeConvert  = 0x0008000 /* Need convertCompoundSelectToSubquery() */
	SF_Converted     = 0x0010000 /* By convertCompoundSelectToSubquery() */
	SF_IncludeHidden = 0x0020000 /* Include hidden columns in output */
	SF_ComplexResult = 0x0040000 /* Result contains subquery or function */
	SF_WhereBegin    = 0x0080000 /* Really a WhereBegin() call.  Debug Only */
	SF_WinRewrite    = 0x0100000 /* Window function rewrite accomplished */
	SF_View          = 0x0200000 /* SELECT statement is a view */
	SF_NoopOrderBy   = 0x0400000 /* ORDER BY is ignored for this query */
	SF_UFSrcCheck    = 0x0800000 /* Check pSrc as required by UPDATE...FROM */
	SF_PushDown      = 0x1000000 /* SELECT has be modified by push-down opt */
	SF_MultiPart     = 0x2000000 /* Has multiple incompatible PARTITIONs */
	SF_CopyCte       = 0x4000000 /* SELECT statement is a copy of a CTE */
	SF_OrderByReqd   = 0x8000000 /* The ORDER BY clause may not be omitted */
)

/* True if S exists and has SF_NestedFrom */
// #define IsNestedFrom(S) ((S)!=0 && ((S)->selFlags&SF_NestedFrom)!=0)

/*
** The results of a SELECT can be distributed in several ways, as defined
** by one of the following macros.  The "SRT" prefix means "SELECT Result
** Type".
**
**     SRT_Union       Store results as a key in a temporary index
**                     identified by pDest->iSDParm.
**
**     SRT_Except      Remove results from the temporary index pDest->iSDParm.
**
**     SRT_Exists      Store a 1 in memory cell pDest->iSDParm if the result
**                     set is not empty.
**
**     SRT_Discard     Throw the results away.  This is used by SELECT
**                     statements within triggers whose only purpose is
**                     the side-effects of functions.
**
**     SRT_Output      Generate a row of output (using the OP_ResultRow
**                     opcode) for each row in the result set.
**
**     SRT_Mem         Only valid if the result is a single column.
**                     Store the first column of the first result row
**                     in register pDest->iSDParm then abandon the rest
**                     of the query.  This destination implies "LIMIT 1".
**
**     SRT_Set         The result must be a single column.  Store each
**                     row of result as the key in table pDest->iSDParm.
**                     Apply the affinity pDest->affSdst before storing
**                     results.  Used to implement "IN (SELECT ...)".
**
**     SRT_EphemTab    Create an temporary table pDest->iSDParm and store
**                     the result there. The cursor is left open after
**                     returning.  This is like SRT_Table except that
**                     this destination uses OP_OpenEphemeral to create
**                     the table first.
**
**     SRT_Coroutine   Generate a co-routine that returns a new row of
**                     results each time it is invoked.  The entry point
**                     of the co-routine is stored in register pDest->iSDParm
**                     and the result row is stored in pDest->nDest registers
**                     starting with pDest->iSdst.
**
**     SRT_Table       Store results in temporary table pDest->iSDParm.
**     SRT_Fifo        This is like SRT_EphemTab except that the table
**                     is assumed to already be open.  SRT_Fifo has
**                     the additional property of being able to ignore
**                     the ORDER BY clause.
**
**     SRT_DistFifo    Store results in a temporary table pDest->iSDParm.
**                     But also use temporary table pDest->iSDParm+1 as
**                     a record of all prior results and ignore any duplicate
**                     rows.  Name means:  "Distinct Fifo".
**
**     SRT_Queue       Store results in priority queue pDest->iSDParm (really
**                     an index).  Append a sequence number so that all entries
**                     are distinct.
**
**     SRT_DistQueue   Store results in priority queue pDest->iSDParm only if
**                     the same record has never been stored before.  The
**                     index at pDest->iSDParm+1 hold all prior stores.
**
**     SRT_Upfrom      Store results in the temporary table already opened by
**                     pDest->iSDParm. If (pDest->iSDParm<0), then the temp
**                     table is an intkey table - in this case the first
**                     column returned by the SELECT is used as the integer
**                     key. If (pDest->iSDParm>0), then the table is an index
**                     table. (pDest->iSDParm) is the number of key columns in
**                     each index record in this case.
 */
const (
	SRT_Union     = 1 /* Store result as keys in an index */
	SRT_Except    = 2 /* Remove result from a UNION index */
	SRT_Exists    = 3 /* Store 1 if the result is not empty */
	SRT_Discard   = 4 /* Do not save the results anywhere */
	SRT_DistFifo  = 5 /* Like SRT_Fifo, but unique results only */
	SRT_DistQueue = 6 /* Like SRT_Queue, but unique results only */

	/* The DISTINCT clause is ignored for all of the above.  Not that
	 ** IgnorableDistinct() implies IgnorableOrderby() */
	// #define IgnorableDistinct(X) ((X->eDest)<=SRT_DistQueue)

	SRT_Queue = 7 /* Store result in an queue */
	SRT_Fifo  = 8 /* Store result as data with an automatic rowid */

	/* The ORDER BY clause is ignored for all of the above */
	// #define IgnorableOrderby(X) ((X->eDest)<=SRT_Fifo)

	SRT_Output    = 9  /* Output each row of result */
	SRT_Mem       = 10 /* Store result in a memory cell */
	SRT_Set       = 11 /* Store results as keys in an index */
	SRT_EphemTab  = 12 /* Create transient tab and store like SRT_Table */
	SRT_Coroutine = 13 /* Generate a single row of result */
	SRT_Table     = 14 /* Store result as data with an automatic rowid */
	SRT_Upfrom    = 15 /* Store result as data with rowid */
)

/*
** An instance of this object describes where to put of the results of
** a SELECT statement.
 */
type SelectDest struct {
	eDest    uint8     /* How to dispose of the results.  One of SRT_* above. */
	iSDParm  int       /* A parameter used by the eDest disposal method */
	iSDParm2 int       /* A second parameter for the eDest disposal method */
	iSdst    int       /* Base register where results are written */
	nSdst    int       /* Number of registers allocated */
	zAffSdst []byte    /* Affinity used when eDest==SRT_Set */
	pOrderBy *ExprList /* Key columns for SRT_Queue and SRT_DistQueue */
}

/*
** The SrcItem object represents a single term in the FROM clause of a query.
** The SrcList object is mostly an array of SrcItems.
**
** Union member validity:
**
**    u1.zIndexedBy          fg.isIndexedBy && !fg.isTabFunc
**    u1.pFuncArg            fg.isTabFunc   && !fg.isIndexedBy
**    u2.pIBIndex            fg.isIndexedBy && !fg.isCte
**    u2.pCteUse             fg.isCte       && !fg.isIndexedBy
 */
type SrcItem struct {
	pSchema     *Schema /* Schema to which this item is fixed */
	zDatabase   []byte  /* Name of database holding this table */
	zName       []byte  /* Name of the table */
	zAlias      []byte  /* The "B" part of a "A AS B" phrase.  zName is the "A" */
	pTab        *Table  /* An SQL table corresponding to zName */
	pSelect     *Select /* A SELECT statement used in place of a table name */
	addrFillSub int     /* Address of subroutine to manifest a subquery */
	regReturn   int     /* Register holding return address of addrFillSub */
	regResult   int     /* Registers holding results of a co-routine */
	fg          struct {
		jointype uint8 /* Type of join between this table and the previous */
		// unsigned notIndexed :1;    /* True if there is a NOT INDEXED clause */
		// unsigned isIndexedBy :1;   /* True if there is an INDEXED BY clause */
		// unsigned isTabFunc :1;     /* True if table-valued-function syntax */
		// unsigned isCorrelated :1;  /* True if sub-query is correlated */
		// unsigned viaCoroutine :1;  /* Implemented as a co-routine */
		// unsigned isRecursive :1;   /* True for recursive reference in WITH */
		// unsigned fromDDL :1;       /* Comes from sqlite_schema */
		// unsigned isCte :1;         /* This is a CTE */
		// unsigned notCte :1;        /* This item may not match a CTE */
		// unsigned isUsing :1;       /* u3.pUsing is valid */
		// unsigned isSynthUsing :1;  /* u3.pUsing is synthensized from NATURAL */
		// unsigned isNestedFrom :1;  /* pSelect is a SF_NestedFrom subquery */
	}
	iCursor int /* The VDBE cursor number used to access this table */
	u3      struct {
		pOn    *Expr   /* fg.isUsing==0 =>  The ON clause of a join */
		pUsing *IdList /* fg.isUsing==1 =>  The USING clause of a join */
	}
	colUsed Bitmask /* Bit N (1<<N) set if column N of pTab is used */
	u1      struct {
		zIndexedBy []byte    /* Identifier from "INDEXED BY <zIndex>" clause */
		pFuncArg   *ExprList /* Arguments to table-valued-function */
	}
	u2 struct {
		IBIndex *Index  /* Index structure corresponding to u1.zIndexedBy */
		pCteUse *CteUse /* CTE Usage info info fg.isCte is true */
	}
}

/*
** The following structure describes the FROM clause of a SELECT statement.
** Each table or subquery in the FROM clause is a separate element of
** the SrcList.a[] array.
**
** With the addition of multiple database support, the following structure
** can also be used to describe a particular table such as the table that
** is modified by an INSERT, DELETE, or UPDATE statement.  In standard SQL,
** such a table must be a simple name: ID.  But in SQLite, the table can
** now be identified by a database name, a dot, then the table name: ID.ID.
**
** The jointype starts out showing the join type between the current table
** and the next table on the list.  The parser builds the list this way.
** But sqlite3SrcListShiftJoinType() later shifts the jointypes so that each
** jointype expresses the join between the table and the previous table.
**
** In the colUsed field, the high-order bit (bit 63) is set if the table
** contains more than 63 columns and the 64-th or later column is used.
 */
type SrcList struct {
	nSrc   int       /* Number of tables or subqueries in the FROM clause */
	nAlloc uint32    /* Number of entries allocated in a[] below */
	a      []SrcItem /* One entry for each identifier on the list */
}

/*
** The schema for each SQL table, virtual table, and view is represented
** in memory by an instance of the following structure.
 */
type Table struct {
	zName   []byte    /* Name of the table or view */
	aCol    *Column   /* Information about each column */
	pIndex  *Index    /* List of SQL indexes on this table. */
	zColAff []byte    /* String defining the affinity of each column */
	pCheck  *ExprList /* All CHECK constraints */
	/*   ... also used as column name list in a VIEW */
	tnum       Pgno   /* Root BTree page for this table */
	nTabRef    uint32 /* Number of pointers to this Table */
	tabFlags   uint32 /* Mask of TF_* values */
	iPKey      int16  /* If not negative, use aCol[iPKey] as the rowid */
	nCol       int16  /* Number of columns in this table */
	nNVCol     int16  /* Number of columns that are not VIRTUAL */
	nRowLogEst LogEst /* Estimated rows in table - from sqlite_stat1 table */
	szTabRow   LogEst /* Estimated size of each table row in bytes */
	costMult   LogEst /* Cost multiplier for using this table */
	keyConf    uint8  /* What to do in case of uniqueness conflict on iPKey */
	eTabType   uint8  /* 0: normal, 1: virtual, 2: view */
	u          struct {
		tab struct { /* Used by ordinary tables: */
			addColOffset int       /* Offset in CREATE TABLE stmt to add a new column */
			pFKey        *FKey     /* Linked list of all foreign keys in this table */
			pDfltList    *ExprList /* DEFAULT clauses on various columns.
			 ** Or the AS clause for generated columns. */
		}
		view struct { /* Used by views: */
			pSelect *Select /* View definition */
		}
		vtab struct { /* Used by virtual tables only: */
			nArg  int     /* Number of arguments to the module */
			azArg string  /* 0: module 1: schema 2: vtab name 3...: args */
			p     *VTable /* List of VTable objects. */
		}
	}
	pTrigger *Trigger /* List of triggers on this object */
	pSchema  *Schema  /* Schema that contains this table */
}

/*
** Each token coming out of the lexer is an instance of
** this structure.  Tokens are also used as part of an expression.
**
** The memory that "z" points to is owned by other objects.  Take care
** that the owner of the "z" string does not deallocate the string before
** the Token goes out of scope!  Very often, the "z" points to some place
** in the middle of the Parse.zSql text.  But it might also point to a
** static string.
 */
type Token struct {
	z []byte /* Text of the token.  Not NULL-terminated! */
	n uint   /* Number of characters in this token */
}

/*
** Each trigger present in the database schema is stored as an instance of
** struct Trigger.
**
** Pointers to instances of struct Trigger are stored in two ways.
** 1. In the "trigHash" hash table (part of the sqlite3* that represents the
**    database). This allows Trigger structures to be retrieved by name.
** 2. All triggers associated with a single table form a linked list, using the
**    pNext member of struct Trigger. A pointer to the first element of the
**    linked list is stored as the "pTrigger" member of the associated
**    struct Table.
**
** The "step_list" member points to the first element of a linked list
** containing the SQL statements specified as the trigger program.
 */
type Trigger struct {
	zName      []byte  /* The name of the trigger                        */
	table      []byte  /* The table or view to which the trigger applies */
	op         uint8   /* One of TK_DELETE, TK_UPDATE, TK_INSERT         */
	tr_tm      uint8   /* One of TRIGGER_BEFORE, TRIGGER_AFTER */
	bReturning uint8   /* This trigger implements a RETURNING clause */
	pWhen      *Expr   /* The WHEN clause of the expression (may be NULL) */
	pColumns   *IdList /* If this is an UPDATE OF <column-list> trigger,
	   the <column-list> is stored here */
	pSchema    *Schema      /* Schema containing the trigger */
	pTabSchema *Schema      /* Schema containing the table */
	step_list  *TriggerStep /* Link list of trigger program steps             */
	pNext      *Trigger     /* Next trigger associated with the table */
}

/*
** An instance of struct TriggerStep is used to store a single SQL statement
** that is a part of a trigger-program.
**
** Instances of struct TriggerStep are stored in a singly linked list (linked
** using the "pNext" member) referenced by the "step_list" member of the
** associated struct Trigger instance. The first element of the linked list is
** the first step of the trigger-program.
**
** The "op" member indicates whether this is a "DELETE", "INSERT", "UPDATE" or
** "SELECT" statement. The meanings of the other members is determined by the
** value of "op" as follows:
**
** (op == TK_INSERT)
** orconf    -> stores the ON CONFLICT algorithm
** pSelect   -> The content to be inserted - either a SELECT statement or
**              a VALUES clause.
** zTarget   -> Dequoted name of the table to insert into.
** pIdList   -> If this is an INSERT INTO ... (<column-names>) VALUES ...
**              statement, then this stores the column-names to be
**              inserted into.
** pUpsert   -> The ON CONFLICT clauses for an Upsert
**
** (op == TK_DELETE)
** zTarget   -> Dequoted name of the table to delete from.
** pWhere    -> The WHERE clause of the DELETE statement if one is specified.
**              Otherwise NULL.
**
** (op == TK_UPDATE)
** zTarget   -> Dequoted name of the table to update.
** pWhere    -> The WHERE clause of the UPDATE statement if one is specified.
**              Otherwise NULL.
** pExprList -> A list of the columns to update and the expressions to update
**              them to. See sqlite3Update() documentation of "pChanges"
**              argument.
**
** (op == TK_SELECT)
** pSelect   -> The SELECT statement
**
** (op == TK_RETURNING)
** pExprList -> The list of expressions that follow the RETURNING keyword.
**
 */
type TriggerStep struct {
	op uint8 /* One of TK_DELETE, TK_UPDATE, TK_INSERT, TK_SELECT,
	 ** or TK_RETURNING */
	orconf    uint8        /* OE_Rollback etc. */
	pTrig     *Trigger     /* The trigger that this step is a part of */
	pSelect   *Select      /* SELECT statement or RHS of INSERT INTO SELECT ... */
	zTarget   []byte       /* Target table for DELETE, UPDATE, INSERT */
	pFrom     *SrcList     /* FROM clause for UPDATE statement (if any) */
	pWhere    *Expr        /* The WHERE clause for DELETE or UPDATE steps */
	pExprList *ExprList    /* SET clause for UPDATE, or RETURNING clause */
	pIdList   *IdList      /* Column names for INSERT */
	pUpsert   *Upsert      /* Upsert clauses on an INSERT */
	zSpan     []byte       /* Original SQL text of this command */
	pNext     *TriggerStep /* Next in the link-list */
	pLast     *TriggerStep /* Last element in link-list. Valid for 1st elem only */
}

/*
** An instance of the following object describes a single ON CONFLICT
** clause in an upsert.
**
** The pUpsertTarget field is only set if the ON CONFLICT clause includes
** conflict-target clause.  (In "ON CONFLICT(a,b)" the "(a,b)" is the
** conflict-target clause.)  The pUpsertTargetWhere is the optional
** WHERE clause used to identify partial unique indexes.
**
** pUpsertSet is the list of column=expr terms of the UPDATE statement.
** The pUpsertSet field is NULL for a ON CONFLICT DO NOTHING.  The
** pUpsertWhere is the WHERE clause for the UPDATE and is NULL if the
** WHERE clause is omitted.
 */
type Upsert struct {
	pUpsertTarget      *ExprList /* Optional description of conflict target */
	pUpsertTargetWhere *Expr     /* WHERE clause for partial index targets */
	pUpsertSet         *ExprList /* The SET clause from an ON CONFLICT UPDATE */
	pUpsertWhere       *Expr     /* WHERE clause for the ON CONFLICT UPDATE */
	pNextUpsert        *Upsert   /* Next ON CONFLICT clause in the list */
	isDoUpdate         uint8     /* True for DO UPDATE.  False for DO NOTHING */
	/* Above this point is the parse tree for the ON CONFLICT clauses.
	 ** The next group of fields stores intermediate data. */
	pToFree interface{} /* Free memory when deleting the Upsert object */
	/* All fields above are owned by the Upsert object and must be freed
	 ** when the Upsert is destroyed.  The fields below are used to transfer
	 ** information from the INSERT processing down into the UPDATE processing
	 ** while generating code.  The fields below are owned by the INSERT
	 ** statement and will be freed by INSERT processing. */
	pUpsertIdx *Index   /* UNIQUE constraint specified by pUpsertTarget */
	pUpsertSrc *SrcList /* Table to be updated */
	regData    int      /* First register holding array of VALUES */
	iDataCur   int      /* Index of the data cursor */
	iIdxCur    int      /* Index of the first index cursor */
}

/*
** An object of this type is created for each virtual table present in
** the database schema.
**
** If the database schema is shared, then there is one instance of this
** structure for each database connection (sqlite3*) that uses the shared
** schema. This is because each database connection requires its own unique
** instance of the sqlite3_vtab* handle used to access the virtual table
** implementation. sqlite3_vtab* handles can not be shared between
** database connections, even when the rest of the in-memory database
** schema is shared, as the implementation often stores the database
** connection handle passed to it via the xConnect() or xCreate() method
** during initialization internally. This database connection handle may
** then be used by the virtual table implementation to access real tables
** within the database. So that they appear as part of the callers
** transaction, these accesses need to be made via the same database
** connection as that used to execute SQL operations on the virtual table.
**
** All VTable objects that correspond to a single table in a shared
** database schema are initially stored in a linked-list pointed to by
** the Table.pVTable member variable of the corresponding Table object.
** When an sqlite3_prepare() operation is required to access the virtual
** table, it searches the list for the VTable that corresponds to the
** database connection doing the preparing so as to use the correct
** sqlite3_vtab* handle in the compiled query.
**
** When an in-memory Table object is deleted (for example when the
** schema is being reloaded for some reason), the VTable objects are not
** deleted and the sqlite3_vtab* handles are not xDisconnect()ed
** immediately. Instead, they are moved from the Table.pVTable list to
** another linked list headed by the sqlite3.pDisconnect member of the
** corresponding sqlite3 structure. They are then deleted/xDisconnected
** next time a statement is prepared using said sqlite3*. This is done
** to avoid deadlock issues involving multiple sqlite3.mutex mutexes.
** Refer to comments above function sqlite3VtabUnlockList() for an
** explanation as to why it is safe to add an entry to an sqlite3.pDisconnect
** list without holding the corresponding sqlite3.mutex mutex.
**
** The memory for objects of this type is always allocated by
** sqlite3DbMalloc(), using the connection handle stored in VTable.db as
** the first argument.
 */
type VTable struct {
	// db *sqlite3;              /* Database connection associated with this table */
	// pMod *Module;             /* Pointer to module implementation */
	// pVtab *sqlite3_vtab;      /* Pointer to vtab instance */
	nRef        int     /* Number of pointers to this structure */
	bConstraint uint8   /* True if constraints are supported */
	eVtabRisk   uint8   /* Riskiness of allowing hacker access */
	iSavepoint  int     /* Depth of the SAVEPOINT stack */
	pNext       *VTable /* Next in linked list (see above) */
}

/*
** This object is used in various ways, most (but not all) related to window
** functions.
**
**   (1) A single instance of this structure is attached to the
**       the Expr.y.pWin field for each window function in an expression tree.
**       This object holds the information contained in the OVER clause,
**       plus additional fields used during code generation.
**
**   (2) All window functions in a single SELECT form a linked-list
**       attached to Select.pWin.  The Window.pFunc and Window.pExpr
**       fields point back to the expression that is the window function.
**
**   (3) The terms of the WINDOW clause of a SELECT are instances of this
**       object on a linked list attached to Select.pWinDefn.
**
**   (4) For an aggregate function with a FILTER clause, an instance
**       of this object is stored in Expr.y.pWin with eFrmType set to
**       TK_FILTER. In this case the only field used is Window.pFilter.
**
** The uses (1) and (2) are really the same Window object that just happens
** to be accessible in two different ways.  Use case (3) are separate objects.
 */
type Window struct {
	zName          []byte    /* Name of window (may be NULL) */
	zBase          []byte    /* Name of base window for chaining (may be NULL) */
	pPartition     *ExprList /* PARTITION BY clause */
	pOrderBy       *ExprList /* ORDER BY clause */
	eFrmType       uint8     /* TK_RANGE, TK_GROUPS, TK_ROWS, or 0 */
	eStart         uint8     /* UNBOUNDED, CURRENT, PRECEDING or FOLLOWING */
	eEnd           uint8     /* UNBOUNDED, CURRENT, PRECEDING or FOLLOWING */
	bImplicitFrame uint8     /* True if frame was implicitly specified */
	eExclude       uint8     /* TK_NO, TK_CURRENT, TK_TIES, TK_GROUP, or 0 */
	pStart         *Expr     /* Expression for "<expr> PRECEDING" */
	pEnd           *Expr     /* Expression for "<expr> FOLLOWING" */
	ppThis         *Window   /* Pointer to this object in Select.pWin list */
	pNextWin       *Window   /* Next window function belonging to this SELECT */
	pFilter        *Expr     /* The FILTER expression */
	pWFunc         *FuncDef  /* The function */
	iEphCsr        int       /* Partition buffer or Peer buffer */
	regAccum       int       /* Accumulator */
	regResult      int       /* Interim result */
	csrApp         int       /* Function cursor (used by min/max) */
	regApp         int       /* Function register (also used by min/max) */
	regPart        int       /* Array of registers for PARTITION BY values */
	pOwner         *Expr     /* Expression object this window is attached to */
	nBufferCol     int       /* Number of columns in buffer table */
	iArgCol        int       /* Offset of first argument for this function */
	regOne         int       /* Register containing constant value 1 */
	regStartRowid  int
	regEndRowid    int
	bExprArgs      uint8 /* Defer evaluation of window function arguments
	 ** due to the SQLITE_SUBTYPE flag */
}

/*
** An instance of the With object represents a WITH clause containing
** one or more CTEs (common table expressions).
 */
type With struct {
	nCte   int   /* Number of CTEs in the WITH clause */
	bView  int   /* Belongs to the outermost Select of a view */
	pOuter *With /* Containing WITH clause, or NULL */
	a      []Cte /* For each CTE in the WITH clause.... */
}

/*
** Each database connection is an instance of the following structure.
 */
type sqlite3 struct {
	//   sqlite3_vfs *pVfs;            /* OS Interface */
	//   struct Vdbe *pVdbe;           /* List of active virtual machines */
	//   CollSeq *pDfltColl;           /* BINARY collseq for the database encoding */
	//   sqlite3_mutex *mutex;         /* Connection mutex */
	//   Db *aDb;                      /* All backends */
	//   int nDb;                      /* Number of backends currently in use */
	//   u32 mDbFlags;                 /* flags recording internal state */
	//   u64 flags;                    /* flags settable by pragmas. See below */
	//   i64 lastRowid;                /* ROWID of most recent insert (see above) */
	//   i64 szMmap;                   /* Default mmap_size setting */
	//   u32 nSchemaLock;              /* Do not reset the schema when non-zero */
	//   unsigned int openFlags;       /* Flags passed to sqlite3_vfs.xOpen() */
	//   int errCode;                  /* Most recent error code (SQLITE_*) */
	//   int errByteOffset;            /* Byte offset of error in SQL statement */
	//   int errMask;                  /* & result codes with this before returning */
	//   int iSysErrno;                /* Errno value from last system error */
	//   u32 dbOptFlags;               /* Flags to enable/disable optimizations */
	//   u8 enc;                       /* Text encoding */
	//   u8 autoCommit;                /* The auto-commit flag. */
	//   u8 temp_store;                /* 1: file 2: memory 0: default */
	//   u8 mallocFailed;              /* True if we have seen a malloc failure */
	//   u8 bBenignMalloc;             /* Do not require OOMs if true */
	//   u8 dfltLockMode;              /* Default locking-mode for attached dbs */
	//   signed char nextAutovac;      /* Autovac setting after VACUUM if >=0 */
	//   u8 suppressErr;               /* Do not issue error messages if true */
	//   u8 vtabOnConflict;            /* Value to return for s3_vtab_on_conflict() */
	//   u8 isTransactionSavepoint;    /* True if the outermost savepoint is a TS */
	//   u8 mTrace;                    /* zero or more SQLITE_TRACE flags */
	//   u8 noSharedCache;             /* True if no shared-cache backends */
	//   u8 nSqlExec;                  /* Number of pending OP_SqlExec opcodes */
	//   u8 eOpenState;                /* Current condition of the connection */
	//   int nextPagesize;             /* Pagesize after VACUUM if >0 */
	//   i64 nChange;                  /* Value returned by sqlite3_changes() */
	//   i64 nTotalChange;             /* Value returned by sqlite3_total_changes() */
	//   int aLimit[SQLITE_N_LIMIT];   /* Limits */
	//   int nMaxSorterMmap;           /* Maximum size of regions mapped by sorter */
	init struct { /* Information used during initialization */
		newTnum Pgno  /* Rootpage of table being initialized */
		iDb     uint8 /* Which db file is being initialized */
		busy    uint8 /* TRUE if currently initializing */
		// unsigned orphanTrigger : 1; /* Last statement is orphaned TEMP trigger */
		// unsigned imposterTable : 1; /* Building an imposter table */
		// unsigned reopenMemdb : 1;   /* ATTACH is really a reopen using MemDB */
		azInit []byte /* "type", "name", and "tbl_name" columns */
	}
	//   int nVdbeActive;              /* Number of VDBEs currently running */
	//   int nVdbeRead;                /* Number of active VDBEs that read or write */
	//   int nVdbeWrite;               /* Number of active VDBEs that read and write */
	//   int nVdbeExec;                /* Number of nested calls to VdbeExec() */
	//   int nVDestroy;                /* Number of active OP_VDestroy operations */
	//   int nExtension;               /* Number of loaded extensions */
	//   void **aExtension;            /* Array of shared library handles */
	//   union {
	//     void (*xLegacy)(void*,const char*);   /* mTrace==SQLITE_TRACE_LEGACY */
	//     int (*xV2)(u32,void*,void*,void*);    /* All other mTrace values */
	//   } trace;
	//   void *pTraceArg;                        /* Argument to the trace function */
	// #ifndef SQLITE_OMIT_DEPRECATED
	//   void (*xProfile)(void*,const char*,u64);  /* Profiling function */
	//   void *pProfileArg;                        /* Argument to profile function */
	// #endif
	//   void *pCommitArg;                 /* Argument to xCommitCallback() */
	//   int (*xCommitCallback)(void*);    /* Invoked at every commit. */
	//   void *pRollbackArg;               /* Argument to xRollbackCallback() */
	//   void (*xRollbackCallback)(void*); /* Invoked at every commit. */
	//   void *pUpdateArg;
	//   void (*xUpdateCallback)(void*,int, const char*,const char*,sqlite_int64);
	//   void *pAutovacPagesArg;           /* Client argument to autovac_pages */
	//   void (*xAutovacDestr)(void*);     /* Destructor for pAutovacPAgesArg */
	//   unsigned int (*xAutovacPages)(void*,const char*,u32,u32,u32);
	//   Parse *pParse;                /* Current parse */
	// #ifdef SQLITE_ENABLE_PREUPDATE_HOOK
	//   void *pPreUpdateArg;          /* First argument to xPreUpdateCallback */
	//   void (*xPreUpdateCallback)(   /* Registered using sqlite3_preupdate_hook() */
	//     void*,sqlite3*,int,char const*,char const*,sqlite3_int64,sqlite3_int64
	//   );
	//   PreUpdate *pPreUpdate;        /* Context for active pre-update callback */
	// #endif /* SQLITE_ENABLE_PREUPDATE_HOOK */
	// #ifndef SQLITE_OMIT_WAL
	//   int (*xWalCallback)(void *, sqlite3 *, const char *, int);
	//   void *pWalArg;
	// #endif
	//   void(*xCollNeeded)(void*,sqlite3*,int eTextRep,const char*);
	//   void(*xCollNeeded16)(void*,sqlite3*,int eTextRep,const void*);
	//   void *pCollNeededArg;
	//   sqlite3_value *pErr;          /* Most recent error message */
	//   union {
	//     volatile int isInterrupted; /* True if sqlite3_interrupt has been called */
	//     double notUsed1;            /* Spacer */
	//   } u1;
	//   Lookaside lookaside;          /* Lookaside malloc configuration */
	// #ifndef SQLITE_OMIT_AUTHORIZATION
	//   sqlite3_xauth xAuth;          /* Access authorization function */
	//   void *pAuthArg;               /* 1st argument to the access auth function */
	// #endif
	// #ifndef SQLITE_OMIT_PROGRESS_CALLBACK
	//   int (*xProgress)(void *);     /* The progress callback */
	//   void *pProgressArg;           /* Argument to the progress callback */
	//   unsigned nProgressOps;        /* Number of opcodes for progress callback */
	// #endif
	// #ifndef SQLITE_OMIT_VIRTUALTABLE
	//   int nVTrans;                  /* Allocated size of aVTrans */
	//   Hash aModule;                 /* populated by sqlite3_create_module() */
	//   VtabCtx *pVtabCtx;            /* Context for active vtab connect/create */
	//   VTable **aVTrans;             /* Virtual tables with open transactions */
	//   VTable *pDisconnect;          /* Disconnect these in next sqlite3_prepare() */
	// #endif
	//   Hash aFunc;                   /* Hash table of connection functions */
	//   Hash aCollSeq;                /* All collating sequences */
	//   BusyHandler busyHandler;      /* Busy callback */
	//   Db aDbStatic[2];              /* Static space for the 2 default backends */
	//   Savepoint *pSavepoint;        /* List of active savepoints */
	//   int nAnalysisLimit;           /* Number of index rows to ANALYZE */
	//   int busyTimeout;              /* Busy handler timeout, in msec */
	//   int nSavepoint;               /* Number of non-transaction savepoints */
	//   int nStatement;               /* Number of nested statement-transactions  */
	//   i64 nDeferredCons;            /* Net deferred constraints this transaction. */
	//   i64 nDeferredImmCons;         /* Net deferred immediate constraints */
	//   int *pnBytesFreed;            /* If not NULL, increment this in DbFree() */
	// #ifdef SQLITE_ENABLE_UNLOCK_NOTIFY
	//   /* The following variables are all protected by the STATIC_MAIN
	//   ** mutex, not by sqlite3.mutex. They are used by code in notify.c.
	//   **
	//   ** When X.pUnlockConnection==Y, that means that X is waiting for Y to
	//   ** unlock so that it can proceed.
	//   **
	//   ** When X.pBlockingConnection==Y, that means that something that X tried
	//   ** tried to do recently failed with an SQLITE_LOCKED error due to locks
	//   ** held by Y.
	//   */
	//   sqlite3 *pBlockingConnection; /* Connection that caused SQLITE_LOCKED */
	//   sqlite3 *pUnlockConnection;           /* Connection to watch for unlock */
	//   void *pUnlockArg;                     /* Argument to xUnlockNotify */
	//   void (*xUnlockNotify)(void **, int);  /* Unlock notify callback */
	//   sqlite3 *pNextBlocked;        /* Next in list of all blocked connections */
	// #endif
	// #ifdef SQLITE_USER_AUTHENTICATION
	//   sqlite3_userauth auth;        /* User authentication information */
	// #endif
}
