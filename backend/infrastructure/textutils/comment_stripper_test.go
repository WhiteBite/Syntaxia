package textutils

import (
	"strings"
	"testing"
)

// Uses mockLogger from context_splitter_test.go

func TestNewCommentStripper(t *testing.T) {
	cs := NewCommentStripper(&mockLogger{})
	if cs == nil {
		t.Fatal("NewCommentStripper returned nil")
	}
}

func TestCommentStripper_Strip_Go(t *testing.T) {
	cs := NewCommentStripper(&mockLogger{})

	content := `package main
// line comment
func main() {
	/* block comment */
	fmt.Println("hello")
}
`
	result := cs.Strip(content, "main.go")

	if strings.Contains(result, "line comment") {
		t.Error("should strip line comments")
	}
	if strings.Contains(result, "block comment") {
		t.Error("should strip block comments")
	}
	if !strings.Contains(result, "package main") {
		t.Error("should keep code")
	}
	if !strings.Contains(result, "fmt.Println") {
		t.Error("should keep code")
	}
}

func TestCommentStripper_Strip_Python(t *testing.T) {
	cs := NewCommentStripper(&mockLogger{})

	content := `def main():
    # this is a comment
    print("hello")
`
	result := cs.Strip(content, "main.py")

	if strings.Contains(result, "this is a comment") {
		t.Error("should strip hash comments")
	}
	if !strings.Contains(result, "def main") {
		t.Error("should keep code")
	}
}

func TestCommentStripper_Strip_HTML(t *testing.T) {
	cs := NewCommentStripper(&mockLogger{})

	content := `<html>
<!-- this is a comment -->
<body>Hello</body>
</html>
`
	result := cs.Strip(content, "index.html")

	if strings.Contains(result, "this is a comment") {
		t.Error("should strip XML comments")
	}
	if !strings.Contains(result, "<body>") {
		t.Error("should keep HTML")
	}
}

func TestCommentStripper_Strip_Unknown(t *testing.T) {
	cs := NewCommentStripper(&mockLogger{})

	content := "some content // with comment"
	result := cs.Strip(content, "file.unknown")

	// Unknown extension should return content unchanged
	if result != content {
		t.Error("unknown extension should return content unchanged")
	}
}

func TestStripCStyleComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
		excludes string
	}{
		{
			name:     "line comment",
			input:    "code // comment\nmore",
			contains: "code",
			excludes: "comment",
		},
		{
			name:     "block comment single line",
			input:    "code /* comment */ more",
			contains: "code",
			excludes: "comment",
		},
		{
			name:     "multiline block comment",
			input:    "code\n/* multi\nline */\nmore",
			contains: "code",
			excludes: "multi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripCStyleComments(tt.input)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("result should contain %q", tt.contains)
			}
			if strings.Contains(result, tt.excludes) {
				t.Errorf("result should not contain %q", tt.excludes)
			}
		})
	}
}

func TestStripHashComments(t *testing.T) {
	input := "code # comment\nmore code"
	result := stripHashComments(input)

	if strings.Contains(result, "comment") {
		t.Error("should strip hash comments")
	}
	if !strings.Contains(result, "code") {
		t.Error("should keep code")
	}
}

func TestStripXMLComments(t *testing.T) {
	input := "<tag><!-- comment --></tag>"
	result := stripXMLComments(input)

	if strings.Contains(result, "comment") {
		t.Error("should strip XML comments")
	}
	if !strings.Contains(result, "<tag>") {
		t.Error("should keep tags")
	}
}

func TestStripCStyleComments_EmptyLines(t *testing.T) {
	input := "code\n// comment\n\nmore"
	result := stripCStyleComments(input)

	// Should not have empty lines from stripped comments
	if strings.Contains(result, "\n\n\n") {
		t.Error("should not have multiple empty lines")
	}
}
