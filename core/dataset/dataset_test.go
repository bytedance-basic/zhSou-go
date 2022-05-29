package dataset

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

var mockIndexFiles = []IndexFile{
	{
		ItemLength: 1,
		SeekInfo: []SeekInfo{
			{0, 1, 4},
		},
	},
	{
		ItemLength: 2,
		SeekInfo: []SeekInfo{
			{0, 4, 8},
			{12, 4, 6},
		},
	},
	{
		ItemLength: 3,
		SeekInfo: []SeekInfo{
			{0, 4, 8},
			{12, 4, 6},
			{22, 3, 2},
		},
	},
}

func buildTestIndexFile() []io.Reader {

	var readers []io.Reader

	for _, file := range mockIndexFiles {
		buf := bytes.NewBuffer([]byte{})
		_ = gob.NewEncoder(buf).Encode(file)
		readers = append(readers, buf)
	}

	return readers
}
func TestNewIndexFileSet(t *testing.T) {
	rs := buildTestIndexFile()
	indexFileSet := NewIndexFileSet(rs)
	assert.NotNil(t, indexFileSet)
}

func TestIndexFileSet_Get(t *testing.T) {
	rs := buildTestIndexFile()
	indexFileSet := NewIndexFileSet(rs)
	// 判断原结构体与反序列化后的结构体是否相等
	assert.Equal(t, mockIndexFiles, indexFileSet.indexFiles)
	// 判断计算的前缀和是否一致
	assert.Equal(t, []uint32{1, 3, 6}, indexFileSet.idArray)
	// 共有6个数据索引记录

	{
		fileId, indexItem := indexFileSet.Get(0)
		assert.Equal(t, 0, fileId)
		assert.Equal(t, mockIndexFiles[0].SeekInfo[0], *indexItem)
	}
	{
		fileId, indexItem := indexFileSet.Get(1)
		assert.Equal(t, 1, fileId)
		assert.Equal(t, mockIndexFiles[1].SeekInfo[0], *indexItem)
	}
	{
		fileId, indexItem := indexFileSet.Get(2)
		assert.Equal(t, 1, fileId)
		assert.Equal(t, mockIndexFiles[1].SeekInfo[1], *indexItem)
	}
	{
		fileId, indexItem := indexFileSet.Get(3)
		assert.Equal(t, 2, fileId)
		assert.Equal(t, mockIndexFiles[2].SeekInfo[0], *indexItem)
	}
	{
		fileId, indexItem := indexFileSet.Get(4)
		assert.Equal(t, 2, fileId)
		assert.Equal(t, mockIndexFiles[2].SeekInfo[1], *indexItem)
	}
	{
		fileId, indexItem := indexFileSet.Get(5)
		assert.Equal(t, 2, fileId)
		assert.Equal(t, mockIndexFiles[2].SeekInfo[2], *indexItem)
	}
}

func TestNewDataReader(t *testing.T) {
	var indexFilePaths, dataFilePaths []string
	for i := 0; i < 256; i++ {
		indexFilePaths = append(indexFilePaths, fmt.Sprintf("D:\\index\\wukong_100m_%d.gob", i))
		dataFilePaths = append(dataFilePaths, fmt.Sprintf("D:\\after\\wukong_100m_%d.dat", i))
	}
	dataReader, err := NewDataReader(indexFilePaths, dataFilePaths)
	fmt.Println(dataReader.indexFileSet.idArray)
	if err != nil {
		panic(err)
	}
	// 按1kw为间隔遍历一亿条数据
	for i := 0; i <= 100000000; i += 10000000 {
		dataRecord, err := dataReader.Read(uint32(i))
		if err != nil {
			panic(fmt.Sprintf("读取错误：%v", err))
		}
		fmt.Println(i, dataRecord)
	}
	dr, _ := dataReader.Read(0)
	assert.Equal(t, "今年能跑赢96不?备战坦克两项俄军开始选拔参赛队员", dr.Text)
	dr, _ = dataReader.Read(100000000)
	assert.Equal(t, "光大银行福州分行成立二十周年发展纪实", dr.Text)
}
