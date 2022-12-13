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

// fmt.Printf("========================\n")
// fmt.Printf("|| %d  ||  %d  ||  %d  ||\n", var1, var2, c)
// fmt.Printf("========================")

package evm

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"unicode/utf8"
)

type FunctionMap map[byte]func(code []byte, stack []*big.Int) []*big.Int
type BytesMap map[byte]int
type uint256 big.Int

var funcs FunctionMap
var bytes BytesMap

// func Push(code []byte, stack []*big.Int) []*big.Int {
// 	fmt.Printf("In Push | %d ", code)
// 	n := new(big.Int)
// 	n.SetBytes([]byte{code})
// 	return append(stack, n)
// }

func popFromStack(code []byte, stack []*big.Int) (*big.Int, []*big.Int) {
	return stack[0], stack[1:]
}

func pushToStack(val *big.Int, stack []*big.Int) []*big.Int {
	return append([]*big.Int{val}, stack...)
}

func overflow(val *big.Int) *big.Int {
	return new(big.Int).Mod(val, floatToBigInt(math.Exp2(256)))
}

func Push(code []byte, stack []*big.Int) []*big.Int {
	// Make a new stack with len of original stack
	n := make([]*big.Int, len(stack))
	// Copy old stack to new stack
	copy(n, stack)
	// Create a new big int value
	val := new(big.Int)
	// Set its bytes
	val.SetBytes(code)

	// Set new big int value to a new stack, merge with old stack
	n = append([]*big.Int{val}, stack...)
	// return new stack
	return n
}

func Pop(code []byte, stack []*big.Int) []*big.Int {
	return stack[1:]
}

func Stop(code []byte, stack []*big.Int) []*big.Int {
	return stack
}

func Add(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	var1 = var1.Add(var1, var2)
	var1 = overflow(var1)
	return pushToStack(var1, stack)
}

func Lt(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)
	c := var1.Cmp(var2)
	if c == 1 || c == 0 {
		return pushToStack(big.NewInt(0), stack)
	} else {
		return pushToStack(big.NewInt(1), stack)
	}
}

func SLt(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1     *big.Int
		var2     *big.Int
		var1Sign bool // Negative true or false
		var2Sign bool
	)
	var1Sign = false
	var2Sign = false
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	if isNegative(var1) {
		var1Sign = true
		var1 = twosComplement(var1)
	}

	if isNegative(var2) {
		var2Sign = true
		var2 = twosComplement(var2)
	}

	if var1Sign && !var2Sign {
		return pushToStack(big.NewInt(1), stack)
	}

	if var2Sign && !var1Sign {
		return pushToStack(big.NewInt(0), stack)
	}

	c := var1.Cmp(var2)

	// If they're both negative, swap the responses
	// if var1Sign && var2Sign {
	// 	if c == 1 || c == 0 {
	// 		return pushToStack(big.NewInt(1), stack)
	// 	} else {
	// 		return pushToStack(big.NewInt(0), stack)
	// 	}
	// }

	if c == 1 || c == 0 {
		return pushToStack(big.NewInt(0), stack)
	} else {
		return pushToStack(big.NewInt(1), stack)
	}
}

func Gt(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)
	c := var1.Cmp(var2)
	if c == -1 || c == 0 {
		return pushToStack(big.NewInt(0), stack)
	} else {
		return pushToStack(big.NewInt(1), stack)
	}
}

func SGt(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1     *big.Int
		var2     *big.Int
		var1Sign bool // Negative true or false
		var2Sign bool
	)
	var1Sign = false
	var2Sign = false
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	if isNegative(var1) {
		var1Sign = true
		var1 = twosComplement(var1)
	}

	if isNegative(var2) {
		var2Sign = true
		var2 = twosComplement(var2)
	}

	if var1Sign && !var2Sign {
		return pushToStack(big.NewInt(1), stack)
	}

	if var2Sign && !var1Sign {
		return pushToStack(big.NewInt(0), stack)
	}

	c := var1.Cmp(var2)

	// If they're both negative, swap the responses
	if var1Sign && var2Sign {
		if c == -1 {
			return pushToStack(big.NewInt(1), stack)
		} else {
			return pushToStack(big.NewInt(0), stack)
		}
	}

	if c == -1 || c == 0 {
		return pushToStack(big.NewInt(0), stack)
	} else {
		return pushToStack(big.NewInt(1), stack)
	}
}

