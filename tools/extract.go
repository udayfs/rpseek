package tools

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/minio/simdjson-go"
)

type PaperMeta struct {
	Id         *simdjson.Element
	Title      *simdjson.Element
	Authors    *simdjson.Element
	Abstract   *simdjson.Element
	Categories *simdjson.Element
}

type IndexEntry struct {
	Id   string         `json:"id"`
	Freq map[string]int `json:"tfs"`
}

func tokenize(text string) []string {
	re := regexp.MustCompile(`\w+`)
	return re.FindAllString(strings.ToUpper(text), -1)
}

func buildTf(text string) map[string]int {
	terms := tokenize(text)

	tf := make(map[string]int)
	for _, term := range terms {
		tf[term]++
	}

	return tf
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

	writer := bufio.NewWriter(of)
	enc := json.NewEncoder(writer)

	defer jf.Close()
	defer of.Close()
	defer writer.Flush()

	resultChan := make(chan simdjson.Stream, 40)

	go func() {
		simdjson.ParseNDStream(jf, resultChan, nil)
	}()

	for r := range resultChan {
		if r.Error != nil {
			if r.Error == io.EOF {
				break
			}
			return r.Error
		}

		var meta PaperMeta

		err := r.Value.ForEach(func(i simdjson.Iter) error {
			var err error

			meta.Id, err = i.FindElement(nil, "id")
			if err != nil {
				return err
			}

			meta.Title, err = i.FindElement(nil, "title")
			if err != nil {
				return err
			}

			meta.Authors, err = i.FindElement(nil, "authors")
			if err != nil {
				return err
			}

			meta.Abstract, err = i.FindElement(nil, "abstract")
			if err != nil {
				return err
			}

			meta.Categories, err = i.FindElement(nil, "categories")
			if err != nil {
				return err
			}

			id, err := meta.Id.Iter.String()
			if err != nil {
				return err
			}

			title, err := meta.Title.Iter.String()
			if err != nil {
				return err
			}

			authors, err := meta.Authors.Iter.String()
			if err != nil {
				return err
			}

			abstract, err := meta.Abstract.Iter.String()
			if err != nil {
				return err
			}

			categories, err := meta.Categories.Iter.String()
			if err != nil {
				return err
			}

			text := title + " " + authors + " " + abstract + " " + categories
			tf := buildTf(text)

			entry := IndexEntry{Id: id, Freq: tf}
			if err := enc.Encode(&entry); err != nil {
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
