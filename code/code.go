package code

import (
	// "github.com/k0kubun/pp"
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + read
	}
	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}
	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// ------------------------

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpPop
	OpSub
	OpMul
	OpDiv
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJumpNotTruthy
	OpJump
	OpNull
	OpGetGlobal
	OpSetGlobal
	OpArray
	OpHash
	OpIndex
	OpCall
	OpReturnValue // for returning value
	OpReturn      // for returning implicit vm.Null
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {
		Name:          "OpConstant",
		OperandWidths: []int{2}, // = two-byte operand
	},
	OpAdd: {
		Name:          "OpAdd",
		OperandWidths: []int{}, // = has no operands
	},
	OpPop: {
		Name:          "OpPop",
		OperandWidths: []int{},
	},
	OpSub: {
		Name:          "OpSub",
		OperandWidths: []int{},
	},
	OpMul: {
		Name:          "OpMul",
		OperandWidths: []int{},
	},
	OpDiv: {
		Name:          "OpDiv",
		OperandWidths: []int{},
	},
	OpTrue: {
		Name:          "OpTrue",
		OperandWidths: []int{},
	},
	OpFalse: {
		Name:          "OpFalse",
		OperandWidths: []int{},
	},
	OpEqual: {
		Name:          "OpEqual",
		OperandWidths: []int{},
	},
	OpNotEqual: {
		Name:          "OpNotEqual",
		OperandWidths: []int{},
	},
	OpGreaterThan: {
		Name:          "OpGreaterThan",
		OperandWidths: []int{},
	},
	OpMinus: {
		Name:          "OpMinus",
		OperandWidths: []int{},
	},
	OpBang: {
		Name:          "OpBang",
		OperandWidths: []int{},
	},
	OpJumpNotTruthy: {
		Name:          "OpJumpNotTruthy",
		OperandWidths: []int{2},
	},
	OpJump: {
		Name:          "OpJump",
		OperandWidths: []int{2},
	},
	OpNull: {
		Name:          "OpNull",
		OperandWidths: []int{},
	},
	OpGetGlobal: {
		Name:          "OpGetGlobal",
		OperandWidths: []int{2},
	},
	OpSetGlobal: {
		Name:          "OpSetGlobal",
		OperandWidths: []int{2},
	},
	OpArray: {
		Name:          "OpArray",
		OperandWidths: []int{2},
	},
	OpHash: {
		Name:          "OpHash",
		OperandWidths: []int{2},
	},
	OpIndex: {
		Name:          "OpIndex",
		OperandWidths: []int{},
	},
	OpCall: {
		Name:          "OpCall",
		OperandWidths: []int{},
	},
	OpReturnValue: {
		Name:          "OpReturnValue",
		OperandWidths: []int{},
	},
	OpReturn: {
		Name:          "OpReturn",
		OperandWidths: []int{},
	},
}

// Lookup takes a byte of Opcode,
// and returns found 'Definition' from 'definitions' table
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1

	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}

// ReadOperands is the function that reverses everything 'Make' does
// the argument 'ins' expects to be given operand part of bytes from a instructions
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
