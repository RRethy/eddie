package search

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearcher_Search_Go(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		query    string
		expected bool
	}{
		{
			name: "basic function search",
			content: `package main

func hello() {
	println("hello")
}`,
			query:    "(function_declaration name: (identifier) @func)",
			expected: true,
		},
		{
			name: "call expression search",
			content: `package main

func main() {
	hello()
}`,
			query:    "(call_expression function: (identifier) @call)",
			expected: true,
		},
		{
			name: "no match",
			content: `package main

var x = 1`,
			query:    "(function_declaration name: (identifier) @func)",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.go")
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			s := &Searcher{}
			err = s.Search(testFile, tt.query)
			assert.NoError(t, err)
		})
	}
}

func TestSearcher_Search_JavaScript(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		content   string
		query     string
		shouldErr bool
	}{
		{
			name:     "function declaration",
			filename: "test.js",
			content: `function hello() {
	console.log("hello");
}`,
			query:     "(function_declaration name: (identifier) @func)",
			shouldErr: false,
		},
		{
			name:     "arrow function",
			filename: "test.js",
			content: `const greet = () => {
	console.log("hello");
};`,
			query:     "(arrow_function) @arrow",
			shouldErr: false,
		},
		{
			name:     "jsx component",
			filename: "test.jsx",
			content: `function Component() {
	return <div>Hello</div>;
}`,
			query:     "(function_declaration name: (identifier) @comp)",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			s := &Searcher{}
			err = s.Search(testFile, tt.query)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearcher_Search_TypeScript(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		content   string
		query     string
		shouldErr bool
	}{
		{
			name:     "function with types",
			filename: "test.ts",
			content: `function add(a: number, b: number): number {
	return a + b;
}`,
			query:     "(function_declaration name: (identifier) @func)",
			shouldErr: false,
		},
		{
			name:     "interface declaration",
			filename: "test.ts",
			content: `interface User {
	name: string;
	age: number;
}`,
			query:     "(interface_declaration name: (type_identifier) @interface)",
			shouldErr: false,
		},
		{
			name:     "tsx component",
			filename: "test.tsx",
			content: `interface Props {
	name: string;
}

function Component({ name }: Props) {
	return <div>{name}</div>;
}`,
			query:     "(function_declaration name: (identifier) @comp)",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			s := &Searcher{}
			err = s.Search(testFile, tt.query)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearcher_Search_Python(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		content   string
		query     string
		shouldErr bool
	}{
		{
			name:     "function definition",
			filename: "test.py",
			content: `def hello():
    print("hello")`,
			query:     "(function_definition name: (identifier) @func)",
			shouldErr: false,
		},
		{
			name:     "class definition",
			filename: "test.py",
			content: `class Person:
    def __init__(self, name):
        self.name = name`,
			query:     "(class_definition name: (identifier) @class)",
			shouldErr: false,
		},
		{
			name:     "python interface file",
			filename: "test.pyi",
			content: `def greet(name: str) -> str: ...`,
			query:     "(function_definition name: (identifier) @func)",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			s := &Searcher{}
			err = s.Search(testFile, tt.query)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearcher_Search_Rust(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		query     string
		shouldErr bool
	}{
		{
			name: "function definition",
			content: `fn main() {
    println!("Hello, world!");
}`,
			query:     "(function_item name: (identifier) @func)",
			shouldErr: false,
		},
		{
			name: "struct definition",
			content: `struct Person {
    name: String,
    age: u32,
}`,
			query:     "(struct_item name: (type_identifier) @struct)",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.rs")
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			s := &Searcher{}
			err = s.Search(testFile, tt.query)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearcher_Search_Java(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		query     string
		shouldErr bool
	}{
		{
			name: "class definition",
			content: `public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}`,
			query:     "(class_declaration name: (identifier) @class)",
			shouldErr: false,
		},
		{
			name: "method definition",
			content: `public class Test {
    public void greet() {
        System.out.println("Hello");
    }
}`,
			query:     "(method_declaration name: (identifier) @method)",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.java")
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			s := &Searcher{}
			err = s.Search(testFile, tt.query)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearcher_Search_C(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		content   string
		query     string
		shouldErr bool
	}{
		{
			name:     "function definition",
			filename: "test.c",
			content: `#include <stdio.h>

int main() {
    printf("Hello, World!\n");
    return 0;
}`,
			query:     "(function_definition declarator: (function_declarator declarator: (identifier) @func))",
			shouldErr: false,
		},
		{
			name:     "header file",
			filename: "test.h",
			content: `#ifndef TEST_H
#define TEST_H

void greet(void);

#endif`,
			query:     "(declaration declarator: (function_declarator declarator: (identifier) @func))",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			s := &Searcher{}
			err = s.Search(testFile, tt.query)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearcher_Search_CPP(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		content   string
		query     string
		shouldErr bool
	}{
		{
			name:     "cpp function",
			filename: "test.cpp",
			content: `#include <iostream>

int main() {
    std::cout << "Hello, World!" << std::endl;
    return 0;
}`,
			query:     "(function_definition declarator: (function_declarator declarator: (identifier) @func))",
			shouldErr: false,
		},
		{
			name:     "class definition",
			filename: "test.cpp",
			content: `class Person {
private:
    std::string name;
public:
    Person(const std::string& n) : name(n) {}
    void greet() {
        std::cout << "Hello, " << name << std::endl;
    }
};`,
			query:     "(class_specifier name: (type_identifier) @class)",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			s := &Searcher{}
			err = s.Search(testFile, tt.query)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearcher_Search_UnsupportedFiles(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  string
		query    string
	}{
		{
			name:     "text file ignored",
			filename: "test.txt",
			content:  "This is just plain text",
			query:    "(function_declaration name: (identifier) @func)",
		},
		{
			name:     "markdown file ignored",
			filename: "test.md",
			content:  "# Header\n\nThis is markdown",
			query:    "(function_declaration name: (identifier) @func)",
		},
		{
			name:     "unknown extension ignored",
			filename: "test.unknown",
			content:  "unknown content",
			query:    "(function_declaration name: (identifier) @func)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			s := &Searcher{}
			err = s.Search(testFile, tt.query)
			assert.NoError(t, err)
		})
	}
}

func TestSearcher_SearchDir_MultiLanguage(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(goFile, []byte(`package main
func hello() {
	println("hello")
}`), 0644)
	require.NoError(t, err)

	jsFile := filepath.Join(tmpDir, "test.js")
	err = os.WriteFile(jsFile, []byte(`function greet() {
	console.log("hello");
}`), 0644)
	require.NoError(t, err)

	txtFile := filepath.Join(tmpDir, "readme.txt")
	err = os.WriteFile(txtFile, []byte("This should be ignored"), 0644)
	require.NoError(t, err)

	s := &Searcher{}
	
	t.Run("function declaration search (Go and JS)", func(t *testing.T) {
		err = s.Search(tmpDir, "(function_declaration name: (identifier) @func)")
		assert.NoError(t, err)
	})
}

func TestSearcher_SearchDir_LanguageSpecific(t *testing.T) {
	tmpDir := t.TempDir()

	pyFile := filepath.Join(tmpDir, "test.py")
	err := os.WriteFile(pyFile, []byte(`def welcome():
    print("hello")`), 0644)
	require.NoError(t, err)

	s := &Searcher{}
	
	t.Run("python function search", func(t *testing.T) {
		err = s.Search(tmpDir, "(function_definition name: (identifier) @func)")
		assert.NoError(t, err)
	})
}

func TestSearcher_SearchDir_NestedDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	require.NoError(t, err)

	goFile := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(goFile, []byte(`package main
func main() {
	println("main")
}`), 0644)
	require.NoError(t, err)

	nestedGoFile := filepath.Join(subDir, "helper.go")
	err = os.WriteFile(nestedGoFile, []byte(`package main
func helper() {
	println("helper")
}`), 0644)
	require.NoError(t, err)

	jsFile := filepath.Join(subDir, "script.js")
	err = os.WriteFile(jsFile, []byte(`function process() {
	console.log("processing");
}`), 0644)
	require.NoError(t, err)

	s := &Searcher{}
	err = s.Search(tmpDir, "(function_declaration name: (identifier) @func)")
	assert.NoError(t, err)
}

func TestGetLanguageFromFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		hasLang  bool
	}{
		{"go file", "test.go", true},
		{"javascript file", "test.js", true},
		{"jsx file", "test.jsx", true},
		{"typescript file", "test.ts", true},
		{"tsx file", "test.tsx", true},
		{"python file", "test.py", true},
		{"python interface", "test.pyi", true},
		{"rust file", "test.rs", true},
		{"java file", "test.java", true},
		{"c file", "test.c", true},
		{"header file", "test.h", true},
		{"cpp file", "test.cpp", true},
		{"hpp file", "test.hpp", true},
		{"text file", "test.txt", false},
		{"markdown file", "test.md", false},
		{"no extension", "testfile", false},
		{"case insensitive", "Test.GO", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang := getLanguageFromFile(tt.filename)
			if tt.hasLang {
				assert.NotNil(t, lang, "expected language for %s", tt.filename)
			} else {
				assert.Nil(t, lang, "expected no language for %s", tt.filename)
			}
		})
	}
}
