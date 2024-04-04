package evaluator

import (
	"slices"
	"testing"

	"github.com/daneofmanythings/calcuroller/pkg/interpreter/lexer"
	"github.com/daneofmanythings/calcuroller/pkg/interpreter/object"
	"github.com/daneofmanythings/calcuroller/pkg/interpreter/parser"
)

func TestRollSingleDie(t *testing.T) {
	testCases := []struct {
		name        string
		val         uint
		repetitions int
	}{
		{"roll d20", 20, 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.repetitions; i++ {
				result := rollSingleDie(tc.val, []uint{})
				if result[0] < 1 || result[0] > tc.val {
					t.Fatalf("got a roll out of range. min=1, max=%d. got=%d", tc.val, result)
				}
			}
		})
	}
}

func TestApplyMaxValue(t *testing.T) {
	testCases := []struct {
		name     string
		rolls    []uint
		val      uint
		expected []uint
	}{
		{"4 rolls, max 5", []uint{7, 1, 6, 5}, 5, []uint{5, 1, 5, 5}},
		{"4 rolls, max 1", []uint{7, 2, 6, 8}, 1, []uint{1, 1, 1, 1}},
		{"5 rolls, max 10", []uint{7, 2, 6, 8}, 10, []uint{7, 2, 6, 8}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := applyMaxValue(tc.rolls, tc.val)
			if slices.Compare(result, tc.expected) != 0 {
				t.Fatalf("expected=%d, got=%d", tc.expected, result)
			}
		})
	}
}

func TestApplyMinValue(t *testing.T) {
	testCases := []struct {
		name     string
		rolls    []uint
		val      uint
		expected []uint
	}{
		{"4 rolls, min 5", []uint{7, 1, 6, 5}, 5, []uint{7, 5, 6, 5}},
		{"4 rolls, min 1", []uint{7, 2, 6, 8}, 1, []uint{7, 2, 6, 8}},
		{"5 rolls, min 10", []uint{7, 2, 6, 8}, 10, []uint{10, 10, 10, 10}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := applyMinValue(tc.rolls, tc.val)
			if slices.Compare(result, tc.expected) != 0 {
				t.Fatalf("expected=%d, got=%d", tc.expected, result)
			}
		})
	}
}

func TestApplyKeepHighest(t *testing.T) {
	testCases := []struct {
		name     string
		rolls    []uint
		val      uint
		expected []uint
	}{
		{"4 rolls, highest 2", []uint{6, 1, 7, 5}, 2, []uint{6, 7}},
		{"2 rolls, highest 1", []uint{2, 20}, 1, []uint{20}},
		{"5 rolls, highest 10", []uint{7, 2, 6, 8, 1}, 10, []uint{7, 2, 6, 8, 1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := applyKeepHighest(tc.rolls, tc.val)
			if slices.Compare(result, tc.expected) != 0 {
				t.Fatalf("expected=%d, got=%d", tc.expected, result)
			}
		})
	}
}

func TestApplyKeepLowest(t *testing.T) {
	testCases := []struct {
		name     string
		rolls    []uint
		val      uint
		expected []uint
	}{
		{"4 rolls, lowest 2", []uint{7, 5, 6, 1}, 2, []uint{5, 1}},
		{"2 rolls, lowest 1", []uint{2, 20}, 1, []uint{2}},
		{"5 rolls, lowest 10", []uint{7, 2, 6, 8, 10}, 10, []uint{7, 2, 6, 8, 10}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := applyKeepLowest(tc.rolls, tc.val)
			if slices.Compare(result, tc.expected) != 0 {
				t.Fatalf("expected=%d, got=%d", tc.expected, result)
			}
		})
	}
}

// NOTE: This method is precarious. random seeding would be better.
func TestEval(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int64
		dicedata map[string]object.DiceData
	}{
		{"sanity1", "5 + 5", 10, map[string]object.DiceData{
			"5(0)": {Literal: "5", Tags: []string{}, RawRolls: []uint{}, FinalRolls: []uint{}, Value: 5},
			"5(1)": {Literal: "5", Tags: []string{}, RawRolls: []uint{}, FinalRolls: []uint{}, Value: 5},
		}},
		{"sanity2", "5 + 5 * 2[test][another one]", 15, map[string]object.DiceData{
			"5(0)": {Literal: "5", Tags: []string{}, RawRolls: []uint{}, FinalRolls: []uint{}, Value: 5},
			"5(1)": {Literal: "5", Tags: []string{}, RawRolls: []uint{}, FinalRolls: []uint{}, Value: 5},
			"2(0)": {Literal: "2", Tags: []string{"test", "another one"}, RawRolls: []uint{}, FinalRolls: []uint{}, Value: 2},
		}},
		{"sanity3", "(5[first] + 5[second]) * 2[third]", 20, map[string]object.DiceData{
			"5(0)": {Literal: "5", Tags: []string{"first"}, RawRolls: []uint{}, FinalRolls: []uint{}, Value: 5},
			"5(1)": {Literal: "5", Tags: []string{"second"}, RawRolls: []uint{}, FinalRolls: []uint{}, Value: 5},
			"2(0)": {Literal: "2", Tags: []string{"third"}, RawRolls: []uint{}, FinalRolls: []uint{}, Value: 2},
		}},
		{"2d1 + 10", "d1qu2[test] + 10", 12, map[string]object.DiceData{
			"2d1[test](0)": {Literal: "2d1[test]", Tags: []string{"test"}, RawRolls: []uint{1, 1}, FinalRolls: []uint{1, 1}, Value: 2},
			"10(0)":        {Literal: "10", Tags: []string{}, RawRolls: []uint{}, FinalRolls: []uint{}, Value: 10},
		}},
		{"4d1kh3 - 2", "d1qu4kh3 - 2", 1, map[string]object.DiceData{
			"4d1kh3(0)": {Literal: "4d1kh3", Tags: []string{}, RawRolls: []uint{1, 1, 1, 1}, FinalRolls: []uint{1, 1, 1}, Value: 3},
			"2(0)":      {Literal: "2", Tags: []string{}, RawRolls: []uint{}, FinalRolls: []uint{}, Value: 2},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := lexer.New(tc.input)
			p := parser.New(l)
			program := p.ParseProgram()
			metadata := object.NewMetadata()
			evaluation := Eval(program, metadata)
			result := evaluation.(*object.Integer).Value
			if result != tc.expected {
				t.Fatalf("expected=%d, got=%d", int(tc.expected), result)
			}
			for key, tcdd := range tc.dicedata {
				dd, ok := metadata.Store[key]
				if !ok {
					t.Fatalf("data not found for key=%s.\nmetadata=%v", key, metadata)
				}
				if !dd.IsEqualTo(tcdd) {
					t.Fatalf("expected=%v, got=%v", tcdd, dd)
				}
			}
		})
	}
}
