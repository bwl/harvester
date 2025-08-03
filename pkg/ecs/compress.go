package ecs

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/json"
)

type SaveOptions struct {
	Password string
	Compress bool
}

type encodedBlob struct {
	Raw []byte `json:"raw"`
}

func EncodeSnapshot(s *Snapshot, opt SaveOptions) ([]byte, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	if opt.Compress {
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		if _, err := zw.Write(b); err != nil {
			return nil, err
		}
		if err := zw.Close(); err != nil {
			return nil, err
		}
		b = buf.Bytes()
	}
	if opt.Password != "" {
		key := sha256.Sum256([]byte(opt.Password))
		blk, err := aes.NewCipher(key[:])
		if err != nil {
			return nil, err
		}
		iv := make([]byte, aes.BlockSize)
		stream := cipher.NewCTR(blk, iv)
		ct := make([]byte, len(b))
		stream.XORKeyStream(ct, b)
		b = ct
	}
	return b, nil
}

func DecodeSnapshot(b []byte, opt SaveOptions) (*Snapshot, error) {
	if opt.Password != "" {
		key := sha256.Sum256([]byte(opt.Password))
		blk, err := aes.NewCipher(key[:])
		if err != nil {
			return nil, err
		}
		iv := make([]byte, aes.BlockSize)
		stream := cipher.NewCTR(blk, iv)
		pt := make([]byte, len(b))
		stream.XORKeyStream(pt, b)
		b = pt
	}
	if opt.Compress {
		zr, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		var out bytes.Buffer
		if _, err := out.ReadFrom(zr); err != nil {
			return nil, err
		}
		if err := zr.Close(); err != nil {
			return nil, err
		}
		b = out.Bytes()
	}
	var s Snapshot
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return &s, nil
}
