# fast-merkle

Golang Implementation of Streaming Merkle Root, Proof, and Verify (single leaf) from Luke Champine's paper: **Streaming Merkle Proofs within Binary Numeral Trees** @ https://eprint.iacr.org/2021/038.pdf

## Usage

```golang
	blkstream := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	root := merkleRoot(blkstream)
	fmt.Printf("Merkle root: %x\n", root)
```