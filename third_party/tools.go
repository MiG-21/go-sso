//go:build third_party
// +build third_party

// See https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md
// for some notes on this file

package third_party

import (
	_ "github.com/onsi/ginkgo/ginkgo"
	_ "github.com/swaggo/swag/cmd/swag"
)
