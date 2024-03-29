package codegen

import (
	"log"
	"strings"
)

type Item interface {
	Code(w *Writer)
}

//--------------------------------------------------------------------------------------------------------------//
// global

type Global struct {
	DoNotEdit string
	Package   string
	Imports   Imports
	Items     []Item
}

func (t *Global) AddImport(item *Import) {
	t.Imports.Add(item)
}

func (t *Global) AddItem(i Item) {
	if t.Items == nil {
		t.Items = make([]Item, 0, 10)
	}
	t.Items = append(t.Items, i)
}

func (t *Global) Code(w *Writer) {
	w.W("%s\n", t.DoNotEdit)         // comment ( do not edit )
	w.W("package %s\n\n", t.Package) // package
	t.Imports.Code(w)                // import

	for _, pt := range t.Items {
		pt.Code(w)
	}
}

//--------------------------------------------------------------------------------------------------------------//
// struct

type Struct struct {
	Name    string
	Fields  *Vars
	Methods []*Function
}

func (t *Struct) AddField(item *Var) {
	if t.Fields == nil {
		t.Fields = &Vars{}
		t.Fields.Init(VarScopeStructField)
	}
	t.Fields.Add(item)
}

func (t *Struct) AddFunction(item *Function) {
	if t.Methods == nil {
		t.Methods = make([]*Function, 0, 10)
	}
	t.Methods = append(t.Methods, item)
}

func (t *Struct) Code(w *Writer) {
	w.W("type %s struct {\n", t.Name)
	w.IndentIn()
	t.Fields.Code(w)
	w.IndentOut()
	w.W("}\n\n")

	// func
	for _, method := range t.Methods {
		method.Code(w)
	}
}

//--------------------------------------------------------------------------------------------------------------//
// function

type Function struct {
	StructName string
	StructType string
	FuncName   string
	InlineCode string

	Args   *Vars
	Rets   *Vars
	Consts *Consts
}

func (t *Function) AddArg(item *Var) {
	if t.Args == nil {
		t.Args = &Vars{}
		t.Args.Init(VarScopeFuncArg)
	}
	t.Args.Add(item)
}

func (t *Function) AddRet(item *Var) {
	if t.Rets == nil {
		t.Rets = &Vars{}
		t.Rets.Init(VarScopeFuncRet)
	}
	t.Rets.Add(item)
}

func (t *Function) AddConst(item *Const) {
	if t.Consts == nil {
		t.Consts = &Consts{}
	}
	t.Consts.Add(item)
}

func (t *Function) Code(w *Writer) {
	w.W("func ")

	if t.StructType != "" {
		w.N("(%s %s) ", t.StructName, t.StructType)
	}

	w.N("%s(", t.FuncName)
	if t.Args != nil {
		w.N("\n")
		w.IndentIn()
		t.Args.Code(w)
		w.IndentOut()
	}
	w.W(")")

	if t.Rets != nil {
		w.N(" ") // 인자와 리턴() 간에 빈칸 1개 추가
		t.Rets.Code(w)
	}

	// 함수 내용 출력
	w.N(" {\n")
	w.IndentIn()
	if t.Consts != nil {
		t.Consts.Code(w)
	}

	// 코드 첫 빈줄 제거
	if len(t.InlineCode) != 0 && t.InlineCode[0:1] == "\n" {
		t.InlineCode = t.InlineCode[1:]
	}

	// 코드 끝 빈줄 제거
	if len(t.InlineCode) != 0 && t.InlineCode[len(t.InlineCode)-1:] == "\n" {
		t.InlineCode = t.InlineCode[:len(t.InlineCode)-1]
	}

	w.W("%s\n", strings.ReplaceAll(t.InlineCode, "\n", "\n"+w.Indent()))
	w.IndentOut()
	w.W("}\n\n")
}

//--------------------------------------------------------------------------------------------------------------//
// var

type VarScope int8

const (
	VarScopeGlobal VarScope = iota + 1
	VarScopeStructField
	VarScopeFuncArg
	VarScopeFuncRet
)

