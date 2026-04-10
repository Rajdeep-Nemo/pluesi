package evaluator

import (
	"fmt"
	"testing"

	"github.com/Rajdeep-Nemo/sugarglaze/internal/lexer"
	"github.com/Rajdeep-Nemo/sugarglaze/internal/object"
	"github.com/Rajdeep-Nemo/sugarglaze/internal/parser"
	"github.com/Rajdeep-Nemo/sugarglaze/internal/token"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		objType  object.ObjectType
	}{
		{"5", 5, object.I32_OBJ},
		{"10", 10, object.I32_OBJ},
		{"-5", -5, object.I32_OBJ},
		{"-10", -10, object.I32_OBJ},
		{"5 + 5 + 5 + 5 - 10", 10, object.I32_OBJ},
		{"2 * 2 * 2 * 2 * 2", 32, object.I32_OBJ},
		{"-50 + 100 + -50", 0, object.I32_OBJ},
		{"5 * 2 + 10", 20, object.I32_OBJ},
		{"5 + 2 * 10", 25, object.I32_OBJ},
		{"20 + 2 * -10", 0, object.I32_OBJ},
		{"50 / 2 * 2 + 10", 60, object.I32_OBJ},
		{"2 * (5 + 10)", 30, object.I32_OBJ},
		{"3 * 3 * 3 + 10", 37, object.I32_OBJ},
		{"3 * (3 * 3) + 10", 37, object.I32_OBJ},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50, object.I32_OBJ},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected, tt.objType)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		objType  object.ObjectType
	}{
		{"5.5", 5.5, object.F32_OBJ},
		{"-5.5", -5.5, object.F32_OBJ},
		{"5.0 + 5.5", 10.5, object.F32_OBJ},
		{"10.5 - 5.0", 5.5, object.F32_OBJ},
		{"2.0 * 2.5", 5.0, object.F32_OBJ},
		{"10.0 / 2.0", 5.0, object.F32_OBJ},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected, tt.objType)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		objType  object.ObjectType
	}{
		{"let a = 5 \n a", 5, object.I32_OBJ},
		{"let a = 5 * 5 \n a", 25, object.I32_OBJ},
		{"let a = 5 \n let b = a \n b", 5, object.I32_OBJ},
		{"let a = 5 \n let b = a \n let c = a + b + 5 \n c", 15, object.I32_OBJ},
		{"let a: i8 = 10 \n a", 10, object.I8_OBJ},
		{"let a: i16 = 200 \n a", 200, object.I16_OBJ},
		{"let a: i64 = 50 \n a", 50, object.I64_OBJ},
		{"let a: u8 = 255 \n a", 255, object.U8_OBJ},
		{"let a: f64 = 5 \n a", 5, object.F64_OBJ},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected, tt.objType)
	}
}

func TestConstStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		objType  object.ObjectType
	}{
		// FIX: Added the mandatory type hints for const to pass the parser check
		{"const a: i32 = 5 \n a", 5, object.I32_OBJ},
		{"const a: i8 = 10 \n a", 10, object.I8_OBJ},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected, tt.objType)
	}
}

func TestAssignStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		objType  object.ObjectType
	}{
		{"let a = 5 \n a = 10 \n a", 10, object.I32_OBJ},
		{"let a = 5 \n a += 10 \n a", 15, object.I32_OBJ},
		{"let a = 10 \n a -= 5 \n a", 5, object.I32_OBJ},
		{"let a = 5 \n a *= 2 \n a", 10, object.I32_OBJ},
		{"let a = 10 \n a /= 2 \n a", 5, object.I32_OBJ},
		{"let a: i8 = 5 \n a = 10 \n a", 10, object.I8_OBJ},
		{"let a: i8 = 5 \n a += 10 \n a", 15, object.I8_OBJ},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected, tt.objType)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"let a: i8 = 300",
			"value 300 out of bounds for i8",
		},
		{
			"let a: u8 = -5",
			"value -5 out of bounds for u8",
		},
		{
			"let a: i32 \n a",
			"cannot access uninitialized variable 'a'",
		},
		{
			"const a: i32 = 5 \n a = 10",
			"cannot reassign to const variable 'a'",
		},
		{
			"let a: i32 = 5 \n a = 5.5",
			"type mismatch: cannot assign 'f32' to variable 'a' (expected 'i32')",
		},
		{
			"let a: f32 = 5.5 \n let b: i32 = a",
			// FIX: Now uses the professional variable-context error string!
			"type mismatch: cannot assign 'f32' to variable 'b' (expected 'i32')",
		},
		{
			"let a: i8 = 100 \n a += 100",
			"value 200 out of bounds for i8",
		},
		{
			"let a: u8 = 5 \n a -= 10",
			"value -5 out of bounds for u8",
		},
		{
			"-true",
			"unknown operator: -bool",
		},
		{
			"let a: u32 = 10 \n let b = -a",
			"cannot negate unsigned integer type u32",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if evaluated == nil {
			t.Errorf("expected error object, got nil for input: %q (Check parser for syntax errors)", tt.input)
			continue
		}

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%v) for input: %q",
				evaluated, evaluated, tt.input)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func testEval(input string) object.Object {
	s := lexer.InitScanner(input)
	var tokens []token.Token

	for {
		tok := s.ScanToken()
		tokens = append(tokens, tok)
		if tok.Type == token.END_OF_FILE {
			break
		}
	}

	p := parser.New(tokens)
	program := p.ParseProgram()

	if program == nil || len(program.Statements) == 0 {
		fmt.Printf("Warning: Parser returned empty program for input: %q\n", input)
		return nil
	}

	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64, expectedType object.ObjectType) bool {
	if obj == nil {
		t.Errorf("object is nil. Expected %d (%s)", expected, expectedType)
		return false
	}

	if obj.Type() != expectedType {
		t.Errorf("object is not %s. got=%T (%s)", expectedType, obj, obj.Type())
		return false
	}

	var val int64
	switch v := obj.(type) {
	case *object.Int8:
		val = int64(v.Value)
	case *object.Int16:
		val = int64(v.Value)
	case *object.Int32:
		val = int64(v.Value)
	case *object.Int64:
		val = v.Value
	case *object.Uint8:
		val = int64(v.Value)
	case *object.Uint16:
		val = int64(v.Value)
	case *object.Uint32:
		val = int64(v.Value)
	case *object.Uint64:
		val = int64(v.Value)
	case *object.Float64:
		val = int64(v.Value)
	default:
		t.Errorf("object is not an integer type. got=%T", obj)
		return false
	}

	if val != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", val, expected)
		return false
	}

	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64, expectedType object.ObjectType) bool {
	if obj == nil {
		t.Errorf("object is nil. Expected %f (%s)", expected, expectedType)
		return false
	}

	if obj.Type() != expectedType {
		t.Errorf("object is not %s. got=%T (%s)", expectedType, obj, obj.Type())
		return false
	}

	var val float64
	switch v := obj.(type) {
	case *object.Float32:
		val = float64(v.Value)
	case *object.Float64:
		val = v.Value
	default:
		t.Errorf("object is not a float type. got=%T", obj)
		return false
	}

	if val != expected {
		t.Errorf("object has wrong value. got=%f, want=%f", val, expected)
		return false
	}

	return true
}

// ParseBool Tests
func TestParseBool_True(t *testing.T) {
	v := ParseBool("true")
	if v == nil || *v != true {
		t.Errorf("Expected true, got %v", v)
	}
}

func TestParseBool_False(t *testing.T) {
	v := ParseBool("false")
	if v == nil || *v != false {
		t.Errorf("Expected false, got %v", v)
	}
}

func TestParseBool_Trimmed(t *testing.T) {
	v := ParseBool("  true  ")
	if v == nil || *v != true {
		t.Errorf("Expected true after trim, got %v", v)
	}
}

func TestParseBool_Invalid(t *testing.T) {
	v := ParseBool("invalid")
	if v != nil {
		t.Errorf("Expected nil, got %v", *v)
	}
}

