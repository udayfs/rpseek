package main

import (
	"os"

	"github.com/udayfs/rpseek/tools"
)

const (
	metaDataFilePath string = "./data/arxiv-metadata.json"
	indexFilePath    string = "./data/index.json"
)

func main() {
	err := tools.BuildJsonIndex(metaDataFilePath, indexFilePath)
	if err != nil {
		os.Exit(1)
	}
}
