package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/minio/sha256-simd"
)

type MerkleTree struct {
	m map[int][]byte
}

func New() *MerkleTree {
	return &MerkleTree{
		m: make(map[int][]byte),
	}
}

func (tree *MerkleTree) leafHash(n []byte) []byte {
	h := sha256.Sum256(n)
	return h[:]
}

func (tree *MerkleTree) parentHash(l, r []byte) []byte {
	h := sha256.New()
	h.Write(l)
	h.Write(r)
	return h.Sum(nil)
}

func (tree *MerkleTree) foldr(f func([]byte, []byte) []byte, coll [][]byte) []byte {
	if len(coll) == 0 {
		return nil
	}
	res := coll[len(coll)-1]
	for i := len(coll) - 2; i >= 0; i-- {
		res = f(coll[i], res)
	}
	return res
}

func (tree *MerkleTree) insert(s map[int][]byte, v []byte, n int) map[int][]byte {
	if _, ok := s[n]; ok {
		p := tree.parentHash(s[n], v)
		return tree.insert(tree.del(s, n), p, n+1)
	}
	s[n] = v
	return s
}

func (tree *MerkleTree) del(s map[int][]byte, n int) map[int][]byte {
	delete(s, n)
	return s
}

func (tree *MerkleTree) Digest() []byte {
	var keys []int
	for k := range tree.m {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	var vals [][]byte
	for _, k := range keys {
		vals = append(vals, tree.m[k])
	}
	return tree.foldr(tree.parentHash, vals)
}

func (tree *MerkleTree) Root(stream [][]byte) []byte {
	tree.populate(stream)
	return tree.Digest()
}

func (tree *MerkleTree) populate(stream [][]byte) {
	if len(stream) == 0 {
		return
	}
	for _, v := range stream {
		tree.Add(v)
	}
}

func (tree *MerkleTree) Add(v []byte) {
	tree.m = tree.insert(tree.m, tree.leafHash(v), 0)
}

func main2342342323() {
	//iterations := 67108864
	// iterations := 33554432
	// iterations := 4194304
	/*
		blkstream := make([][]byte, iterations)
		for i := 0; i < iterations; i++ {
			blkstream[i] = []byte("42") //fmt.Sprint(i)
		}
	*/

	m := New()
	iterations := 8 << 30 / 256
	// blkstream := make([][]byte, iterations)
	buf := make([]byte, 128)
	for i := 0; i < iterations; i++ {
		_, err := rand.Read(buf)
		if err != nil {
			log.Fatalf("error while generating random string: %s", err)
		}
		//blkstream[i] = buf
		m.Add(buf)
	}

	start := time.Now()

	root := m.Digest()
	elapsed := time.Since(start)
	fmt.Printf("Merkle root: %x\n", root)
	fmt.Printf("Elapsed time: %s\n", elapsed)

}
