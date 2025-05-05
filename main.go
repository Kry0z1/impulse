package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/Kry0z1/impulse/app"
	"os"

	"github.com/Kry0z1/impulse/config"
	"github.com/Kry0z1/impulse/fs"
	"github.com/Kry0z1/impulse/lib"
)

var (
	cfgPath = flag.String("c", "", "path to config file")
	inPath  = flag.String("i", "", "path to input file")
	logPath = flag.String("l", "", "path to log file")
	outPath = flag.String("o", "", "path to output file")
)

func main() {
	flag.Parse()

	cfg := config.MustLoad(*cfgPath)

	in := fs.MustLoadInputFile(*inPath)
	out, err := fs.OpenOutputFile(*outPath)
	if err != nil {
		fmt.Println(err.Error())
		out = bufio.NewWriter(os.Stdout)
	}
	log, err := fs.OpenLogFile(*logPath)
	if err != nil {
		fmt.Println(err.Error())
		log = bufio.NewWriter(os.Stdout)
	}

	orch := lib.NewOrchestrator(cfg)

	application := app.New(in, out, log, cfg, orch)
	err = application.Run()

	if err != nil {
		fmt.Println(err.Error())
	}
}
