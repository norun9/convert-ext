package main

import (
	"flag"
	"fmt"
	"github.com/first-task/internal/cvt"
	"log"
)

var i, o, be, ae string

func init() {
	flag.StringVar(&i, "i", "", "input dir")
	flag.StringVar(&o, "o", "", "output dir")
	flag.StringVar(&be, "be", "jpg", "before ext")
	flag.StringVar(&ae, "ae", "png", "after ext")
	flag.Parse()
}

func main() {
	c := cvt.NewImageCvtGlue(i, o, be, ae)
	if err := c.Exec(); err != nil {
		log.Fatalf("Failed to execute image convert", fmt.Sprintf("%+v", err))
	}
}
