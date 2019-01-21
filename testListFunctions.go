package main

import (
	"git.esentri.com/fn/ext-zeebe/zeebe"
)

func main() {
	zeebe.GetFunctionsWithZeebeJobType("http://localhost:8080")
}
