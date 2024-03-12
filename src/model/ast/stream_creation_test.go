package ast

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

func Test_Duration_Duration(t *testing.T) {
	tests := []struct {
		duration         Duration
		expectedDuration time.Duration
	}{
		{duration: Duration{Amount: &amount1, Unit: &durationUnitSecond}, expectedDuration: time.Second},
		{duration: Duration{Amount: &amount3, Unit: &durationUnitSeconds}, expectedDuration: 3 * time.Second},
		{duration: Duration{Amount: &amount3, Unit: &durationUnitMinute}, expectedDuration: 3 * time.Minute},
		{duration: Duration{Amount: &amount1, Unit: &durationUnitMinutes}, expectedDuration: time.Minute},
		{duration: Duration{Amount: &amount1, Unit: &durationUnitHours}, expectedDuration: time.Hour},
	}
	for _, test := range tests {
		amount := "nil"
		unit := "nil"
		if test.duration.Amount != nil {
			amount = fmt.Sprintf("%d", *test.duration.Amount)
		}
		if test.duration.Unit != nil {
			unit = fmt.Sprintf("%s", *test.duration.Unit)
		}
		name := fmt.Sprintf("Duration should return %d when Duration is %s %s", test.expectedDuration, amount, unit)
		t.Run(name, func(t *testing.T) {
			// GIVEN
			// WHEN
			result := test.duration.Duration()
			// THEN
			assert.Equal(t, test.expectedDuration, result)
		})
	}
}
