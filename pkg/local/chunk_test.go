package local

import (
	"bytes"
	"testing"
)

func TestSaveChunk(t *testing.T) {
	s := &chunkStorage{}
	if err := s.SaveChunk("./test", 1, bytes.NewReader([]byte("1")), 0); err != nil {
		t.Fatal(err)
	}
	if err := s.SaveChunk("./test", 2, bytes.NewReader([]byte("2")), 0); err != nil {
		t.Fatal(err)
	}
	if err := s.SaveChunk("./test", 3, bytes.NewReader([]byte("3")), 0); err != nil {
		t.Fatal(err)
	}
	if path, err := s.MergeChunk("./test", len("123")); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("success upload, path: %v", path)
	}

}
