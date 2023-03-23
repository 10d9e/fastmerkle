package main

import (
	"bufio"
	"bytes"
	"io"
	"math/big"

	"github.com/minio/sha256-simd"
)

type SubRoot struct {
	I    int
	SubR []byte
}

func Limit(r io.Reader, n int) [][]byte {
	var lines [][]byte
	scanner := bufio.NewScanner(r)
	for scanner.Scan() && (n == 0 || len(lines) < n) {
		lines = append(lines, scanner.Bytes())
	}
	return lines
}

func LeftmostOneBitPos(i *big.Int) int {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	pos := 0
	for i.Cmp(zero) > 0 {
		if i.And(i, i.Sub(i, one)).Cmp(zero) == 0 {
			return pos
		}
		i.Rsh(i, 1)
		pos++
	}
	return -1
}

func RightmostZeroBitPos(i *big.Int, size int) int {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	pos := size - 1
	for i.Cmp(one.Lsh(big.NewInt(int64(size)), uint(pos))) >= 0 && pos >= 0 {
		if i.And(i, one.Lsh(big.NewInt(int64(size)), uint(pos))).Cmp(zero) == 0 {
			return pos
		}
		pos--
	}
	return pos
}

func ProveLeaf(stream io.Reader, index int) map[string]interface{} {
	var pre, post []SubRoot
	for i := LeftmostOneBitPos(big.NewInt(int64(index))); i >= 0; i-- {
		subr := merkleRoot(Limit(stream, 1<<uint(i)))
		pre = append(pre, SubRoot{i, subr})
	}
	leaf := Limit(stream, 1)[0]
	size := len(leaf)
	for i := RightmostZeroBitPos(new(big.Int).SetBytes(leaf), size); i >= 0; i-- {
		subr := merkleRoot(Limit(stream, 1<<uint(i)))
		post = append(post, SubRoot{i, subr})
	}
	return map[string]interface{}{
		"pre":  pre,
		"leaf": leaf,
		"post": post,
	}
}

func LoadStack(stk map[int][]byte, blks []SubRoot) map[int][]byte {
	for _, blk := range blks {
		stk[blk.I] = blk.SubR
	}
	return stk
}

func RootFromProofAndLeaf(leaf []byte, proof map[string]interface{}) []byte {
	stk := make(map[int][]byte)
	for _, pre := range proof["pre"].([]SubRoot) {
		stk = LoadStack(stk, []SubRoot{pre})
	}
	stk = insert(stk, sha256Hash(leaf), 0)
	for _, post := range proof["post"].([]SubRoot) {
		stk = LoadStack(stk, []SubRoot{post})
	}
	return finalize(stk)
}

func VerifyLeaf(knownroot []byte, leaf []byte, proof map[string]interface{}) bool {
	return bytes.Equal(knownroot, RootFromProofAndLeaf(leaf, proof))
}

func sha256Hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}
