package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dominichamon/golisp/golisp"
)

var (
	verbose = flag.Bool("verbose", false, "enable to get verbose logging")
)

func main() {
	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(ioutil.Discard)
	if *verbose {
		log.SetOutput(os.Stdout)
	}

	if len(flag.Args()) == 0 {
		if err := golisp.Repl(); err != nil {
			log.Fatalf("%s", err)
		}
	}

	res, err := golisp.Exec(flag.Arg(0))
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Printf("%s\n", res)
}
