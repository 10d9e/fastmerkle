package main

import (
	"crypto/rand"
	"log"
	"reflect"
	"testing"
)

func BenchmarkMerkleRoot(b *testing.B) {
	// equivalent hashes of a 8G file, chunked into 256 byte tranches

	iterations := 8 << 30 / 256
	blkstream := make([][]byte, iterations)
	buf := make([]byte, 128)
	for i := 0; i < iterations; i++ {
		_, err := rand.Read(buf)
		if err != nil {
			log.Fatalf("error while generating random string: %s", err)
		}
		blkstream[i] = buf
	}

	m := New()
	_ = m.Root(blkstream)
}

func TestMerkleRoot(t *testing.T) {
	type args struct {
		stream [][]byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test 1",
			args: args{
				stream: [][]byte{[]byte("a"), []byte("b"), []byte("c")},
			},
			want: []byte{0x70, 0x75, 0x15, 0x2d, 0x03, 0xa5, 0xcd, 0x92, 0x10, 0x48, 0x87, 0xb4, 0x76, 0x86, 0x27, 0x78, 0xec, 0x0c, 0x87, 0xbe, 0x5c, 0x2f, 0xa1, 0xc0, 0xa9, 0x0f, 0x87, 0xc4, 0x9f, 0xad, 0x6e, 0xff},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			if got := m.Root(tt.args.stream); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("merkleRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerkleRootWithAdd(t *testing.T) {
	type args struct {
		stream [][]byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test 1",
			args: args{
				stream: [][]byte{[]byte("a"), []byte("b"), []byte("c")},
			},
			want: []byte{0x70, 0x75, 0x15, 0x2d, 0x03, 0xa5, 0xcd, 0x92, 0x10, 0x48, 0x87, 0xb4, 0x76, 0x86, 0x27, 0x78, 0xec, 0x0c, 0x87, 0xbe, 0x5c, 0x2f, 0xa1, 0xc0, 0xa9, 0x0f, 0x87, 0xc4, 0x9f, 0xad, 0x6e, 0xff},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			for i := 0; i < len(tt.args.stream); i++ {
				m.Add(tt.args.stream[i])
			}

			if got := m.Digest(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("merkleRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}
