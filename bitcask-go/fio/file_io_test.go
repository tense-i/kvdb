package fio

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestNewFileIOManager(t *testing.T) {
	fio, err := NewFileIOManager(filepath.Join("./../tmp", "a.data"))
	assert.Nil(t, err)
	assert.NotNil(t, fio)
}

func TestFileIO_Write(t *testing.T) {
	fio, err := NewFileIOManager(filepath.Join("./../tmp", "a.data"))
	assert.Nil(t, err)
	assert.NotNil(t, fio)
	n, err := fio.Write([]byte(""))
	assert.Equal(t, 0, n)
	n, err = fio.Write([]byte("data"))
	assert.Equal(t, len("data"), n)

}

func TestFileIO_Read(t *testing.T) {
	fio, err := NewFileIOManager(filepath.Join("./../tmp", "c.data"))
	assert.Nil(t, err)
	assert.NotNil(t, fio)
	dat1 := make([]byte, 6)
	n, err := fio.Write([]byte("abcde"))
	assert.Equal(t, len("abcde"), n)
	n, err = fio.Write([]byte("fghij"))
	assert.Equal(t, len("fghij"), n)
	_, err = fio.Read(dat1, 0)
	t.Log(dat1)
	assert.Nil(t, err)
	dat2 := make([]byte, 6)
	_, err = fio.Read(dat2, 3)
	assert.Nil(t, err)
	t.Log(dat2)

}
