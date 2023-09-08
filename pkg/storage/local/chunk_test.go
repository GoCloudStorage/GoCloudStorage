package local

import (
	"bytes"
	"testing"
)

func TestSaveChunk(t *testing.T) {
	s := &chunkUploader{}
	if err := s.saveChunk("./test", 1, bytes.NewReader([]byte("1"))); err != nil {
		t.Fatal(err)
	}
	if err := s.saveChunk("./test", 2, bytes.NewReader([]byte("2"))); err != nil {
		t.Fatal(err)
	}
	if err := s.saveChunk("./test", 3, bytes.NewReader([]byte("3"))); err != nil {
		t.Fatal(err)
	}
	if err := s.mergeChunk("./test", 3, len("123")); err != nil {
		t.Fatal(err)
	}
}
