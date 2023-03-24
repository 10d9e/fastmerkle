package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
)

type proof struct {
	pre  []subRoot
	leaf []byte
	post []subRoot
}

type subRoot struct {
	i    int
	subr []byte
}

func limit(stream [][]byte, n int) [][]byte {
	r := make([][]byte, 0)
	for _, s := range stream {
		if n == 0 {
			break
		}
		r = append(r, []byte(s))
		n--
	}
	return r
}

func subroot(stream [][]byte, k int) []byte {
	return MerkleRoot(limit(stream, 1<<k))
}

func readLeaf(stream [][]byte) []byte {
	if len(stream) == 0 {
		return nil
	}
	return leafHash(stream[0])
}

func stringToByteSlice(s string) []byte {
	return []byte(s)
}

func bitTest(i int, n int) bool {
	return (i & (1 << n)) != 0
}

func ones(i int) []int {
	r := []int{}
	for n := 0; n < 32; n++ {
		if bitTest(i, n) {
			r = append(r, n)
		}
	}
	return r
}

func zeros(i int) []int {
	r := []int{}
	for n := 0; n < 32; n++ {
		if !bitTest(i, n) {
			r = append(r, n)
		}
	}
	return r
}

func proveLeaf(stream [][]byte, index int) *proof {
	pre := make([]subRoot, 0, len(ones(index)))
	for _, k := range ones(index) {
		pre = append(pre, subRoot{i: k, subr: subroot(stream, k)})
	}
	post := make([]subRoot, 0, len(zeros(index)))
	for _, k := range zeros(index) {
		post = append(post, subRoot{i: k, subr: subroot(stream, k)})
	}
	return &proof{pre: pre, leaf: readLeaf(stream), post: post}
}

func loadStack(stk map[int][]byte, blks []subRoot) map[int][]byte {
	for _, blk := range blks {
		stk = insert(stk, blk.subr, blk.i)
	}
	return stk
}

func rootFromProofAndLeaf(leaf []byte, proof *proof) []byte {
	stk := make(map[int][]byte)
	stk = loadStack(stk, proof.pre)
	stk = insert(stk, leaf, 0)
	stk = loadStack(stk, proof.post)
	return finalize(stk)
}

func verifyLeaf(knownroot []byte, leaf []byte, proof *proof) bool {
	return bytes.Equal(knownroot, rootFromProofAndLeaf(leaf, proof))
}

func main() {
	// arbitrary example block stream with 12 'blocks'

	blkstream := make([][]byte, 12)
	for i := 0; i < len(blkstream); i++ {
		blkstream[i] = []byte(strconv.Itoa(i))
	}

	// blkstream := [][]byte{[]byte("a"), []byte("b"), []byte("c")}

	// generate proof for leaf at index 6
	proof := proveLeaf(blkstream, 6)
	fmt.Printf("proof: %+v\n", proof)

	// calculate Merkle root
	expectedRoot := "c2ff85521db556cc6b72381f33dbfed5e570f9137acd7e195a496201cf23500e"

	// verify the proof
	verificationHash := rootFromProofAndLeaf(proof.leaf, proof)
	fmt.Printf("verification hash: %s\n", hex.EncodeToString(verificationHash))
	fmt.Printf("expected root: %s\n", expectedRoot)

	// compare the verification hash to the expected Merkle root
	if hex.EncodeToString(verificationHash) == expectedRoot {
		fmt.Println("Verification successful!")
	} else {
		fmt.Println("Verification failed.")
	}
}
