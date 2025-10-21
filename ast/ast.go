package ast

import (
	"bytes"
	"strings"
)

// Node represents any node in the AST
type Node interface {
	String() string
}

// Statement represents statement nodes
type Statement interface {
	Node
	statementNode()
}

// Expression represents expression nodes
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root of every AST
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Variable Declaration
type VarStatement struct {
	Names      []*Identifier // support multiple variable assignment (backward compatible with single variable)
	Type       *TypeAnnotation
	Value      Expression
	IsConstant bool
}

func (vs *VarStatement) statementNode() {}
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	if vs.IsConstant {
		out.WriteString("const ")
	} else {
		out.WriteString("var ")
	}
	// Handle multiple names
	names := []string{}
	for _, name := range vs.Names {
		names = append(names, name.String())
	}
	out.WriteString(strings.Join(names, ", "))
	if vs.Type != nil {
		out.WriteString(": ")
		out.WriteString(vs.Type.String())
	}
	out.WriteString(" = ")
	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}
	return out.String()
}

// Function Declaration
type FnStatement struct {
	Name       *Identifier
	Parameters []*Parameter
	ReturnType *TypeAnnotation
	Body       *BlockStatement
	WhereBlock *WhereBlock
	Receiver   *TypeAnnotation // for methods like Person.greet()
}

func (fs *FnStatement) statementNode() {}
func (fs *FnStatement) String() string {
	var out bytes.Buffer
	out.WriteString("fn ")
	if fs.Receiver != nil {
		out.WriteString(fs.Receiver.String())
		out.WriteString(".")
	}
	out.WriteString(fs.Name.String())
	out.WriteString("(")
	params := []string{}
	for _, p := range fs.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	if fs.ReturnType != nil {
		out.WriteString(": ")
		out.WriteString(fs.ReturnType.String())
	}
	out.WriteString(" ::")
	if fs.Body != nil {
		out.WriteString(fs.Body.String())
	}
	out.WriteString("end")
	return out.String()
}

// TODO: still not sure about struct
// Struct Declaration
type StructStatement struct {
	Name   *Identifier
	Fields []*StructField
}

func (ss *StructStatement) statementNode() {}
func (ss *StructStatement) String() string {
	var out bytes.Buffer
	out.WriteString("struct ")
	out.WriteString(ss.Name.String())
	out.WriteString(" ::")
	for _, field := range ss.Fields {
		out.WriteString("\n  ")
		out.WriteString(field.String())
	}
	out.WriteString("\nend")
	return out.String()
}

// Type Declaration
type TypeStatement struct {
	Name *Identifier
	Type *TypeAnnotation
}

// Module Declaration
type ModuleStatement struct {
	Name *Identifier
	Body *BlockStatement
}

// Using Statement (for external module imports)
type UsingStatement struct {
	Path  *StringLiteral
	Alias *Identifier // Optional alias for module
}

// Component Declaration (UI component)
type ComponentStatement struct {
	Name       *Identifier
	Parameters []*Parameter
	Body       *ComponentBody
}

func (cs *ComponentStatement) statementNode() {}
func (cs *ComponentStatement) String() string {
	var out bytes.Buffer
	out.WriteString("component ")
	out.WriteString(cs.Name.String())
	out.WriteString("(")
	params := []string{}
	for _, p := range cs.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ::")
	if cs.Body != nil {
		out.WriteString(cs.Body.String())
	}
	out.WriteString("\nend")
	return out.String()
}

// ComponentBody holds both state (var declarations) and UI tree
type ComponentBody struct {
	Statements []Statement // var declarations, assignments
	Root       *UIElement  // root UI element (e.g., Window)
}

func (cb *ComponentBody) String() string {
	var out bytes.Buffer
	for _, stmt := range cb.Statements {
		out.WriteString("\n  ")
		out.WriteString(stmt.String())
	}
	if cb.Root != nil {
		out.WriteString("\n  ")
		out.WriteString(cb.Root.String())
	}
	return out.String()
}

// UIElement represents a UI element (Window, VBox, Text, Button, etc.)
type UIElement struct {
	Type       *Identifier            // Window, VBox, Button, etc.
	Properties map[string]Expression  // title: "Hello", width: 400px, onClick: fn() :: ... end
	Children   []*UIElement           // nested UI elements
}

func (ue *UIElement) expressionNode() {}
func (ue *UIElement) String() string {
	var out bytes.Buffer
	out.WriteString(ue.Type.String())
	out.WriteString(" {")

	// Properties
	if len(ue.Properties) > 0 {
		out.WriteString("\n")
		for key, value := range ue.Properties {
			out.WriteString("    ")
			out.WriteString(key)
			out.WriteString(": ")
			out.WriteString(value.String())
			out.WriteString(",\n")
		}
	}

	// Children
	for _, child := range ue.Children {
		out.WriteString("    ")
		out.WriteString(child.String())
		out.WriteString("\n")
	}

	out.WriteString("  }")
	return out.String()
}

func (ts *TypeStatement) statementNode() {}
func (ts *TypeStatement) String() string {
	var out bytes.Buffer
	out.WriteString("type ")
	out.WriteString(ts.Name.String())
	out.WriteString(" = ")
	out.WriteString(ts.Type.String())
	return out.String()
}

