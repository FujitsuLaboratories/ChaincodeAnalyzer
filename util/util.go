package util

import (
	"go/ast"
	"go/token"
)

func SrcLine(src []byte, p token.Position) string {
	// run to end of line in both directions if not at line start/end
	lo, hi := p.Offset, p.Offset+1
	for lo > 0 && src[lo-1] != '\n' {
		lo--
	}
	for hi < len(src) && src[hi-1] != '\n' {
		hi++
	}
	return string(src[lo:hi])
}

func IdName(n ast.Node) string {
	var name string
	if n != nil {
		switch t := n.(type) {
		case *ast.FuncDecl:
			name = t.Name.Name
		case *ast.ImportSpec:
			name = t.Name.Name
		case *ast.BasicLit:
			name = t.Value
		case *ast.Ident:
			name = t.Name
		case *ast.IndexExpr:
			name = IdName(t.X) + "[" + IdName(t.Index) + "]"
		case *ast.SelectorExpr:
			name = IdName(t.X) + "." + IdName(t.Sel)
		default:
			return ""
		}
	}
	return name
}
