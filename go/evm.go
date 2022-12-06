// Package evm is an **incomplete** implementation of the Ethereum Virtual
// Machine for the "EVM From Scratch" course:
// https://github.com/w1nt3r-eth/evm-from-scratch
//
// To work on EVM From Scratch In Go:
//
// - Install Golang: https://golang.org/doc/install
// - Go to the `go` directory: `cd go`
// - Edit `evm.go` (this file!), see TODO below
// - Run `go test ./...` to run the tests
package evm

import (
	"fmt"
	"math/big"
)

type FunctionMap map[byte]func(code byte, stack []*big.Int) []*big.Int
type BytesMap map[byte]int

var funcs FunctionMap
var bytes BytesMap

// func byteToBigInt(b byte) *[]*big.Int {
// 	n := new(big.Int)
// 	n.SetBytes([]byte{b})

// 	numbers := []big.Int{}
// 	numbers = append(numbers, *n)

// 	return numbers
// }

// func addToStack(code byte, stack []*big.Int) []*big.Int {
// 	stackValue := stack[0]
// }

func Push(code byte, stack []*big.Int) []*big.Int {
	fmt.Printf("In Push | %d ", code)
	n := new(big.Int)
	n.SetBytes([]byte{code})
	return append(stack, n)
}

func Stop(code byte, stack []*big.Int) []*big.Int {
	return stack
}

// Run runs the EVM code and returns the stack and a success indicator.
func Evm(code []byte) ([]*big.Int, bool) {
	var stack []*big.Int
	pc := 0
	funcs, bytes = buildMaps()
	fmt.Printf("Funcs: %T", funcs[96])

	op := code[pc] // First bytes in code
	bytesRequired := bytes[op]

	fmt.Printf("\n==================\n")
	fmt.Printf("%d || %d || %d ", op, code[pc:], bytesRequired)
	fmt.Printf("\n==================\n")

	pc++
	for pc <= bytesRequired {
		val := code[pc]
		stack = funcs[op](val, stack)
		pc++
	}
	return stack, true
}

func buildMaps() (FunctionMap, BytesMap) {
	funcs = make(map[byte]func(code byte, stack []*big.Int) []*big.Int)
	bytes = make(map[byte]int)

	funcs[0] = Stop
	bytes[0] = 0
	funcs[96] = Push
	bytes[96] = 1
	funcs[97] = Push
	bytes[97] = 2
	return funcs, bytes
}
