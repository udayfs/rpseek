package internal

import (
	"fmt"
	"io"
	"os"

	"github.com/minio/simdjson-go"
)

type FreqTable map[string]int

type SearchQuery struct {
	Query string `json:"query"`
}

// See: https://en.wikipedia.org/wiki/Tf%E2%80%93idf#Term_frequency
func BuildTermFreq(term string, docFreq FreqTable) float64 {
	var sum int

	for _, val := range docFreq {
		sum += val
	}

	return float64(docFreq[term]) / float64(sum)
}

func SearchDoc(indexFilePath string, tokens []string) error {
	index, err := os.Open(indexFilePath)
	if err != nil {
		return err
	}

	defer index.Close()

	parsed := make(chan simdjson.Stream, 1)

	go func() {
		simdjson.ParseNDStream(index, parsed, nil)
	}()

	for r := range parsed {
		if r.Error != nil {
			if r.Error == io.EOF {
				break
			}
			return err
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

			freq, err := obj.FindKey("freq", nil).Iter.Object(nil)
			if err != nil {
				return err
			}

			tf_freq := make(FreqTable)
			err = freq.ForEach(func(key []byte, i simdjson.Iter) {
				val, err := i.Int()
				if err != nil {
					return
				}

				tf_freq[string(key)] = int(val)

			}, nil)

			var total_tf float64 = 0.0
			for _, tok := range tokens {
				total_tf += BuildTermFreq(tok, tf_freq)
			}

			fmt.Println(id, "=>", total_tf)
	
			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}
