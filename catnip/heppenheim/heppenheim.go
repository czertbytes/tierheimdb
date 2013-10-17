package main

import (
	"os"

	c "github.com/czertbytes/tierheimdb/catnip"
)

func main() {
	c.NewCatnip(NewParser()).Run(os.Args)
}
