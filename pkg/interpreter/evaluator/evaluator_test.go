package evaluator

import (
	"slices"
	"testing"

	"github.com/daneofmanythings/diceroni/pkg/interpreter/lexer"
	"github.com/daneofmanythings/diceroni/pkg/interpreter/object"
	"github.com/daneofmanythings/diceroni/pkg/interpreter/parser"
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
		name        string
		repetitions int
		input       string
		expected    float64
	}{
		{"sanity1", 1, "5 + 5", 10},
		{"sanity2", 1, "5 + 5 * 2", 15},
		{"sanity3", 1, "(5 + 5) * 2", 20},
		{"2d20 + 10", 100000, "d20mq2 + 10", 31},
		{"4d8 - 2 max6", 100000, "d8mq4ma6 - 2", (33*4)/8.0 - 2},
		{"attack at adv", 100000, "d20mq2mh1 + 5", 18.5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultSum := 0
			for i := 0; i < tc.repetitions; i++ {
				l := lexer.New(tc.input)
				p := parser.New(l)
				program := p.ParseProgram()
				evaluation := Eval(program, &object.Environment{})
				resultSum += int(evaluation.(*object.Integer).Value)
			}
			result := resultSum / tc.repetitions
			if result != int(tc.expected) {
				t.Fatalf("expected=%d, got=%d", int(tc.expected), result)
			}
		})
	}
}
