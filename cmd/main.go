package main

import (
	"flag"
	"fmt"
	"github.com/first-task/internal/cvt"
	"log"
)

var i, o, f, t string

func init() {
	flag.StringVar(&i ,"i", "", "input dir")
	flag.StringVar(&o ,"o", "", "output dir")
	flag.StringVar(&f, "f", "jpg", "input ext")
	flag.StringVar(&t, "t", "png", "output ext")
	flag.Parse()
}

func main() {
	c := cvt.NewImageCvtGlue(i, o, f, t)
	if err := c.Exec(); err != nil {
		log.Fatalf("Failed to execute image convert", fmt.Sprintf("%+v", err))
	}
}