func (ms *ModuleStatement) statementNode() {}
func (ms *ModuleStatement) String() string {
	var out bytes.Buffer
	out.WriteString("module ")
	out.WriteString(ms.Name.String())
	out.WriteString(" ::")
	if ms.Body != nil {
		out.WriteString(ms.Body.String())
	}
	out.WriteString("\nend")
	return out.String()
}

func (us *UsingStatement) statementNode() {}
func (us *UsingStatement) String() string {
	var out bytes.Buffer
	out.WriteString("using ")
	if us.Path != nil {
		out.WriteString(us.Path.String())
	}
	if us.Alias != nil {
		out.WriteString(" as ")
		out.WriteString(us.Alias.String())
	}
	return out.String()
}

// If Statement
type IfStatement struct {
	Condition Expression
	ThenBlock *BlockStatement
	ElseIfs   []*ElseIfClause
	ElseBlock *BlockStatement
}

func (ifs *IfStatement) statementNode() {}
func (ifs *IfStatement) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(ifs.Condition.String())
	out.WriteString(" ::")
	out.WriteString(ifs.ThenBlock.String())
	for _, elif := range ifs.ElseIfs {
		out.WriteString("else if ")
		out.WriteString(elif.Condition.String())
		out.WriteString(" ::")
		out.WriteString(elif.Block.String())
	}
	if ifs.ElseBlock != nil {
		out.WriteString("else ::")
		out.WriteString(ifs.ElseBlock.String())
	}
	out.WriteString("end")
	return out.String()
}

// Case Statement
type CaseStatement struct {
	Expression Expression
	Branches   []*CaseBranch
}

func (cs *CaseStatement) statementNode() {}
func (cs *CaseStatement) String() string {
	var out bytes.Buffer
	out.WriteString("case ")
	out.WriteString(cs.Expression.String())
	out.WriteString(" ::")
	for _, branch := range cs.Branches {
		out.WriteString("\n  ")
		out.WriteString(branch.String())
	}
	out.WriteString("\nend")
	return out.String()
}

// TODO: case statement and case expression have the code maybe i can remove this duplication
// Case Expression (like case statement but returns a value)
type CaseExpression struct {
	Expression Expression
	Branches   []*CaseBranch
}

func (ce *CaseExpression) expressionNode() {}
func (ce *CaseExpression) String() string {
	var out bytes.Buffer
	out.WriteString("case ")
	out.WriteString(ce.Expression.String())
	out.WriteString(" ::")
	for _, branch := range ce.Branches {
		out.WriteString("\n  ")
		out.WriteString(branch.String())
	}
	out.WriteString("\nend")
	return out.String()
}

// For Statement
type ForStatement struct {
	Variable *Identifier
	Index    *Identifier // optional, for index, value syntax
	Iterable Expression
	Body     *BlockStatement
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	if fs.Index != nil {
		out.WriteString(fs.Index.String())
		out.WriteString(", ")
	}
	out.WriteString(fs.Variable.String())
	out.WriteString(" in ")
	out.WriteString(fs.Iterable.String())
	out.WriteString(" ::")
	out.WriteString(fs.Body.String())
	out.WriteString("end")
	return out.String()
}

// Check Block
type CheckStatement struct {
	Statements []Statement  // statements (e.g. var declarations) before assertions
	Assertions []*Assertion
	Label      string // optional label for test group
}

func (cs *CheckStatement) statementNode() {}
func (cs *CheckStatement) String() string {
	var out bytes.Buffer
	out.WriteString("check")
	if cs.Label != "" {
		out.WriteString(" \"")
		out.WriteString(cs.Label)
		out.WriteString("\"")
	}
	out.WriteString(" ::")
	for _, stmt := range cs.Statements {
		out.WriteString("\n  ")
		out.WriteString(stmt.String())
	}
	for _, assertion := range cs.Assertions {
		out.WriteString("\n  ")
		out.WriteString(assertion.String())
	}
	out.WriteString("\nend")
	return out.String()
}

// Return Statement
type ReturnStatement struct {
	Values []Expression // support multiple return values (backward compatible with single value)
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString("return")
	if len(rs.Values) > 0 {
		out.WriteString(" ")
		values := []string{}
		for _, v := range rs.Values {
			values = append(values, v.String())
		}
		out.WriteString(strings.Join(values, ", "))
	}
	return out.String()
}

// Break Statement
type BreakStatement struct {
}

func (bs *BreakStatement) statementNode() {}
func (bs *BreakStatement) String() string {
	return "break"
}

// Expression Statement
type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// Block Statement
type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString("\n  ")
		out.WriteString(s.String())
	}
	return out.String()
}

// Supporting structures

// Parameter represents function parameters
type Parameter struct {
	Name *Identifier
	Type *TypeAnnotation
}

func (p *Parameter) String() string {
	var out bytes.Buffer
	out.WriteString(p.Name.String())
	if p.Type != nil {
		out.WriteString(": ")
		out.WriteString(p.Type.String())
	}
	return out.String()
}

