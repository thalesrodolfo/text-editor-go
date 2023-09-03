package main

const (
	PACKAGE = "package"
	IMPORT  = "import"
	FUNC    = "func"
	VAR     = "var"
	TYPE    = "type"
	STRUCT  = "struct"
	FOR     = "for"
	IF      = "if"
	ELSE    = "else"
	RANGE   = "range"
)

type KeywordPos struct {
	line      int
	init, end int
}

func getKeywords() [10]string {
	return [10]string{PACKAGE, IMPORT, FUNC, VAR, TYPE, STRUCT, FOR, IF, ELSE, RANGE}
}

/*

   {
   package: [(0, 8)]
   }
*/