func Mul(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	var1 = var1.Mul(var1, var2)
	var1 = overflow(var1)
	return pushToStack(var1, stack)
}

func Sub(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	var1 = var1.Sub(var1, var2)
	var1 = overflow(var1)
	return pushToStack(var1, stack)
}

func Div(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	if var2.BitLen() == 0 {
		return pushToStack(var2, stack)
	} else {
		var1 = var1.Div(var1, var2)
		var1 = overflow(var1)
		return pushToStack(var1, stack)
	}
}

func twosComplement(n *big.Int) *big.Int {
	// Check if n is zero
	if n.Sign() == 0 {
		return n
	}

	// Calculate the maximum value for the given bit size of n
	max := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(n.BitLen())), nil)

	// Subtract n from the maximum value
	return new(big.Int).Sub(max, n)
}

func SDiv(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1     *big.Int
		var2     *big.Int
		var1Sign bool
		var2Sign bool
	)
	var1Sign = false
	var2Sign = false
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	if isNegative(var1) {
		var1Sign = true
		var1 = twosComplement(var1)
	}

	if isNegative(var2) {
		var2Sign = true
		var2 = twosComplement(var2)
	}

	if var2.BitLen() == 0 {
		return pushToStack(var2, stack)
	} else {
		var1 = var1.Div(var1, var2)
		var1 = overflow(var1)
		if var1Sign != var2Sign {
			var1 = flipToNegative(var1)
		}
		return pushToStack(var1, stack)
	}
}

func SMod(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1     *big.Int
		var2     *big.Int
		var1Sign bool
		var2Sign bool
	)
	var1Sign = false
	var2Sign = false
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	if isNegative(var1) {
		var1Sign = true
		var1 = twosComplement(var1)
	}

	if isNegative(var2) {
		var2Sign = true
		var2 = twosComplement(var2)
	}

	var res *big.Int
	if var2.BitLen() == 0 {
		res = new(big.Int).Mod(var2, var1)
		return pushToStack(res, stack)
	} else {
		res = new(big.Int).Mod(var1, var2)
	}

	if var1Sign || var2Sign {
		res = flipToNegative(res)
	}
	return pushToStack(res, stack)
}

func flipToNegative(v *big.Int) *big.Int {
	y := floatToBigInt(math.Exp2(256))
	y = y.Sub(y, v)
	return y
}

func Mod(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	var res *big.Int
	if var2.BitLen() == 0 {
		res = new(big.Int).Mod(var2, var1)
	} else {
		res = new(big.Int).Mod(var1, var2)
	}
	return pushToStack(res, stack)
}

func AddMod(code []byte, stack []*big.Int) []*big.Int {
	stack = Add(code, stack)
	stack = Mod(code, stack)
	return stack
}

func MulMod(code []byte, stack []*big.Int) []*big.Int {
	// stack = Mul(code, stack)
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	var1 = var1.Mul(var1, var2)
	stack = pushToStack(var1, stack)
	stack = Mod(code, stack)

	var1, stack = popFromStack(code, stack)
	var1 = overflow(var1)
	return pushToStack(var1, stack)
}

func Exp(code []byte, stack []*big.Int) []*big.Int {
	// stack = Mul(code, stack)
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	var1 = new(big.Int).Exp(var1, var2, big.NewInt(0))
	var1 = overflow(var1)
	return pushToStack(var1, stack)
}

func checkFirstBit(s string) int {
	r, _ := utf8.DecodeRuneInString(s)
	if string(r) == "0" {
		return 0
	} else {
		return 1
	}
	// return string(r) == "0"? 0 : 1
}

