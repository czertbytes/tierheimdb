package main

import (
	"log"
	"os"

	c "github.com/czertbytes/tierheimdb/catnip"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	c.NewCatnip(NewParser()).Run(os.Args)
}
