package ast

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var amount0 = 0
var amount1 = 1
var amount3 = 3

func Test_Size_Bytes(t *testing.T) {
	tests := []struct {
		size          Size
		expectedBytes uint64
	}{
		{size: Size{Amount: nil, Unit: nil}, expectedBytes: 0},
		{size: Size{Amount: nil, Unit: &SizeKb}, expectedBytes: 0},
		{size: Size{Amount: &amount3, Unit: nil}, expectedBytes: 3},
		{size: Size{Amount: &amount0, Unit: &SizeKb}, expectedBytes: 0},
		{size: Size{Amount: &amount0, Unit: &SizeMb}, expectedBytes: 0},
		{size: Size{Amount: &amount0, Unit: &SizeGb}, expectedBytes: 0},
		{size: Size{Amount: &amount0, Unit: &SizeTb}, expectedBytes: 0},
		{size: Size{Amount: &amount1, Unit: &SizeKb}, expectedBytes: 1024},
		{size: Size{Amount: &amount1, Unit: &SizeMb}, expectedBytes: 1024 * 1024},
		{size: Size{Amount: &amount1, Unit: &SizeGb}, expectedBytes: 1024 * 1024 * 1024},
		{size: Size{Amount: &amount1, Unit: &SizeTb}, expectedBytes: 1024 * 1024 * 1024 * 1024},
		{size: Size{Amount: &amount3, Unit: &SizeKb}, expectedBytes: 3 * 1024},
		{size: Size{Amount: &amount3, Unit: &SizeMb}, expectedBytes: 3 * 1024 * 1024},
		{size: Size{Amount: &amount3, Unit: &SizeGb}, expectedBytes: 3 * 1024 * 1024 * 1024},
		{size: Size{Amount: &amount3, Unit: &SizeTb}, expectedBytes: 3 * 1024 * 1024 * 1024 * 1024},
	}
	for _, test := range tests {
		amount := "nil"
		unit := "nil"
		if test.size.Amount != nil {
			amount = fmt.Sprintf("%d", *test.size.Amount)
		}
		if test.size.Unit != nil {
			unit = fmt.Sprintf("%s", *test.size.Unit)
		}
		name := fmt.Sprintf("Bytes should return %d when size is %s %s", test.expectedBytes, amount, unit)
		t.Run(name, func(t *testing.T) {
			// GIVEN
			// WHEN
			result := test.size.Bytes()
			// THEN
			assert.Equal(t, test.expectedBytes, result)
		})
	}
}