func paddedBinary(x *big.Int, n int) string {
	// res := fmt.Sprintf("%0*b", n, x)
	s := strings.Repeat("0", n)
	return s
}

func padNegative(x *big.Int, n int) string {
	// res := fmt.Sprintf("%0*b", n, x)
	s := strings.Repeat("1", n)
	return s
}

func FullBinary(x *big.Int) string {
	// Convert the big.Int to a binary string
	s := x.Text(2)

	// Pad the binary string with zeros to make it 32 bits
	s = strings.Repeat("0", 32-len(s)) + s

	return s
}

func isNegative(value *big.Int) bool {
	hold := make([]big.Word, value.BitLen())
	copy(hold, value.Bits())
	binaryLength := len(value.Text(2))
	additionalLength := 0
	if binaryLength%8 != 0 {
		additionalLength = 8 - binaryLength%8
	}
	binaryString := paddedBinary(value, additionalLength)
	if checkFirstBit(binaryString) == 0 {
		return false
	} else {
		return true
	}
}

func SignNumber(value *big.Int) *big.Int {
	hold := make([]big.Word, value.BitLen())
	copy(hold, value.Bits())
	res := padNegative(big.NewInt(0).SetBits(hold), 256)
	x, _ := new(big.Int).SetString(res, 2)
	return x
}

func SignExtend(code []byte, stack []*big.Int) []*big.Int {
	var (
		value *big.Int
	)
	_, stack = popFromStack(code, stack)
	value, stack = popFromStack(code, stack)
	hold := make([]big.Word, value.BitLen())
	copy(hold, value.Bits())
	binaryLength := len(value.Text(2))
	additionalLength := 0
	if binaryLength%8 != 0 {
		additionalLength = 8 - binaryLength%8
	}
	binaryString := paddedBinary(value, additionalLength)
	if checkFirstBit(binaryString) == 0 {
		stack = pushToStack(big.NewInt(0).SetBits(hold), stack)
		return stack
	} else {
		res := padNegative(big.NewInt(0).SetBits(hold), 256)
		x, _ := new(big.Int).SetString(res, 2)
		return pushToStack(x, stack)
	}
}

func SignExtendSingle(value *big.Int) *big.Int {
	hold := make([]big.Word, value.BitLen())
	copy(hold, value.Bits())
	// binaryLength := len(value.Text(2))
	// additionalLength := 0
	// if binaryLength%8 != 0 {
	// 	additionalLength = 8 - binaryLength%8
	// }
	// binaryString := paddedBinary(value, additionalLength)
	// if checkFirstBit(binaryString) == 0 {
	// 	return big.NewInt(0).SetBits(hold)
	// } else {
	res := padNegative(big.NewInt(0).SetBits(hold), 256)
	x, _ := new(big.Int).SetString(res, 2)
	return x
	// }
}

func Eq(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	c := var1.Cmp(var2)
	if c == 0 {
		return pushToStack(big.NewInt(1), stack)
	} else {
		return pushToStack(big.NewInt(0), stack)
	}
}

func IsZero(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
	)
	var1, stack = popFromStack(code, stack)

	c := var1.Cmp(big.NewInt(0))
	if c == 0 {
		return pushToStack(big.NewInt(1), stack)
	} else {
		return pushToStack(big.NewInt(0), stack)
	}
}

// func Not(code []byte, stack []*big.Int) []*big.Int {
// 	var (
// 		var1 *big.Int
// 	)
// 	var1, stack = popFromStack(code, stack)
// 	// fmt.Printf("========================\n")
// 	// fmt.Printf("|| %d  ||\n", var1)
// 	// fmt.Printf("========================")
// 	// // y := floatToBigInt(math.Exp2(256))
// 	// // var1 = y.Sub(y, var1)
// 	// var1 = var1.Or(var1, floatToBigInt(math.Exp2(256)))
// 	// fmt.Printf("========================\n")
// 	// fmt.Printf("|| %d  ||\n", var1)
// 	// fmt.Printf("========================")
// 	// var1 = SignExtendSingle(var1)
// 	// fmt.Printf("========================\n")
// 	// fmt.Printf("|| %d  ||\n", var1)
// 	// fmt.Printf("========================")
// 	var1 = var1.Not(var1)

