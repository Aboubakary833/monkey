package object

import "fmt"


type ObjectType int


const (
	NULL_OBJ		ObjectType = iota
	INTEGER_OBJ
	FLOAT_OBJ
	BOOLEAN_OBJ
)


type Object interface {
	Type()		ObjectType
	Inspect()	string
}


type Integer struct {
	Value		int64
}
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }



type Float struct {
	Value		float64
}
func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string { return fmt.Sprintf("%g", f.Value ) }



type Boolean struct {
	Value		bool
}
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }



type Null struct {}
func (b *Null) Type() ObjectType { return NULL_OBJ }
func (b *Null) Inspect() string { return "null" }




