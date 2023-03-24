# fast-merkle

Golang Implementation of Streaming Merkle Root, Proof, and Verify (single leaf) from Luke Champine's paper: **Streaming Merkle Proofs within Binary Numeral Trees** @ https://eprint.iacr.org/2021/038.pdf

## Usage

```golang
import "github.com/jlogelin/fast-merkle"

func main() {
  blkstream := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
  root := MerkleRoot(blkstream)
  fmt.Printf("Merkle root: %x\n", root)
}
```

## Benchmark

The following root node calculation benchmark was run on an `Apple M1 Pro 32 GB`.

Parameters: 
- 33554432 elements
- sha256 hash

| Project  | Execution Time |
| ------------- | ------------- |
| [fast-merkle](https://github.com/jlogelin/fastmerkle) | 8.320988875s |
| [cbergoon/merkletree](https://github.com/cbergoon/merkletree) | 26.294219875s |
| [wealdtech/go-merkletree](https://github.com/wealdtech/go-merkletree) | 30.399641917s |
| [txaty/go-merkletree](https://github.com/txaty/go-merkletree) | 2m 37s |