//     // Convert the result back to a big.Int.
//     noVar1, ok := new(big.Int).SetString(var1, 16)
// 	return pushToStack(var1, stack)

// }

func Not(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	max := floatToBigInt(math.Exp2(256))
	max = max.Sub(max, big.NewInt(1))
	var1 = max.Sub(max, var1)

	return pushToStack(var1, stack)

}

func And(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)
	// max := floatToBigInt(math.Exp2(256))
	// max = max.Sub(max, big.NewInt(1))
	var1 = var2.And(var2, var1)
	fmt.Printf("========================\n")
	fmt.Printf("|| %s  || %s || \n", var1.Text(16), var2.Text(16))
	fmt.Printf("========================")
	return pushToStack(var1, stack)
}

func Or(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	var1 = var2.Or(var2, var1)
	return pushToStack(var1, stack)
}

func Xor(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	var1 = var2.Xor(var2, var1)
	return pushToStack(var1, stack)
}

func Shl(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
		var2 *big.Int
	)

	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	max := floatToBigInt(math.Exp2(256))
	max = max.Sub(max, big.NewInt(1))

	// If we're pushing by more than 1024 (16 * 64) number of Fs in hex
	if var1.Cmp(big.NewInt(1024)) == 1 {
		return pushToStack(big.NewInt(0), stack)
	}
	var1 = var2.Lsh(var2, uint(var1.Int64()))

	res := new(big.Int).And(var1, max)
	fmt.Printf("========================\n")
	fmt.Printf("|| %s \n|| %s \n|| %s || \n", var1.Text(16), res.Text(16), max.Text(16))
	fmt.Printf("========================")
	return pushToStack(res, stack)
}

func Shr(code []byte, stack []*big.Int) []*big.Int {
	var (
		var1 *big.Int
		var2 *big.Int
	)
	var1, stack = popFromStack(code, stack)
	var2, stack = popFromStack(code, stack)

	var1 = var2.Rsh(var2, uint(var1.Int64()))
	return pushToStack(var1, stack)
}

// Run runs the EVM code and returns the stack and a success indicator.
func Evm(code []byte) ([]*big.Int, bool) {
	var stack []*big.Int
	var remainder []byte
	funcs, bytes = buildMaps()

Opcode:
	pc := 0

	op := code[pc] // First bytes in code
	if op == 0 {
		goto Done
	}

	// PUSH1 - PUSH32
	if 96 <= op && op <= 127 {

		size := int(op) - 96 + 1
		val := code[pc+1 : pc+1+int(size)]
		remainder = code[pc+1+int(size):]
		stack = Push(val, stack)
	} else {
		size := bytes[op]
		val := code[pc+1 : pc+1+int(size)]
		remainder = code[pc+1+int(size):]
		stack = funcs[op](val, stack)
	}

	if len(remainder) > 0 {
		code = remainder
		goto Opcode
	} else {
		return stack, true
	}

Done:
	return stack, true

}

func buildMaps() (FunctionMap, BytesMap) {
	funcs = make(map[byte]func(code []byte, stack []*big.Int) []*big.Int)
	bytes = make(map[byte]int)

	funcs[0] = Stop
	bytes[0] = 0

	funcs[80] = Pop
	bytes[80] = 0

	funcs[1] = Add
	bytes[1] = 0

	funcs[2] = Mul
	bytes[2] = 0

	funcs[3] = Sub
	bytes[3] = 0

	funcs[4] = Div
	bytes[4] = 0

	funcs[5] = SDiv
	bytes[5] = 0

	funcs[6] = Mod
	bytes[6] = 0

	funcs[7] = SMod
	bytes[7] = 0

	funcs[8] = AddMod
	bytes[8] = 0

	funcs[9] = MulMod
	bytes[9] = 0

	funcs[10] = Exp
	bytes[10] = 0

	funcs[11] = SignExtend
	bytes[11] = 0

	funcs[16] = Lt
	bytes[16] = 0

	funcs[17] = Gt
	bytes[17] = 0

	funcs[18] = SLt
	bytes[18] = 0

	funcs[19] = SGt
	bytes[19] = 0

	funcs[20] = Eq
	bytes[20] = 0

	funcs[21] = IsZero
	bytes[21] = 0

	funcs[22] = And
	bytes[22] = 0

	funcs[23] = Or
	bytes[23] = 0

	funcs[24] = Xor
	bytes[24] = 0

	funcs[25] = Not
	bytes[25] = 0

	funcs[27] = Shl
	bytes[27] = 0

	funcs[28] = Shr
	bytes[28] = 0

	return funcs, bytes
}