// StructField represents struct field definitions
type StructField struct {
	Name *Identifier
	Type *TypeAnnotation
}

func (sf *StructField) String() string {
	return sf.Name.String() + ": " + sf.Type.String()
}

// ElseIfClause represents else if clauses
type ElseIfClause struct {
	Condition Expression
	Block     *BlockStatement
}

// CaseBranch represents case statement branches
type CaseBranch struct {
	Pattern Expression // could be literal, identifier, or wildcard
	Result  Expression
}

func (cb *CaseBranch) String() string {
	return cb.Pattern.String() + " => " + cb.Result.String()
}

// WhereBlock represents function test blocks
type WhereBlock struct {
	Statements []Statement  // statements (e.g. var declarations) before assertions
	Assertions []*Assertion
}

func (wb *WhereBlock) String() string {
	var out bytes.Buffer
	out.WriteString("where ::")
	for _, stmt := range wb.Statements {
		out.WriteString("\n  ")
		out.WriteString(stmt.String())
	}
	for _, assertion := range wb.Assertions {
		out.WriteString("\n  ")
		out.WriteString(assertion.String())
	}
	out.WriteString("\nend")
	return out.String()
}

// Assertion represents test assertions
type Assertion struct {
	Left     Expression
	Operator string // "is", "isA", "contains", etc.
	Right    Expression
}

func (a *Assertion) String() string {
	if a.Right != nil {
		return a.Left.String() + " " + a.Operator + " " + a.Right.String()
	}
	// Unary assertion (no right operand)
	return a.Left.String() + " " + a.Operator
}

// Expression implementations

// Identifier
type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value }

// Number Literal
type NumberLiteral struct {
	Value string
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string  { return nl.Value }

// String Literal
type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return "\"" + sl.Value + "\"" }

// Interpolated String - string with embedded expressions like "Hello #{name}"
type InterpolatedString struct {
	Parts []Expression // Mix of StringLiteral and other expressions
}

func (is *InterpolatedString) expressionNode() {}
func (is *InterpolatedString) String() string {
	var out bytes.Buffer
	out.WriteString("\"")
	for _, part := range is.Parts {
		if strLit, ok := part.(*StringLiteral); ok {
			out.WriteString(strLit.Value)
		} else {
			out.WriteString("#{")
			out.WriteString(part.String())
			out.WriteString("}")
		}
	}
	out.WriteString("\"")
	return out.String()
}

// Boolean Literal
type BooleanLiteral struct {
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) String() string {
	if bl.Value {
		return "true"
	}
	return "false"
}

// Nil Literal
type NilLiteral struct{}

func (nl *NilLiteral) expressionNode() {}
func (nl *NilLiteral) String() string {
	return "nil"
}

// Array Literal
type ArrayLiteral struct {
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range al.Elements {
		elements = append(elements, e.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// Map Literal
type MapLiteral struct {
	Pairs []MapPair
}

type MapPair struct {
	Key   Expression
	Value Expression
}

func (ml *MapLiteral) expressionNode() {}
func (ml *MapLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, p := range ml.Pairs {
		pairs = append(pairs, p.Key.String()+": "+p.Value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// Function Literal (anonymous functions)
type FunctionLiteral struct {
	Parameters []*Parameter
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ::")
	out.WriteString(fl.Body.String())
	out.WriteString("end")
	return out.String()
}

// Prefix Expression
type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// Infix Expression
type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// Call Expression
type CallExpression struct {
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// Index Expression
type IndexExpression struct {
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

// Dot Expression (property access)
type DotExpression struct {
	Left     Expression
	Property *Identifier
}

func (de *DotExpression) expressionNode() {}
func (de *DotExpression) String() string {
	var out bytes.Buffer
	out.WriteString(de.Left.String())
	out.WriteString(".")
	out.WriteString(de.Property.String())
	return out.String()
}

// Assignment Expression
type AssignmentExpression struct {
	Left  Expression
	Value Expression
}

func (ae *AssignmentExpression) expressionNode() {}
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ae.Left.String())
	out.WriteString(" = ")
	out.WriteString(ae.Value.String())
	return out.String()
}

// Range Expression
type RangeExpression struct {
	Start     Expression
	End       Expression
	Inclusive bool // true for ..., false for ..
}

func (re *RangeExpression) expressionNode() {}
func (re *RangeExpression) String() string {
	var out bytes.Buffer
	out.WriteString(re.Start.String())
	if re.Inclusive {
		out.WriteString("...")
	} else {
		out.WriteString("..")
	}
	out.WriteString(re.End.String())
	return out.String()
}

// Type Annotation
type TypeAnnotation struct {
	Name       string
	Parameters []*TypeAnnotation // for generic types like Array[String]
}

func (ta *TypeAnnotation) String() string {
	var out bytes.Buffer
	out.WriteString(ta.Name)
	if len(ta.Parameters) > 0 {
		params := []string{}
		for _, p := range ta.Parameters {
			params = append(params, p.String())
		}
		out.WriteString("[")
		out.WriteString(strings.Join(params, ", "))
		out.WriteString("]")
	}
	return out.String()
}
