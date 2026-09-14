package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := "test-pass"
	cases := [][]byte{
		[]byte("hello"),
		bytes.Repeat([]byte{0xab}, chunkSize),
		bytes.Repeat([]byte{0xcd}, chunkSize+17),
		{},
	}
	for i, in := range cases {
		got := decrypt(encrypt(in, key), key)
		if !bytes.Equal(got, in) {
			t.Fatalf("case %d: round trip mismatch, in=%d got=%d", i, len(in), len(got))
		}
	}
}

func TestFileRoundTripWithSplit(t *testing.T) {
	dir := t.TempDir()
	inDir := filepath.Join(dir, "in")
	outDir := filepath.Join(dir, "out")
	backDir := filepath.Join(dir, "back")
	if err := os.MkdirAll(filepath.Join(inDir, "sub"), 0755); err != nil {
		t.Fatal(err)
	}

	raw := bytes.Repeat([]byte("0123456789abcdef"), chunkSize/8+100)
	src := filepath.Join(inDir, "sub", "demo.bin")
	if err := os.WriteFile(src, raw, 0644); err != nil {
		t.Fatal(err)
	}

	p := &Progress{}
	if err := dojob(p, inDir, outDir, true, "secret", chunkSize*2); err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	enc := filepath.Join(outDir, "sub", "demo.bin"+encSuffix)
	if _, err := os.Stat(enc); err != nil {
		t.Fatalf("missing encrypted file: %v", err)
	}
	if _, err := os.Stat(enc + ".0"); err != nil {
		t.Fatalf("expected split part .0: %v", err)
	}

	p2 := &Progress{}
	if err := dojob(p2, outDir, backDir, false, "secret", chunkSize*2); err != nil {
		t.Fatalf("decrypt: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(backDir, "sub", "demo.bin"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, raw) {
		t.Fatalf("decrypted content mismatch: %d vs %d", len(got), len(raw))
	}
}
