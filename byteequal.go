// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
	"go/token"
)

func init() {
	register(stringequalOpt)
}

var stringequalOpt = fix{
	name: "stringequal",
	date: "2024-06-05",
	f:    stringequal,
	desc: `Replaces string(a) == string(b) with a version that doesn't allocate.`,
}

func stringequal(f *ast.File) bool {
	return false


	// This is broken!
	fixed := false
	walk(f, func(n any) {
		binOp, ok := n.(*ast.BinaryExpr)
		if !ok {
			return
		}

		if binOp.Op != token.EQL {
			return
		}

		binOp.X = &ast.BinaryExpr{
			X: &ast.CallExpr{
				Fun: ast.NewIdent("len"),
				Args: []ast.Expr{
					binOp.X,
				},
			},
			Op: token.EQL,
			Y: &ast.CallExpr{
				Fun: ast.NewIdent("len"),
				Args: []ast.Expr{
					binOp.Y,
				},
			},
		}
		binOp.Op = token.LAND;
		binOp.Y = &ast.CallExpr{
			Fun: ast.NewIdent("memequal"),
			Args: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("unsafe.SliceData"),
					Args: []ast.Expr{
						binOp.X,
					},
				},
				&ast.CallExpr{
					Fun: ast.NewIdent("unsafe.SliceData"),
					Args: []ast.Expr{
						binOp.Y,
					},
				},
				&ast.CallExpr{
					Fun: ast.NewIdent("len"),
					Args: []ast.Expr{
						binOp.X,
					},
				},
			},
		}

		fixed = true
	})
	return fixed
}
