package hnsw

import "fmt"

type BinaryString []byte

func (bs BinaryString) String() string {
	return fmt.Sprintf("%08b", bs)
}