func floatToBigInt(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)

	coin := new(big.Float)
	coin.SetInt(big.NewInt(1))

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result) // store converted number in result

	return result
}

// padSigned pads a signed value with the specified number of bits to a two's complement representation.
// func padSigned(value []byte, bits int) []byte {
// 	// Calculate the required padding
// 	padding := (bits - len(value)*8) % 8

// 	// If the value is negative and there is padding, add the padding to the value
// 	if value[0]&0x80 != 0 && padding > 0 {
// 		value = append([]byte{0xff}, value...)
// 	}

// 	// Return the padded value
// 	return value
// }

// func SignExtEnd(code []byte, stack []*big.Int) []*big.Int {
// 	// Get the second and third items from the stack.
// 	// These are the position and value to sign extend, respectively.
// 	pos := stack[len(stack)-2]
// 	val := stack[len(stack)-1]

// 	// Check if the position is within the bounds of the value.
// 	if pos.Cmp(big.NewInt(31)) < 0 && pos.Cmp(big.NewInt(0)) >= 0 {
// 		// Create a new big.Int to store the result.
// 		res := big.NewInt(0)

// 		// Convert the value to two's complement representation.
// 		val.Neg(val)
// 		val.Add(val, big.NewInt(1))

// 		// Check if the value at the specified position is equal to 1 or 0.
// 		// If it is 1, the value should be sign extended for positive numbers.
// 		// If it is 0, the value should be sign extended for negative numbers.
// 		if val.Bit(int(pos.Int64())) == 1 || val.Bit(int(pos.Int64())) == 0 {
// 			// Create a mask with all bits set to 1.
// 			mask := big.NewInt(0)
// 			mask.SetBit(mask, 256, 1)

// 			// Convert the result of the subtraction to uint before passing it to the Rsh function.
// 			mask.Rsh(mask, uint(255-pos.Int64()))

// 			// Perform a bitwise OR operation between the value and the mask.
// 			// This will set all bits after the specified position to 1 for positive numbers,
// 			// or all bits before the specified position to 1 for negative numbers.
// 			res.Or(val, mask)

// 			// Convert the value back to its original representation.
// 			res.Neg(res)
// 			res.Add(res, big.NewInt(1))
// 		} else {
// 			// If the value at the specified position is not 1 or 0,
// 			// then no sign extension is needed, so we can just return the value.
// 			res = val
// 		}

// 		// Pop the last two items from the stack and append the result to the stack.
// 		return append(stack[:len(stack)-2], res)
// 	} else {
// 		// If the position is not within the bounds of the value,
// 		// we return the original stack without modifying it.
// 		return stack
// 	}
// }

// SignExtend converts the negative value at the specified byte offset to a two's complement representation.
// func SignExtend(code []byte, stack []*big.Int) []*big.Int {
// 	// Get the byte offset and bit offset from the code
// 	offset := new(big.Int).SetBytes(code[1:])
// 	bitoffset := code[0]

// 	// Get the value at the specified offset on the stack
// 	value := stack[offset.Int64()]

// 	// Convert the value to a two's complement representation with the specified number of bits
// 	result := new(big.Int).SetBytes(padSigned(value.Bytes(), int(bitoffset)))

// 	// Return the result on the stack
// 	return append(stack, result)
// }
