package common

import (
	"testing"
	"fmt"
)

func TestStrIsMatch(t *testing.T) {
	// reg := fmt.Sprintf(`(?i)\s*CREATE\s*TABLE\s*[%v\w\d_]+\s*LIKE\s*[%v\w\d_;]+`, "`", "`")
	reg := fmt.Sprintf(`(?i)^\s*CREATE\s*TABLE\s*[0-9a-z_%s\.]+\s*LIKE\s*[0-9a-z_%s\.]+\s*;?\s*$`, "`")
	fmt.Println(reg)
	sql := "create table `order_invoice_1`.`invoice_detail__1023` like order_invoice_0.invoice_detail_0;   "
	fmt.Println(sql)

	matchd := StrIsMatch(sql, reg)

	fmt.Println(matchd)
}
