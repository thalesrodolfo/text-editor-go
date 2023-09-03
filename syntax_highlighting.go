package main

// keywords
const (
	PACKAGE = "package"
	IMPORT  = "import"
	FUNC    = "func"
	VAR     = "var"
	CONST   = "const"
	TYPE    = "type"
	STRUCT  = "struct"
	FOR     = "for"
	IF      = "if"
	ELSE    = "else"
	RANGE   = "range"
	RETURN  = "return"
)

// types
const (
	INT     = " int"
	INT32   = " int32"
	INT64   = " int64"
	STRING  = " string"
	CHAR    = " char"
	BYTE    = " byte"
	RUNE    = " rune"
	FLOAT   = " float"
	FLOAT32 = " float32"
	FLOAT64 = " float64"
	ERROR   = " error"
	NIL     = " nil"
)

type KeywordPos struct {
	line      int
	init, end int
}

func getKeywords() [12]string {
	return [12]string{PACKAGE, IMPORT, FUNC, VAR, CONST, TYPE, STRUCT, FOR, IF, ELSE, RANGE, RETURN}
}

func getTypes() [12]string {
	return [12]string{INT, INT32, INT64, STRING, CHAR, BYTE, RUNE, FLOAT, FLOAT32, FLOAT64, ERROR, NIL}
}

/*

   {
   package: [(0, 8)]
   }
*/
