package object

import (
	"fmt"
	"slices"
)

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	DICE_OBJ    = "DICE"
	ERROR_OBJ   = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type DiceData struct {
	Literal    string
	Tags       []string
	RawRolls   []uint
	FinalRolls []uint
	Value      int64
}

func (dd *DiceData) IsEqualTo(other DiceData) bool {
	// for testing purposes to compare equality
	isLit := dd.Literal == other.Literal
	isTags := slices.Compare(dd.Tags, other.Tags) == 0
	isRawRolls := slices.Compare(dd.RawRolls, other.RawRolls) == 0
	isFinalRolls := slices.Compare(dd.FinalRolls, other.FinalRolls) == 0
	isValue := dd.Value == other.Value

	return isLit && isTags && isRawRolls && isFinalRolls && isValue
}

type Metadata struct {
	Store map[string]DiceData
}

func NewMetadata() *Metadata {
	s := make(map[string]DiceData)
	return &Metadata{Store: s}
}

func (m *Metadata) Add(literal string, val DiceData) {
	// NOTE: handles collisions very sloppily. There shouldn't be too many though.
	collisionCounter := 0
	quantifier := "(%d)"
	literal += fmt.Sprintf(quantifier, collisionCounter)
	for {
		if _, ok := m.Store[literal]; ok {
			collisionCounter += 1
			literal = literal[:len(literal)-3] + fmt.Sprintf(quantifier, collisionCounter)
		} else {
			break
		}
	}
	m.Store[literal] = val
}
