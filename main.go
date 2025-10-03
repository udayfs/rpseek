package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/udayfs/rpseek/internal"
	"github.com/udayfs/rpseek/tools"
)

var (
	indexOutputPath string
	indexInputPath  string
	serveIndexPath  string
	port            int
	host            string
)

func main() {
	var err error

	flag.Usage = func() {
		usage :=
			"Usage: rpseek [command] [flags]\n" +
				"Available commands:\n" +
				" > index\n" +
				" > serve\n" +
				"Use '<command> -h' for command specific usage"

		fmt.Fprintln(os.Stderr, usage)
	}

	indexCmd := flag.NewFlagSet("index", flag.ExitOnError)
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)

	indexCmd.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: index -i <file> [-o <outPath>]")
		indexCmd.PrintDefaults()
	}

	serveCmd.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: serve -f <index-file> [-port <port>] [-host <host>]")
		serveCmd.PrintDefaults()
	}

	indexCmd.StringVar(&indexInputPath, "i", "", "file to be indexed")
	indexCmd.StringVar(&indexOutputPath, "o", "index.json", "path to put the indexed output")

	serveCmd.StringVar(&serveIndexPath, "f", "", "index file to be searched upon")
	serveCmd.StringVar(&host, "host", "127.0.0.1", "web server host")
	serveCmd.IntVar(&port, "port", 42069, "web server port")

	if len(os.Args) > 1 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			flag.Usage()
			os.Exit(0)
		}

		switch os.Args[1] {
		case "index":
			indexCmd.Parse(os.Args[2:])
			if indexInputPath == "" {
				internal.ExitOnError("no input file specified with '-i' flag")
			}
			err = tools.BuildJsonIndex(indexInputPath, indexOutputPath)

		case "serve":
			serveCmd.Parse(os.Args[2:])
			if serveIndexPath == "" {
				internal.ExitOnError("no index file specified with '-f' flag")
			}
			server := &tools.Server{Host: host, Port: port, IndexFile: serveIndexPath}
			err = server.Serve()

		default:
			internal.ExitOnError(fmt.Sprintf("unrecognized command '%s'", os.Args[1]))
		}
	} else {
		flag.Usage()
	}

	if err != nil {
		internal.ExitOnError(err.Error())
	}
}
