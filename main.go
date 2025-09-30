package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/udayfs/rpseek/tools"
)

func main() {
	var err error

	flag.Usage = func() {
		usage :=
			"Usage: rpseek [command] [flags]\n" +
				"Available commands:\n" +
				" + index\n" +
				" + serve\n" +
				"Use '<command> -h' for command-specific usage"

		fmt.Fprintln(os.Stderr, usage)
	}

	indexCmd := flag.NewFlagSet("index", flag.ExitOnError)
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)

	indexCmd.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: index -i <file> [-o outPath]")
		indexCmd.PrintDefaults()
	}

	serveCmd.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: serve [-p port]")
		serveCmd.PrintDefaults()
	}

	var indexOutputPath string
	var indexInputPath string
	var port int

	indexCmd.StringVar(&indexInputPath, "i", "", "file to be indexed")
	indexCmd.StringVar(&indexOutputPath, "o", "index.json", "path to put the indexed output")
	serveCmd.IntVar(&port, "port", 42069, "port to be used for serving the web interface")

	if len(os.Args) > 1 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			flag.Usage()
			os.Exit(0)
		}

		switch os.Args[1] {
		case "index":
			indexCmd.Parse(os.Args[2:])
			err = tools.BuildJsonIndex(indexInputPath, indexOutputPath)

		case "serve":
			serveCmd.Parse(os.Args[2:])
		}
	} else {
		flag.Usage()
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
}
