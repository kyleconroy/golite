package internal

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
	// sqlite3 *db;         /* The main database structure */
	// char *zErrMsg;       /* An error message */
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
	//   Token constraintName;/* Name of the constraint currently being parsed */
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
	//   Token sLastToken;       /* The last token parsed */
	//   ynVar nVar;               /* Number of '?' variables seen in the SQL so far */
	//   u8 iPkSortOrder;          /* ASC or DESC for INTEGER PRIMARY KEY */
	//   u8 explain;               /* True if the EXPLAIN flag is found on the query */
	//   u8 eParseMode;            /* PARSE_MODE_XXX constant */
	// // #ifndef SQLITE_OMIT_VIRTUALTABLE
	//   int nVtabLock;            /* Number of virtual tables to lock */
	// // #endif
	//   int nHeight;              /* Expression tree height of current sub-select */
	// // #ifndef SQLITE_OMIT_EXPLAIN
	//   int addrExplain;          /* Address of current OP_Explain opcode */
	// // j#endif
	//   VList *pVList;            /* Mapping between variable names and numbers */
	//   Vdbe *pReprepare;         /* VM being reprepared (sqlite3Reprepare()) */
	//   const char *zTail;        /* All SQL text past the last semicolon parsed */
	//   Table *pNewTable;         /* A table being constructed by CREATE TABLE */
	//   Index *pNewIndex;         /* An index being constructed by CREATE INDEX.
	//                             ** Also used to hold redundant UNIQUE constraints
	//                             ** during a RENAME COLUMN */
	//   Trigger *pNewTrigger;     /* Trigger under construct by a CREATE TRIGGER */
	//   const char *zAuthContext; /* The 6th parameter to db->xAuth callbacks */
	// #ifndef SQLITE_OMIT_VIRTUALTABLE
	//   Token sArg;               /* Complete text of a module argument */
	//   Table **apVtabLock;       /* Pointer to virtual tables needing locking */
	// #endif
	//   With *pWith;              /* Current WITH clause, or NULL */
	// #ifndef SQLITE_OMIT_ALTERTABLE
	//   RenameToken *pRename;     /* Tokens subject to renaming by ALTER TABLE */
	// #endif
}
