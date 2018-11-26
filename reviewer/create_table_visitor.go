package reviewer

import (
	"fmt"

	"github.com/daiguadaidai/blingbling/ast"
)

type CreateTableVisitor struct {
}

func NewCreateTableVisitor() *CreateTableVisitor {
	createTableVisitor := new(CreateTableVisitor)

	return createTableVisitor
}

func (this *CreateTableVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	fmt.Printf("Enter: %T, %[1]v\n", in)

	return in, false
}

func (this *CreateTableVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	fmt.Printf("Leave: %T\n", in)

	return in, true
}
