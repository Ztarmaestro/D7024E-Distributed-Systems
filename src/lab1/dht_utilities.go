package dht

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"math/big"
)

func distance(a, b []byte, bits int) *big.Int {
	var ring big.Int
	ring.Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)

	var a_int, b_int big.Int
	(&a_int).SetBytes(a)
	(&b_int).SetBytes(b)

	var dist big.Int
	(&dist).Sub(&b_int, &a_int)

	(&dist).Mod(&dist, &ring)
	return &dist
}

func between(id1, id2, key []byte) bool {
	// 0 if a==b, -1 if a < b, and +1 if a > b

	if bytes.Compare(key, id2) == 0 { // key == id1
		return true
	}

	if bytes.Compare(id2, id1) == 1 { // id2 > id1
		if bytes.Compare(key, id2) == -1 && bytes.Compare(key, id1) == 1 { // key < id2 && key > id1
			return true
		} else {
			return false
		}
	} else { // id1 > id2
		if bytes.Compare(key, id1) == 1 || bytes.Compare(key, id2) == -1 { // key > id1 || key < id2
			return true
		} else {
			return false
		}
	}
}

// (n + 2^(k-1)) mod (2^m)
func calcFinger(n []byte, k int, m int) (string, []byte) {
	// convert the n to a bigint
	nBigInt := big.Int{}
	nBigInt.SetBytes(n)

	// get the right addend, i.e. 2^(k-1)
	two := big.NewInt(2)
	addend := big.Int{}
	addend.Exp(two, big.NewInt(int64(k-1)), nil)

	// calculate sum
	sum := big.Int{}
	sum.Add(&nBigInt, &addend)

	// calculate 2^m
	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(m)), nil)

	// apply the mod
	result := big.Int{}
	result.Mod(&sum, &ceil)

	resultBytes := result.Bytes()
	resultHex := fmt.Sprintf("%x", resultBytes)

	return resultHex, resultBytes
}

func generateNodeId() string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	// calculate sha-1 hash
	hasher := sha1.New()
	hasher.Write([]byte(u.String()))

	return fmt.Sprintf("%x", hasher.Sum(nil))
}
