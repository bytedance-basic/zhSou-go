package invertedindex

import (
	"encoding/gob"
	"github.com/zhSou/zhSou-go/core/dict"
	"github.com/zhSou/zhSou-go/util/algorithm/set"
	"io"
	"sort"
)

type invertedIndex struct {
	Data map[int][]int
	dict *dict.Dict
}

func NewInvertedIndex(dict *dict.Dict) *invertedIndex {
	return &invertedIndex{
		Data: make(map[int][]int),
		dict: dict,
	}
}

func LoadInvertedIndexFromDisk(r io.Reader) (*invertedIndex, error) {
	ii := invertedIndex{}
	err := gob.NewDecoder(r).Decode(&ii)
	if err != nil {
		return nil, err
	}
	return &ii, nil
}

func (i *invertedIndex) SaveToDisk(w io.Writer) error {
	err := gob.NewEncoder(w).Encode(*i)
	if err != nil {
		return err
	}
	return nil
}

func (i *invertedIndex) Add(word string, id int) {
	wordId := i.dict.Put(word)
	i.Data[wordId] = append(i.Data[wordId], id)
}

func (i *invertedIndex) AddWords(words []string, id int) {
	for _, word := range set.Deduplication[string](words) {
		wordId := i.dict.Put(word)
		i.Data[wordId] = append(i.Data[wordId], id)
	}
}

func (i *invertedIndex) Get(word string) []int {
	wordId := i.dict.Put(word)

	return i.Data[wordId]
}

func (i *invertedIndex) Sort() {
	for _, ids := range i.Data {
		sort.Ints(ids)
	}
}
