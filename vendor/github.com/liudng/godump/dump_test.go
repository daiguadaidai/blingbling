// Copyright 2014 The godump Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package godump

import (
	"go/parser"
	"go/token"
	"testing"
)

var emptyString = ""

type S struct {
	A int
	B int
}

type T struct {
	S
	C int
}

type Circular struct {
	c *Circular
}

func TestDump(t *testing.T) {
	Dump(nil)
	Dump(token.STRING)
	Dump(&emptyString)
	Dump([3]int{1, 2, 3})
	Dump([]int{1, 2, 3})
	Dump(&[][]int{[]int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}})
	Dump(map[string]int{"satu": 1, "dua": 2})
	Dump(T{S{1, 2}, 3})
	//Dump(T{A: 1, B: 2, C: 3})

	/*bulet := make([]Circular, 3)
	bulet[0].c = &bulet[1]
	bulet[1].c = &bulet[2]
	bulet[2].c = &bulet[0]
	Dump(struct{ a []Circular }{bulet})*/

	fset := token.NewFileSet() // positions are relative to fset
	file, err := parser.ParseFile(fset, "dump_test.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	Dump(file)
}
