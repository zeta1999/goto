package object

import (
	"fmt"
	"strings"

	"github.com/pandeykartikey/goto/ast"
)

type Type string
type BuiltinFunction func(args ...Object) Object

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	LOOP_CONTROL_OBJ = "LOOP_CONTROL"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	LIST_OBJ         = "LIST"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%v", b.Value)
}

type Null struct{}

func (n *Null) Type() Type {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

type String struct {
	Value string
}

func (s *String) Type() Type {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() Type {
	return RETURN_VALUE_OBJ
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

type LoopControl struct {
	Value string
}

func (l *LoopControl) Type() Type {
	return LOOP_CONTROL_OBJ
}

func (l *LoopControl) Inspect() string {
	return l.Value
}

type Function struct {
	ParameterList *ast.IdentifierList
	FuncBody      *ast.BlockStatement
}

func (f *Function) Type() Type {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out strings.Builder

	out.WriteString(f.ParameterList.String())
	out.WriteString(" ")
	out.WriteString(f.FuncBody.String())

	return out.String()
}

type List struct {
	Value []Object
}

func (l *List) Type() Type {
	return LIST_OBJ
}

func (l *List) Inspect() string {
	var out strings.Builder

	out.WriteString("[")
	for idx, param := range l.Value {
		if idx > 0 {
			out.WriteString(", ")
		}
		out.WriteString(param.Inspect())
	}

	out.WriteString("]")

	return out.String()
}

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type {
	return BUILTIN_OBJ
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

type Error struct {
	Message string
}

func (e *Error) Type() Type {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return "Error: " + e.Message
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (env *Environment) Get(id string) (Object, bool) {
	value, ok := env.store[id]
	if !ok && env.outer != nil {
		value, ok = env.outer.Get(id)
	}
	return value, ok
}

func (env *Environment) Create(id string, obj Object) (Object, bool) {
	_, ok := env.store[id]
	if ok {
		return nil, false
	}
	env.store[id] = obj
	return env.store[id], true
}

func (env *Environment) Update(id string, obj Object) (Object, bool) {
	_, ok := env.Get(id)
	if !ok {
		return nil, false
	}
	if _, ok = env.store[id]; ok {
		env.store[id] = obj
		return env.store[id], true
	}

	return env.outer.Update(id, obj)
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store, outer: nil}
}

func ExtendEnv(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
