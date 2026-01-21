package main

import (
	"flag"
	"log"
)

func main() {
	n := flag.Int("n", 3, "number of overwrite passes")
	z := flag.Bool("z", false, "add final zero overwrite")
	u := flag.Bool("u", false, "remove file after overwriting")
	v := flag.Bool("v", false, "verbose output")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("usage: shred [options] file")
	}

	if err := shredPath(flag.Arg(0), *n, *z, *u, *v); err != nil {
		log.Fatal(err)
	}
}