func TestParseBool_Empty(t *testing.T) {
	v := ParseBool("")
	if v != nil {
		t.Errorf("Expected nil, got %v", *v)
	}
}

// ParseChar Tests
func TestParseChar_Valid(t *testing.T) {
	v := ParseChar("a")
	if v == nil || *v != 'a' {
		t.Errorf("Expected 'a', got %v", v)
	}
}

func TestParseChar_Trimmed(t *testing.T) {
	v := ParseChar("  a  ")
	if v == nil || *v != 'a' {
		t.Errorf("Expected 'a' after trim, got %v", v)
	}
}

func TestParseChar_Multiple(t *testing.T) {
	v := ParseChar("ab")
	if v != nil {
		t.Errorf("Expected nil, got %v", *v)
	}
}

func TestParseChar_Empty(t *testing.T) {
	v := ParseChar("")
	if v != nil {
		t.Errorf("Expected nil, got %v", *v)
	}
}

// ParseU8 Tests
func TestParseU8_Valid(t *testing.T) {
	v := ParseU8("255")
	if v == nil || *v != 255 {
		t.Errorf("Expected 255, got %v", v)
	}
}

func TestParseU8_Zero(t *testing.T) {
	v := ParseU8("0")
	if v == nil || *v != 0 {
		t.Errorf("Expected 0, got %v", v)
	}
}

func TestParseU8_Trimmed(t *testing.T) {
	v := ParseU8("  42  ")
	if v == nil || *v != 42 {
		t.Errorf("Expected 42 after trim, got %v", v)
	}
}

func TestParseU8_Overflow(t *testing.T) {
	v := ParseU8("256")
	if v != nil {
		t.Errorf("Expected nil on overflow, got %v", *v)
	}
}

func TestParseU8_Negative(t *testing.T) {
	v := ParseU8("-1")
	if v != nil {
		t.Errorf("Expected nil on negative, got %v", *v)
	}
}

func TestParseU8_Partial(t *testing.T) {
	v := ParseU8("123abc")
	if v != nil {
		t.Errorf("Expected nil on partial parse, got %v", *v)
	}
}

func TestParseU8_Empty(t *testing.T) {
	v := ParseU8("")
	if v != nil {
		t.Errorf("Expected nil on empty, got %v", *v)
	}
}

// ParseU16 Tests
func TestParseU16_Valid(t *testing.T) {
	v := ParseU16("65535")
	if v == nil || *v != 65535 {
		t.Errorf("Expected 65535, got %v", v)
	}
}

func TestParseU16_Zero(t *testing.T) {
	v := ParseU16("0")
	if v == nil || *v != 0 {
		t.Errorf("Expected 0, got %v", v)
	}
}

func TestParseU16_Trimmed(t *testing.T) {
	v := ParseU16("  100  ")
	if v == nil || *v != 100 {
		t.Errorf("Expected 100 after trim, got %v", v)
	}
}

func TestParseU16_Overflow(t *testing.T) {
	v := ParseU16("65536")
	if v != nil {
		t.Errorf("Expected nil on overflow, got %v", *v)
	}
}

func TestParseU16_Negative(t *testing.T) {
	v := ParseU16("-1")
	if v != nil {
		t.Errorf("Expected nil on negative, got %v", *v)
	}
}

func TestParseU16_Partial(t *testing.T) {
	v := ParseU16("123abc")
	if v != nil {
		t.Errorf("Expected nil on partial parse, got %v", *v)
	}
}

func TestParseU16_Empty(t *testing.T) {
	v := ParseU16("")
	if v != nil {
		t.Errorf("Expected nil on empty, got %v", *v)
	}
}

// ParseU32 Tests
func TestParseU32_Valid(t *testing.T) {
	v := ParseU32("4294967295")
	if v == nil || *v != 4294967295 {
		t.Errorf("Expected 4294967295, got %v", v)
	}
}

func TestParseU32_Zero(t *testing.T) {
	v := ParseU32("0")
	if v == nil || *v != 0 {
		t.Errorf("Expected 0, got %v", v)
	}
}