type Vars struct {
	Scope VarScope
	Items []*Var
}

func (t *Vars) Init(scope VarScope) {
	t.Scope = scope
}

func (t *Vars) Add(item *Var) {
	if t.Items == nil {
		t.Items = make([]*Var, 0, 10)
	}
	t.Items = append(t.Items, item)
}

func (t *Vars) maxLenField() int {
	var max int
	for _, pt := range t.Items {
		if max < len(pt.Name) {
			max = len(pt.Name)
		}
	}
	return max
}

func (t *Vars) Code(w *Writer) {
	switch t.Scope {
	case VarScopeGlobal:
		w.W("var (\n")
		w.IndentIn()
		for _, pt := range t.Items {
			w.W("")
			pt.Code(w)
			w.N("\n")
		}
		w.IndentOut()
		w.W(")\n")
	case VarScopeStructField:
		for _, pt := range t.Items {
			// set blanks
			pt.Type = strings.Repeat(" ", t.maxLenField()-len(pt.Name)) + pt.Type
			w.W("")
			pt.Code(w)
			w.N("\n")
		}
	case VarScopeFuncArg:
		for i := 0; i < len(t.Items); i++ {
			w.W("")
			t.Items[i].Code(w)
			w.N(",\n")
		}
	case VarScopeFuncRet:
		if len(t.Items) == 0 {
			break
		} else if len(t.Items) == 1 && t.Items[0].Name == "" {
			t.Items[0].Code(w)
			break
		}

		w.N("(\n")
		w.IndentIn()
		for i := 0; i < len(t.Items); i++ {
			w.W("")
			t.Items[i].Code(w)
			w.N(",\n")
		}
		w.IndentOut()
		w.W(")")
	default:
		log.Fatalf("input undefined type of scope - %v", t.Scope)
	}
}

//--------------------------------------------------------------------------------------------------------------//
// var

type Var struct {
	Name string
	Type string
}

func (t *Var) Code(w *Writer) {
	if t.Name != "" {
		w.N("%s %s", t.Name, t.Type)
	} else {
		w.N("%s", t.Type)
	}
}

//--------------------------------------------------------------------------------------------------------------//
// const

type Consts struct {
	items []*Const
}

func (t *Consts) Add(item *Const) {
	if t.items == nil {
		t.items = make([]*Const, 0, 10)
	}
	t.items = append(t.items, item)
}

func (t *Consts) Code(w *Writer) {
	if len(t.items) > 0 {
		w.W("const (\n")
		w.IndentIn()
		for _, pt := range t.items {
			pt.Code(w)
		}
		w.IndentOut()
		w.W(")\n")
	}
}

//--------------------------------------------------------------------------------------------------------------//
// const - item

type Const struct {
	Name  string
	Type  string // type 이 없으면 value 로 형이 결정 됨
	Value string // value 가 없으면 상위 개체가 iota 일 수 있음
}

func (t *Const) Code(w *Writer) {
	if t.Value == "" {
		w.W("%s %s\n", t.Name, t.Type) // iota
	} else {
		w.W("%s %s = %s\n", t.Name, t.Type, t.Value)
	}
}

//--------------------------------------------------------------------------------------------------------------//
// imports

type Imports struct {
	items []*Import
}

func (t *Imports) Add(item *Import) {
	if t.items == nil {
		t.items = make([]*Import, 0, 10)
	}
	t.items = append(t.items, item)
}

func (t *Imports) Code(w *Writer) {
	if len(t.items) > 0 {
		w.W("import (\n")
		w.IndentIn()
		for _, pt := range t.items {
			pt.Code(w)
		}
		w.IndentOut()
		w.W(")\n\n")
	}
}

//--------------------------------------------------------------------------------------------------------------//
// import

type Import struct {
	Path  string
	Alias string // 기본 ""(없음)  or . or  임의이름
}

func (t *Import) Code(w *Writer) {
	if t.Alias != "" {
		w.W("%s \"%s\"\n", t.Alias, t.Path)
	} else {
		w.W("\"%s\"\n", t.Path)
	}
}
