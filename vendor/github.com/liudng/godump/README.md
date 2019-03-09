# godump

Dumps information about a variable. Now godump is part of [zhgo](https://github.com/zhgo/console) project.

[![Build Status](https://travis-ci.org/liudng/godump.svg)](https://travis-ci.org/liudng/godump)
[![Coverage](http://gocover.io/_badge/github.com/liudng/godump)](http://gocover.io/github.com/liudng/godump)
[![GoDoc](https://godoc.org/github.com/liudng/godump?status.png)](http://godoc.org/github.com/liudng/godump)

## Install

```bash
go get github.com/liudng/godump
```

## Sample code

```go
package main

import (
	"github.com/liudng/godump"
)

func main() {
	a := make(map[string]int64)

	a["A"] = 1
	a["B"] = 2

	godump.Dump(a)
}
```

Then Print:

```bash
(map[string]int64)
  A(int64) 1
  B(int64) 2
```