func TestParseU32_Trimmed(t *testing.T) {
	v := ParseU32("  100  ")
	if v == nil || *v != 100 {
		t.Errorf("Expected 100 after trim, got %v", v)
	}
}

func TestParseU32_Overflow(t *testing.T) {
	v := ParseU32("4294967296")
	if v != nil {
		t.Errorf("Expected nil on overflow, got %v", *v)
	}
}

func TestParseU32_Negative(t *testing.T) {
	v := ParseU32("-1")
	if v != nil {
		t.Errorf("Expected nil on negative, got %v", *v)
	}
}

func TestParseU32_Partial(t *testing.T) {
	v := ParseU32("123abc")
	if v != nil {
		t.Errorf("Expected nil on partial parse, got %v", *v)
	}
}

func TestParseU32_Empty(t *testing.T) {
	v := ParseU32("")
	if v != nil {
		t.Errorf("Expected nil on empty, got %v", *v)
	}
}

// ParseU64 Tests
func TestParseU64_Valid(t *testing.T) {
	v := ParseU64("18446744073709551615")
	if v == nil || *v != 18446744073709551615 {
		t.Errorf("Expected 18446744073709551615, got %v", v)
	}
}

func TestParseU64_Zero(t *testing.T) {
	v := ParseU64("0")
	if v == nil || *v != 0 {
		t.Errorf("Expected 0, got %v", v)
	}
}

func TestParseU64_Trimmed(t *testing.T) {
	v := ParseU64("  100  ")
	if v == nil || *v != 100 {
		t.Errorf("Expected 100 after trim, got %v", v)
	}
}

func TestParseU64_Overflow(t *testing.T) {
	v := ParseU64("18446744073709551616")
	if v != nil {
		t.Errorf("Expected nil on overflow, got %v", *v)
	}
}

func TestParseU64_Negative(t *testing.T) {
	v := ParseU64("-1")
	if v != nil {
		t.Errorf("Expected nil on negative, got %v", *v)
	}
}

func TestParseU64_Partial(t *testing.T) {
	v := ParseU64("123abc")
	if v != nil {
		t.Errorf("Expected nil on partial parse, got %v", *v)
	}
}

func TestParseU64_Empty(t *testing.T) {
	v := ParseU64("")
	if v != nil {
		t.Errorf("Expected nil on empty, got %v", *v)
	}
}

// ParseI8 Tests
func TestParseI8_Valid(t *testing.T) {
	v := ParseI8("127")
	if v == nil || *v != 127 {
		t.Errorf("Expected 127, got %v", v)
	}
}

func TestParseI8_Negative(t *testing.T) {
	v := ParseI8("-128")
	if v == nil || *v != -128 {
		t.Errorf("Expected -128, got %v", v)
	}
}

func TestParseI8_Zero(t *testing.T) {
	v := ParseI8("0")
	if v == nil || *v != 0 {
		t.Errorf("Expected 0, got %v", v)
	}
}

func TestParseI8_Trimmed(t *testing.T) {
	v := ParseI8("  42  ")
	if v == nil || *v != 42 {
		t.Errorf("Expected 42 after trim, got %v", v)
	}
}

func TestParseI8_Overflow(t *testing.T) {
	v := ParseI8("128")
	if v != nil {
		t.Errorf("Expected nil on overflow, got %v", *v)
	}
}

func TestParseI8_Underflow(t *testing.T) {
	v := ParseI8("-129")
	if v != nil {
		t.Errorf("Expected nil on underflow, got %v", *v)
	}
}

func TestParseI8_Partial(t *testing.T) {
	v := ParseI8("12abc")
	if v != nil {
		t.Errorf("Expected nil on partial parse, got %v", *v)
	}
}

func TestParseI8_Empty(t *testing.T) {
	v := ParseI8("")
	if v != nil {
		t.Errorf("Expected nil on empty, got %v", *v)
	}
}

