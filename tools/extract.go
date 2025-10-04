package tools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/minio/simdjson-go"
	"github.com/udayfs/rpseek/internal"
)

type indexEntry struct {
	Id   string             `json:"id"`
	Freq internal.FreqTable `json:"freq"`
}

func buildDocFreq(text string) map[string]int {
	terms := internal.Tokenize(text)

	df := make(map[string]int)
	for _, term := range terms {
		df[term]++
	}

	return df
}

func BuildJsonIndex(metaDataFilePath string, indexFilePath string) error {
	if !simdjson.SupportedCPU() {
		return fmt.Errorf("simdjson: unsupported CPU")
	}

	jf, err := os.Open(metaDataFilePath)
	if err != nil {
		return err
	}

	of, err := os.Create(indexFilePath)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(jf)
	writer := bufio.NewWriter(of)
	enc := json.NewEncoder(writer)

	defer jf.Close()
	defer of.Close()
	defer writer.Flush()

	resultChan := make(chan simdjson.Stream, 1)

	go func() {
		simdjson.ParseNDStream(reader, resultChan, nil)
	}()

	for r := range resultChan {
		if r.Error != nil {
			if r.Error == io.EOF {
				break
			}
			return r.Error
		}

		err := r.Value.ForEach(func(i simdjson.Iter) error {
			var err error

			obj, err := i.Object(nil)
			if err != nil {
				return err
			}

			id, err := obj.FindKey("id", nil).Iter.String()
			if err != nil {
				return err
			}

			title, err := obj.FindKey("title", nil).Iter.String()
			if err != nil {
				return err
			}

			authors, err := obj.FindKey("authors", nil).Iter.String()
			if err != nil {
				return err
			}

			abstract, err := obj.FindKey("abstract", nil).Iter.String()
			if err != nil {
				return err
			}

			categories, err := obj.FindKey("categories", nil).Iter.String()
			if err != nil {
				return err
			}

			text := strings.Join([]string{title, authors, abstract, categories}, " ")
			df := buildDocFreq(text)

			if err := enc.Encode(&indexEntry{Id: id, Freq: df}); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
