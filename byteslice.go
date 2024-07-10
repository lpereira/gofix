// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
)

func init() {
	register(byteslicestrOpt)
}

var byteslicestrOpt = fix{
	name: "byteslicestr",
	date: "2024-06-04",
	f:    byteslicestr,
	desc: `Convert []byte(str) to something tinygo likes better.`,
}

func byteslicestr(f *ast.File) bool {
	fixed := false
	walk(f, func(n any) {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		arrayTy, ok := call.Fun.(*ast.ArrayType)
		if !ok {
			return
		}
		if arrayTy.Len != nil {
			return
		}

		eltTy, ok := arrayTy.Elt.(*ast.Ident)
		if !ok {
			return
		}

		if eltTy.Name != "byte" {
			return
		}

		if len(call.Args) != 1 {
			return
		}

		// Only transform basic cases, i.e. []byte(s), and not cases where
		// TinyGo may already have optimizations (e.g. []byte("some literal")),
		// or cases where the argument would be the result of a function call,
		// as we don't know if that has side-effects.
		switch arg := call.Args[0].(type) {
		case *ast.Ident, *ast.SelectorExpr:
			call.Fun = ast.NewIdent("unsafe.Slice")
			call.Args = []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("unsafe.StringData"),
					Args: []ast.Expr{
						arg,
					},
				},
				&ast.CallExpr{
					Fun: ast.NewIdent("len"),
					Args: []ast.Expr{
						arg,
					},
				},
			}
		default:
			return
		}

		fixed = true
	})
	return fixed
}