// ParseI16 Tests
func TestParseI16_Valid(t *testing.T) {
	v := ParseI16("32767")
	if v == nil || *v != 32767 {
		t.Errorf("Expected 32767, got %v", v)
	}
}

func TestParseI16_Negative(t *testing.T) {
	v := ParseI16("-32768")
	if v == nil || *v != -32768 {
		t.Errorf("Expected -32768, got %v", v)
	}
}

func TestParseI16_Zero(t *testing.T) {
	v := ParseI16("0")
	if v == nil || *v != 0 {
		t.Errorf("Expected 0, got %v", v)
	}
}

func TestParseI16_Trimmed(t *testing.T) {
	v := ParseI16("  100  ")
	if v == nil || *v != 100 {
		t.Errorf("Expected 100 after trim, got %v", v)
	}
}

func TestParseI16_Overflow(t *testing.T) {
	v := ParseI16("32768")
	if v != nil {
		t.Errorf("Expected nil on overflow, got %v", *v)
	}
}

func TestParseI16_Underflow(t *testing.T) {
	v := ParseI16("-32769")
	if v != nil {
		t.Errorf("Expected nil on underflow, got %v", *v)
	}
}

func TestParseI16_Partial(t *testing.T) {
	v := ParseI16("12abc")
	if v != nil {
		t.Errorf("Expected nil on partial parse, got %v", *v)
	}
}

func TestParseI16_Empty(t *testing.T) {
	v := ParseI16("")
	if v != nil {
		t.Errorf("Expected nil on empty, got %v", *v)
	}
}

// ParseI32 Tests
func TestParseI32_Valid(t *testing.T) {
	v := ParseI32("2147483647")
	if v == nil || *v != 2147483647 {
		t.Errorf("Expected 2147483647, got %v", v)
	}
}

func TestParseI32_Negative(t *testing.T) {
	v := ParseI32("-2147483648")
	if v == nil || *v != -2147483648 {
		t.Errorf("Expected -2147483648, got %v", v)
	}
}

func TestParseI32_Zero(t *testing.T) {
	v := ParseI32("0")
	if v == nil || *v != 0 {
		t.Errorf("Expected 0, got %v", v)
	}
}

func TestParseI32_Trimmed(t *testing.T) {
	v := ParseI32("  100  ")
	if v == nil || *v != 100 {
		t.Errorf("Expected 100 after trim, got %v", v)
	}
}

func TestParseI32_Overflow(t *testing.T) {
	v := ParseI32("2147483648")
	if v != nil {
		t.Errorf("Expected nil on overflow, got %v", *v)
	}
}

func TestParseI32_Underflow(t *testing.T) {
	v := ParseI32("-2147483649")
	if v != nil {
		t.Errorf("Expected nil on underflow, got %v", *v)
	}
}

func TestParseI32_Partial(t *testing.T) {
	v := ParseI32("12abc")
	if v != nil {
		t.Errorf("Expected nil on partial parse, got %v", *v)
	}
}

func TestParseI32_Empty(t *testing.T) {
	v := ParseI32("")
	if v != nil {
		t.Errorf("Expected nil on empty, got %v", *v)
	}
}

// ParseI64 Tests
func TestParseI64_Valid(t *testing.T) {
	v := ParseI64("9223372036854775807")
	if v == nil || *v != 9223372036854775807 {
		t.Errorf("Expected 9223372036854775807, got %v", v)
	}
}

func TestParseI64_Negative(t *testing.T) {
	v := ParseI64("-9223372036854775808")
	if v == nil || *v != -9223372036854775808 {
		t.Errorf("Expected -9223372036854775808, got %v", v)
	}
}

func TestParseI64_Zero(t *testing.T) {
	v := ParseI64("0")
	if v == nil || *v != 0 {
		t.Errorf("Expected 0, got %v", v)
	}
}

func TestParseI64_Trimmed(t *testing.T) {
	v := ParseI64("  100  ")
	if v == nil || *v != 100 {
		t.Errorf("Expected 100 after trim, got %v", v)
	}
}

