package evaluator

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/daneofmanythings/diceroni/pkg/interpreter/ast"
	"github.com/daneofmanythings/diceroni/pkg/interpreter/object"
)

var NULL = &object.Null{}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	// Expressions
	case *ast.DiceLiteral:
		return evalDiceExpression(node, env) // evaluates the roll and records all metadata in the env

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.Error:
			return result
		}
	}

	return result
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)
	}
	return result
}

// TODO: here is the dice evaluation!
func evalDiceExpression(node ast.Expression, env *object.Environment) object.Object {
	dice, ok := node.(*ast.DiceLiteral)
	if !ok {
		return newError("expected DiceLiteral, got=%v", node.TokenLiteral())
	}

	rawRolls := []uint{}

	if dice.Quantity > 0 {
		for i := 0; i < int(dice.Quantity); i++ {
			rawRolls = rollSingleDie(dice.Size, rawRolls)
		}
	} else {
		rawRolls = rollSingleDie(dice.Size, rawRolls)
	}

	adjustedRolls := slices.Clone(rawRolls)

	if dice.MaxValue > 0 {
		adjustedRolls = applyMaxValue(adjustedRolls, dice.MaxValue)
	}
	if dice.MinValue > 0 {
		adjustedRolls = applyMinValue(adjustedRolls, dice.MinValue)
	}
	if dice.KeepHighest > 0 {
		adjustedRolls = applyKeepHighest(adjustedRolls, dice.KeepHighest)
	}
	if dice.KeepLowest > 0 {
		adjustedRolls = applyKeepLowest(adjustedRolls, dice.KeepLowest)
	}

	return &object.Integer{Value: sumRolls(adjustedRolls)}
}

func rollSingleDie(size uint, rawRolls []uint) []uint {
	roll := rand.Intn(int(size))
	rawRolls = append(rawRolls, uint(roll+1))
	return rawRolls
}

func applyMaxValue(rolls []uint, val uint) []uint {
	for i := 0; i < len(rolls); i++ {
		if rolls[i] > val {
			rolls[i] = val
		}
	}
	return rolls
}

func applyMinValue(rolls []uint, val uint) []uint {
	for i := 0; i < len(rolls); i++ {
		if rolls[i] < val {
			rolls[i] = val
		}
	}
	return rolls
}

func applyKeepHighest(rolls []uint, val uint) []uint {
	return applyKeepFunc(rolls, val, slices.Max)
}

func applyKeepLowest(rolls []uint, val uint) []uint {
	return applyKeepFunc(rolls, val, slices.Min)
}

func applyKeepFunc(rolls []uint, val uint, f func([]uint) uint) []uint {
	resultRolls := []uint{}          // what will be returned
	rollsCopy := slices.Clone(rolls) // what will be updated to track remaining rolls after grabbing a min
	// rolls will be used to sort the returning slice in the order the rolls happened

	if int(val) >= len(rolls) {
		return rolls // more rolls that the keep value
	}

	// grabbing mins, putting them into lowestRolls, and removing them from the copy.
	for i := 0; i < int(val); i++ {
		nextRoll := f(rollsCopy)
		resultRolls = append(resultRolls, nextRoll)
		idx := slices.Index(rollsCopy, nextRoll)
		rollsCopy = slices.Delete(rollsCopy, idx, idx+1)
	}

	// sorting the result in roll order
	slices.SortFunc(resultRolls, func(a, b uint) int {
		if slices.Index(rolls, a) < slices.Index(rolls, b) {
			return -1
		} else {
			return 1
		}
	})

	return resultRolls
}

func sumRolls(rolls []uint) int64 {
	var result int64 = 0
	for _, roll := range rolls {
		result += int64(roll)
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	} else if right.Type() == object.DICE_OBJ {
		return newError("dice not implmented yet: %s", right.Type())
	}
	return newError("unknown operator: -%s", right.Type())
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
