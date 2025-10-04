package internal

import (
	"bufio"
	"io"
	"math"
	"os"
	"sort"

	"github.com/minio/simdjson-go"
)

type FreqTable map[string]int

type SearchQuery struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Doc_id string  `json:"doc_id"`
	Rank   float64 `json:"rank"`
}

// See: https://en.wikipedia.org/wiki/Tf%E2%80%93idf#Term_frequency
func BuildTermFreq(term string, docFreq FreqTable) float64 {
	var sum int

	for _, val := range docFreq {
		sum += val
	}

	return float64(docFreq[term]) / float64(sum)
}

func BuildInverseDocFreq(term string, docs map[string]FreqTable) float64 {
	N := len(docs)
	M := 0

	for _, tf := range docs {
		if tf[term] != 0 {
			M += 1
		}
	}

	return math.Log10(float64(N) / (1.0 + float64(M)))
}

func SearchDoc(indexFilePath string, tokens []string) ([]QueryResponse, error) {
	index, err := os.Open(indexFilePath)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(index)
	defer index.Close()

	parsed := make(chan simdjson.Stream, 1)
	docs := make(map[string]FreqTable)

	go func() {
		simdjson.ParseNDStream(reader, parsed, nil)
	}()

	for r := range parsed {
		if r.Error != nil {
			if r.Error == io.EOF {
				break
			}
			return nil, err
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

			docs[id] = tf_freq

			return err
		})

		if err != nil {
			return nil, err
		}
	}

	res := make([]QueryResponse, 0)
	for id, freq := range docs {
		var rank float64 = 0.0
		for _, tok := range tokens {
			rank += BuildTermFreq(tok, freq) * BuildInverseDocFreq(tok, docs)
		}
		res = append(res, QueryResponse{id, rank})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Rank > res[j].Rank
	})

	return res, nil
}