func TestParseI64_Overflow(t *testing.T) {
	v := ParseI64("9223372036854775808")
	if v != nil {
		t.Errorf("Expected nil on overflow, got %v", *v)
	}
}

func TestParseI64_Underflow(t *testing.T) {
	v := ParseI64("-9223372036854775809")
	if v != nil {
		t.Errorf("Expected nil on underflow, got %v", *v)
	}
}

func TestParseI64_Partial(t *testing.T) {
	v := ParseI64("12abc")
	if v != nil {
		t.Errorf("Expected nil on partial parse, got %v", *v)
	}
}

func TestParseI64_Empty(t *testing.T) {
	v := ParseI64("")
	if v != nil {
		t.Errorf("Expected nil on empty, got %v", *v)
	}
}

// ParseF32 Tests
func TestParseF32_Valid(t *testing.T) {
	v := ParseF32("3.14")
	if v == nil || *v != float32(3.14) {
		t.Errorf("Expected 3.14, got %v", v)
	}
}

func TestParseF32_Negative(t *testing.T) {
	v := ParseF32("-3.14")
	if v == nil || *v != float32(-3.14) {
		t.Errorf("Expected -3.14, got %v", v)
	}
}

func TestParseF32_Zero(t *testing.T) {
	v := ParseF32("0")
	if v == nil || *v != 0 {
		t.Errorf("Expected 0, got %v", v)
	}
}

func TestParseF32_Trimmed(t *testing.T) {
	v := ParseF32("  3.14  ")
	if v == nil || *v != float32(3.14) {
		t.Errorf("Expected 3.14 after trim, got %v", v)
	}
}

func TestParseF32_Overflow(t *testing.T) {
	v := ParseF32("3.4028235e+39")
	if v != nil {
		t.Errorf("Expected nil on overflow, got %v", *v)
	}
}

func TestParseF32_Partial(t *testing.T) {
	v := ParseF32("3.14abc")
	if v != nil {
		t.Errorf("Expected nil on partial parse, got %v", *v)
	}
}

func TestParseF32_Empty(t *testing.T) {
	v := ParseF32("")
	if v != nil {
		t.Errorf("Expected nil on empty, got %v", *v)
	}
}

func TestParseF32_Integer(t *testing.T) {
	v := ParseF32("42")
	if v == nil || *v != float32(42) {
		t.Errorf("Expected 42, got %v", v)
	}
}

// ParseF64 Tests
func TestParseF64_Valid(t *testing.T) {
	v := ParseF64("3.141592653589793")
	if v == nil || *v != 3.141592653589793 {
		t.Errorf("Expected 3.141592653589793, got %v", v)
	}
}

func TestParseF64_Negative(t *testing.T) {
	v := ParseF64("-3.141592653589793")
	if v == nil || *v != -3.141592653589793 {
		t.Errorf("Expected -3.141592653589793, got %v", v)
	}
}

func TestParseF64_Zero(t *testing.T) {
	v := ParseF64("0")
	if v == nil || *v != 0 {
		t.Errorf("Expected 0, got %v", v)
	}
}

func TestParseF64_Trimmed(t *testing.T) {
	v := ParseF64("  3.14  ")
	if v == nil || *v != 3.14 {
		t.Errorf("Expected 3.14 after trim, got %v", v)
	}
}

func TestParseF64_Overflow(t *testing.T) {
	v := ParseF64("1.7976931348623157e+309")
	if v != nil {
		t.Errorf("Expected nil on overflow, got %v", *v)
	}
}

func TestParseF64_Partial(t *testing.T) {
	v := ParseF64("3.14abc")
	if v != nil {
		t.Errorf("Expected nil on partial parse, got %v", *v)
	}
}

func TestParseF64_Empty(t *testing.T) {
	v := ParseF64("")
	if v != nil {
		t.Errorf("Expected nil on empty, got %v", *v)
	}
}

func TestParseF64_Integer(t *testing.T) {
	v := ParseF64("42")
	if v == nil || *v != 42 {
		t.Errorf("Expected 42, got %v", v)
	}
}
