package main

import (
	"git.esentri.com/fn/ext-zeebe/zeebe"
)

func main() {
	zeebe.ListFunctions("http://localhost:8080")
}