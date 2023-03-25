package fastmerkle

import "bytes"

type Proof struct {
	pre  []subRoot
	leaf []byte
	post []subRoot
}

type subRoot struct {
	i    int
	subr []byte
}

func (tree *MerkleTree) limit(stream [][]byte, n int) [][]byte {
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

func (tree *MerkleTree) subroot(stream [][]byte, k int) []byte {
	return tree.root(tree.limit(stream, 1<<k))
}

func (tree *MerkleTree) readLeaf(stream [][]byte) []byte {
	if len(stream) == 0 {
		return nil
	}
	return tree.leafHash(stream[0])
}

func (tree *MerkleTree) stringToByteSlice(s string) []byte {
	return []byte(s)
}

func (tree *MerkleTree) bitTest(i int, n int) bool {
	return (i & (1 << n)) != 0
}

func (tree *MerkleTree) ones(i int) []int {
	r := []int{}
	for n := 0; n < 32; n++ {
		if tree.bitTest(i, n) {
			r = append(r, n)
		}
	}
	return r
}

func (tree *MerkleTree) zeros(i int) []int {
	r := []int{}
	for n := 0; n < 32; n++ {
		if !tree.bitTest(i, n) {
			r = append(r, n)
		}
	}
	return r
}

func (tree *MerkleTree) proveLeaf(stream [][]byte, index int) *Proof {
	pre := make([]subRoot, 0, len(tree.ones(index)))
	for _, k := range tree.ones(index) {
		pre = append(pre, subRoot{i: k, subr: tree.subroot(stream, k)})
	}
	post := make([]subRoot, 0, len(tree.zeros(index)))
	for _, k := range tree.zeros(index) {
		post = append(post, subRoot{i: k, subr: tree.subroot(stream, k)})
	}
	return &Proof{pre: pre, leaf: tree.readLeaf(stream), post: post}
}

func (tree *MerkleTree) loadStack(stk map[int][]byte, blks []subRoot) map[int][]byte {
	for _, blk := range blks {
		stk = tree.insert(stk, blk.subr, blk.i)
	}
	return stk
}

func (tree *MerkleTree) rootFromProofAndLeaf(leaf []byte, proof *Proof) []byte {
	stk := make(map[int][]byte)
	stk = tree.loadStack(stk, proof.pre)
	stk = tree.insert(stk, leaf, 0)
	stk = tree.loadStack(stk, proof.post)
	stk = tree.insert(stk, leaf, 0)
	return tree.finalize()
}

func (tree *MerkleTree) verifyLeaf(knownroot []byte, leaf []byte, proof *Proof) bool {
	return bytes.Equal(knownroot, tree.rootFromProofAndLeaf(leaf, proof))
}
