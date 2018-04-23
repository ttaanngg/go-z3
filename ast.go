package z3

// #include <stdlib.h>
// #include "go-z3.h"
import "C"

// AST represents an AST value in Z3.
//
// AST memory management is automatically managed by the Context it
// is contained within. When the Context is freed, so are the AST nodes.
type AST struct {
	rawCtx C.Z3_context
	rawAST C.Z3_ast
}

// String returns a human-friendly string version of the AST.
func (a *AST) String() string {
	return C.GoString(C.Z3_ast_to_string(a.rawCtx, a.rawAST))
}

// DeclName returns the name of a declaration. The AST value must be a
// func declaration for this to work.
func (a *AST) DeclName() *Symbol {
	return &Symbol{
		rawCtx: a.rawCtx,
		rawSymbol: C.Z3_get_decl_name(
			a.rawCtx, C.Z3_to_func_decl(a.rawCtx, a.rawAST)),
	}
}

//-------------------------------------------------------------------
// Var, Literal Creation
//-------------------------------------------------------------------

// Const declares a variable. It is called "Const" since internally
// this is equivalent to create a function that always returns a constant
// value. From an initial user perspective this may be confusing but go-z3
// is following identical naming convention.
func (c *Context) Const(s *Symbol, typ *Sort) *AST {
	return &AST{
		rawCtx: c.raw,
		rawAST: C.Z3_mk_const(c.raw, s.rawSymbol, typ.rawSort),
	}
}

// Int creates an integer type.
//
// Maps: Z3_mk_int
func (c *Context) Int(v int, typ *Sort) *AST {
	return &AST{
		rawCtx: c.raw,
		rawAST: C.Z3_mk_int(c.raw, C.int(v), typ.rawSort),
	}
}

// Float creates an float type.
//
// Maps: Z3_mk_real
func (c *Context) Float(v float64) *AST {
	//TODO: test if this could work
	return &AST{
		rawCtx: c.raw,
		rawAST: C.Z3_mk_real(c.raw, C.int(v), C.int(1)),
	}
}

// Str creates an string type.
//
// Maps: Z3_mk_string
func (c *Context) Str(str string) *AST {
	//TODO: test if this could work
	return &AST{
		rawCtx: c.raw,
		rawAST: C.Z3_mk_string(c.raw, C.CString(str)),
	}
}

// True creates the value "true".
//
// Maps: Z3_mk_true
func (c *Context) True() *AST {
	return &AST{
		rawCtx: c.raw,
		rawAST: C.Z3_mk_true(c.raw),
	}
}

// False creates the value "false".
//
// Maps: Z3_mk_false
func (c *Context) False() *AST {
	return &AST{
		rawCtx: c.raw,
		rawAST: C.Z3_mk_false(c.raw),
	}
}

// RealSeq returns the seq type number.
func (c *Context) RealSet(reals ...float64) *AST {
	set := &AST{
		rawCtx: c.raw,
		rawAST: C.Z3_mk_empty_set(
			c.raw,
			c.RealSort().rawSort,
		),
	}
	for _, content := range reals {
		C.Z3_mk_set_add(
			c.raw,
			set.rawAST,
			c.Float(content).rawAST,
		)
	}
	return set
}

// StringSeqSort returns the seq type string.
func (c *Context) StringSet(strings ...string) *AST {
	set := &AST{
		rawCtx: c.raw,
		rawAST: C.Z3_mk_empty_set(
			c.raw,
			c.StringSort().rawSort,
		),
	}
	for _, content := range strings {
		C.Z3_mk_set_add(
			c.raw,
			set.rawAST,
			c.Str(content).rawAST,
		)
	}
	return set
}

//-------------------------------------------------------------------
// Value Readers
//-------------------------------------------------------------------

// Int gets the integer value of this AST. The value must be able to fit
// into a machine integer.
func (a *AST) Int() int {
	var dst C.int
	C.Z3_get_numeral_int(a.rawCtx, a.rawAST, &dst)
	return int(dst)
}
