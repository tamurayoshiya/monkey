package code

import (
	// "github.com/k0kubun/pp"
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{
			OpConstant,
			[]int{65534},
			[]byte{ // 3バイト使う
				byte(OpConstant), // // 1バイト目はopcode(OpConstant)
				255,              // 残りの2バイトはbig-endian encodingで65534を表現
				254,
			},
		},
		{
			OpConstant,
			[]int{254},
			[]byte{
				byte(OpConstant), 0, 254,
			},
		},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want=%d, got=%d",
				len(tt.expected), len(instruction))
		}
		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, b, instruction[i])
			}
		}
	}
}
