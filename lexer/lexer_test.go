// This test case is AI generated - Gemini 3.1 Pro
package lexer

import (
	"Me/tokens"
	"testing"
)

// expectedToken is a small helper struct just for our tests
type expectedToken struct {
	expectedType   tokens.TokenType
	expectedLexeme string
	expectedLine   int
}

func TestScanner(t *testing.T) {
	// 1. Define our Table-Driven Tests
	tests := []struct {
		name     string          // Name of the test
		input    string          // The raw source code to lex
		expected []expectedToken // The exact tokens we expect to get back
	}{
		{
			name:  "Single Character Punctuation",
			input: "(){}[],.;?",
			expected: []expectedToken{
				{tokens.OPEN_PAREN, "(", 1},
				{tokens.CLOSE_PAREN, ")", 1},
				{tokens.OPEN_BRACE, "{", 1},
				{tokens.CLOSE_BRACE, "}", 1},
				{tokens.OPEN_BRACKET, "[", 1},
				{tokens.CLOSE_BRACKET, "]", 1},
				{tokens.COMMA, ",", 1},
				{tokens.DOT, ".", 1},
				{tokens.SEMICOLON, ";", 1},
				{tokens.QUESTION, "?", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Two-Character Operators",
			input: "!= == <= >= -> ..",
			expected: []expectedToken{
				{tokens.BANG_EQUAL, "!=", 1},
				{tokens.EQUAL_EQUAL, "==", 1},
				{tokens.LESS_EQUAL, "<=", 1},
				{tokens.GREATER_EQUAL, ">=", 1},
				{tokens.ARROW, "->", 1},
				{tokens.DOT_DOT, "..", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Mathematical Operators",
			input: "+ - * / %",
			expected: []expectedToken{
				{tokens.PLUS, "+", 1},
				{tokens.MINUS, "-", 1},
				{tokens.STAR, "*", 1},
				{tokens.SLASH, "/", 1},
				{tokens.PERCENT, "%", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Compound Assignment Operators",
			input: "+= -= *= /= %=",
			expected: []expectedToken{
				{tokens.PLUS_EQUAL, "+=", 1},
				{tokens.MINUS_EQUAL, "-=", 1},
				{tokens.STAR_EQUAL, "*=", 1},
				{tokens.SLASH_EQUAL, "/=", 1},
				{tokens.PERCENT_EQUAL, "%=", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Bitwise and Logical Operators",
			input: "& && | || ^ ~ << >> !",
			expected: []expectedToken{
				{tokens.BIT_AND, "&", 1},
				{tokens.AND, "&&", 1},
				{tokens.BIT_OR, "|", 1},
				{tokens.OR, "||", 1},
				{tokens.BIT_XOR, "^", 1},
				{tokens.BIT_NOT, "~", 1},
				{tokens.LEFT_SHIFT, "<<", 1},
				{tokens.RIGHT_SHIFT, ">>", 1},
				{tokens.BANG, "!", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Control Flow Keywords",
			input: "if else loop in break continue match return",
			expected: []expectedToken{
				{tokens.IF, "if", 1},
				{tokens.ELSE, "else", 1},
				{tokens.LOOP, "loop", 1},
				{tokens.IN, "in", 1},
				{tokens.BREAK, "break", 1},
				{tokens.CONTINUE, "continue", 1},
				{tokens.MATCH, "match", 1},
				{tokens.RETURN, "return", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Data Type Keywords",
			input: "i8 u32 f64 char string bool struct enum union",
			expected: []expectedToken{
				{tokens.I8, "i8", 1},
				{tokens.U32, "u32", 1},
				{tokens.F64, "f64", 1},
				{tokens.CHAR, "char", 1},
				{tokens.STRING, "string", 1},
				{tokens.BOOL, "bool", 1},
				{tokens.STRUCT, "struct", 1},
				{tokens.ENUM, "enum", 1},
				{tokens.UNION, "union", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Boolean and Nil Literals",
			input: "true false NIL",
			expected: []expectedToken{
				{tokens.TRUE, "true", 1},
				{tokens.FALSE, "false", 1},
				{tokens.NIL_LITERAL, "NIL", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Character Literals",
			input: "'a' '\\n'",
			expected: []expectedToken{
				{tokens.CHAR_LITERAL, "'a'", 1},
				{tokens.CHAR_LITERAL, "'\\n'", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Complex Identifiers",
			input: "_privateVar camelCase123 another_var",
			expected: []expectedToken{
				{tokens.IDENTIFIER, "_privateVar", 1},
				{tokens.IDENTIFIER, "camelCase123", 1},
				{tokens.IDENTIFIER, "another_var", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "String Literal and Whitespace",
			input: "   \n \"hello world\"  \n",
			expected: []expectedToken{
				{tokens.STRING_LITERAL, "\"hello world\"", 2},
				{tokens.END_OF_FILE, "", 3},
			},
		},
		{
			name:  "Unterminated String Error",
			input: "\"this string never ends",
			expected: []expectedToken{
				{tokens.ERROR_TOKEN, "Unterminated string.", 1}, // Added the period based on your earlier code, adjust if yours lacks it
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name: "Full Snippet with Comments",
			input: `
                // This is a comment
                fn calculate() i32 {
                    return 100;
                }
            `,
			expected: []expectedToken{
				{tokens.FN, "fn", 3},
				{tokens.IDENTIFIER, "calculate", 3},
				{tokens.OPEN_PAREN, "(", 3},
				{tokens.CLOSE_PAREN, ")", 3},
				{tokens.I32, "i32", 3},
				{tokens.OPEN_BRACE, "{", 3},
				{tokens.RETURN, "return", 4},
				{tokens.INT_LITERAL, "100", 4},
				{tokens.SEMICOLON, ";", 4},
				{tokens.CLOSE_BRACE, "}", 5},
				{tokens.END_OF_FILE, "", 6},
			},
		},
		{
			name: "EVIL: Alien Characters",
			// Someone drops a random @ and $ in the code
			input: "let @x = $5;",
			expected: []expectedToken{
				{tokens.LET, "let", 1},
				{tokens.ERROR_TOKEN, "Unexpected character.", 1}, // The @
				{tokens.IDENTIFIER, "x", 1},
				{tokens.EQUAL, "=", 1},
				{tokens.ERROR_TOKEN, "Unexpected character.", 1}, // The $
				{tokens.INT_LITERAL, "5", 1},                     // It should recover and keep reading!
				{tokens.SEMICOLON, ";", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name: "EVIL: EOF Triggered Mid-Escape Sequence",
			// A string that ends literally right after the escape backslash
			input: "\"This string ends with an escape \\",
			expected: []expectedToken{
				{tokens.ERROR_TOKEN, "Unterminated string after escape.", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name: "EVIL: EOF Triggered Mid-Comment",
			// A comment with no newline at the end of the file
			input: "let x = 10; // This file ends right he...",
			expected: []expectedToken{
				{tokens.LET, "let", 1},
				{tokens.IDENTIFIER, "x", 1},
				{tokens.EQUAL, "=", 1},
				{tokens.INT_LITERAL, "10", 1},
				{tokens.SEMICOLON, ";", 1},
				{tokens.END_OF_FILE, "", 1}, // It should gracefully hit EOF, not freeze
			},
		},
		{
			name: "EVIL: Whitespace Chaos",
			// Mixing tabs, spaces, Windows carriage returns (\r), and Linux newlines (\n)
			input: " \t\r\n  \n\n  let\r\n\tx",
			expected: []expectedToken{
				{tokens.LET, "let", 4}, // Should correctly track lines despite the mess
				{tokens.IDENTIFIER, "x", 5},
				{tokens.END_OF_FILE, "", 5},
			},
		},
		{
			name: "EVIL: Operator Soup (Maximal Munch Test)",
			// Lexers use "Maximal Munch" (grab the longest match possible).
			// ===>= should become `==`, `=`, `>=`
			input: "===>=",
			expected: []expectedToken{
				{tokens.EQUAL_EQUAL, "==", 1},
				{tokens.EQUAL, "=", 1},
				{tokens.GREATER_EQUAL, ">=", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name: "EVIL: Numbers glued to Identifiers",
			// Standard C-like languages separate these into two tokens
			input: "99bottles",
			expected: []expectedToken{
				{tokens.INT_LITERAL, "99", 1},
				{tokens.IDENTIFIER, "bottles", 1},
				{tokens.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "EVIL: Completely Empty File",
			input: "    \t   \n  ",
			expected: []expectedToken{
				{tokens.END_OF_FILE, "", 2}, // Should just skip whitespace and exit safely
			},
		},
	}

	// 2. Loop through all the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Initialize scanner with the test input
			scanner := InitScanner(tt.input)

			// Check every expected token
			for i, exp := range tt.expected {
				tok := scanner.ScanToken()

				// Check Token Type
				if tok.Type != exp.expectedType {
					t.Errorf("[%s] Token %d: Expected Type %v, got %v (Lexeme: '%s')",
						tt.name, i, exp.expectedType, tok.Type, tok.Lexeme)
				}

				// Check Lexeme
				if tok.Lexeme != exp.expectedLexeme {
					t.Errorf("[%s] Token %d: Expected Lexeme '%s', got '%s'",
						tt.name, i, exp.expectedLexeme, tok.Lexeme)
				}

				// Check Line Number
				if tok.Line != uint(exp.expectedLine) {
					t.Errorf("[%s] Token %d: Expected Line %d, got %d (Lexeme: '%s')",
						tt.name, i, exp.expectedLine, tok.Line, tok.Lexeme)
				}

				// Break early if we hit EOF so the loop stops
				if tok.Type == tokens.END_OF_FILE {
					break
				}
			}
		})
	}
}
