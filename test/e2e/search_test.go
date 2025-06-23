package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchCommand_Go(t *testing.T) {
	tmpDir := t.TempDir()

	goFile := filepath.Join(tmpDir, "main.go")
	goContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	greet("Eddie")
}

func greet(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

type Person struct {
	Name string
	Age  int
}

func (p Person) sayHello() {
	fmt.Printf("Hi, I'm %s\n", p.Name)
}`
	require.NoError(t, os.WriteFile(goFile, []byte(goContent), 0644))

	tests := []struct {
		name         string
		args         []string
		wantContains []string
		wantErr      bool
	}{
		{
			name: "function declaration search",
			args: []string{"search", tmpDir, "--tree-sitter-query", "(function_declaration name: (identifier) @func)"},
			wantContains: []string{
				"@func: func main()",
				"@func: func greet(name string)",
			},
			wantErr: false,
		},

		{
			name: "call expression search",
			args: []string{"search", tmpDir, "--tree-sitter-query", "(call_expression function: (identifier) @call)"},
			wantContains: []string{
				"@call: greet(\"Eddie\")",
			},
			wantErr: false,
		},
		{
			name: "struct type search",
			args: []string{"search", tmpDir, "--tree-sitter-query", "(type_declaration (type_spec name: (type_identifier) @struct type: (struct_type)))"},
			wantContains: []string{
				"@struct: type Person struct",
			},
			wantErr: false,
		},
		{
			name: "single file search",
			args: []string{"search", goFile, "--tree-sitter-query", "(function_declaration name: (identifier) @func)"},
			wantContains: []string{
				"@func: func main()",
				"@func: func greet(name string)",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runEddie(t, tt.args...)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, "stderr: %s", stderr)
				for _, want := range tt.wantContains {
					assert.Contains(t, stdout, want, "stdout should contain: %s", want)
				}
			}
		})
	}
}

func TestSearchCommand_JavaScript(t *testing.T) {
	tmpDir := t.TempDir()

	jsFile := filepath.Join(tmpDir, "script.js")
	jsContent := `function greet(name) {
    console.log("Hello, " + name + "!");
}

const add = (a, b) => {
    return a + b;
};

class Person {
    constructor(name) {
        this.name = name;
    }
    
    sayHello() {
        greet(this.name);
    }
}

greet("World");`
	require.NoError(t, os.WriteFile(jsFile, []byte(jsContent), 0644))

	tests := []struct {
		name         string
		args         []string
		wantContains []string
		wantErr      bool
	}{
		{
			name: "function declaration search",
			args: []string{"search", tmpDir, "--tree-sitter-query", "(function_declaration name: (identifier) @func)"},
			wantContains: []string{
				"@func: function greet(name)",
			},
			wantErr: false,
		},
		{
			name: "arrow function search",
			args: []string{"search", tmpDir, "--tree-sitter-query", "(arrow_function) @arrow"},
			wantContains: []string{
				"@arrow: const add = (a, b) =>",
			},
			wantErr: false,
		},
		{
			name: "class declaration search",
			args: []string{"search", tmpDir, "--tree-sitter-query", "(class_declaration name: (identifier) @class)"},
			wantContains: []string{
				"@class: class Person",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runEddie(t, tt.args...)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, "stderr: %s", stderr)
				for _, want := range tt.wantContains {
					assert.Contains(t, stdout, want, "stdout should contain: %s", want)
				}
			}
		})
	}
}

func TestSearchCommand_Python(t *testing.T) {
	tmpDir := t.TempDir()

	pyFile := filepath.Join(tmpDir, "script.py")
	pyContent := `def greet(name):
    print(f"Hello, {name}!")

class Person:
    def __init__(self, name):
        self.name = name
    
    def say_hello(self):
        greet(self.name)

def main():
    person = Person("Alice")
    person.say_hello()

if __name__ == "__main__":
    main()`
	require.NoError(t, os.WriteFile(pyFile, []byte(pyContent), 0644))

	tests := []struct {
		name         string
		args         []string
		wantContains []string
		wantErr      bool
	}{
		{
			name: "function definition search",
			args: []string{"search", tmpDir, "--tree-sitter-query", "(function_definition name: (identifier) @func)"},
			wantContains: []string{
				"@func: def greet(name):",
				"@func: def __init__(self, name):",
				"@func: def say_hello(self):",
				"@func: def main():",
			},
			wantErr: false,
		},
		{
			name: "class definition search",
			args: []string{"search", tmpDir, "--tree-sitter-query", "(class_definition name: (identifier) @class)"},
			wantContains: []string{
				"@class: class Person:",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runEddie(t, tt.args...)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, "stderr: %s", stderr)
				for _, want := range tt.wantContains {
					assert.Contains(t, stdout, want, "stdout should contain: %s", want)
				}
			}
		})
	}
}

func TestSearchCommand_MultiLanguage(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple language files
	goFile := filepath.Join(tmpDir, "main.go")
	goContent := `package main
func hello() {
	println("Go Hello")
}`
	require.NoError(t, os.WriteFile(goFile, []byte(goContent), 0644))

	jsFile := filepath.Join(tmpDir, "script.js")
	jsContent := `function hello() {
    console.log("JS Hello");
}`
	require.NoError(t, os.WriteFile(jsFile, []byte(jsContent), 0644))

	pyFile := filepath.Join(tmpDir, "script.py")
	pyContent := `def hello():
    print("Python Hello")`
	require.NoError(t, os.WriteFile(pyFile, []byte(pyContent), 0644))

	// Create unsupported file that should be ignored
	txtFile := filepath.Join(tmpDir, "readme.txt")
	require.NoError(t, os.WriteFile(txtFile, []byte("This is plain text"), 0644))

	tests := []struct {
		name         string
		args         []string
		wantContains []string
		wantErr      bool
		description  string
	}{
		{
			name: "function declaration search on specific files",
			args: []string{"search", goFile, "--tree-sitter-query", "(function_declaration name: (identifier) @func)"},
			wantContains: []string{
				"main.go:",
				"@func:",
			},
			wantErr:     false,
			description: "Should find functions in Go file",
		},
		{
			name: "function definition search on Python file",
			args: []string{"search", pyFile, "--tree-sitter-query", "(function_definition name: (identifier) @func)"},
			wantContains: []string{
				"script.py:",
				"@func:",
			},
			wantErr:     false,
			description: "Should find function in Python file only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runEddie(t, tt.args...)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, "stderr: %s", stderr)
				for _, want := range tt.wantContains {
					assert.Contains(t, stdout, want, "stdout should contain: %s", want)
				}
				// Verify txt file is not mentioned
				assert.NotContains(t, stdout, "readme.txt", "should not search unsupported files")
			}
		})
	}
}

func TestSearchCommand_NestedDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create nested directory structure
	subDir := filepath.Join(tmpDir, "src")
	require.NoError(t, os.Mkdir(subDir, 0755))
	
	nestedDir := filepath.Join(subDir, "utils")
	require.NoError(t, os.Mkdir(nestedDir, 0755))

	// Create files in different directories
	mainFile := filepath.Join(tmpDir, "main.go")
	mainContent := `package main
func main() {
	println("main")
}`
	require.NoError(t, os.WriteFile(mainFile, []byte(mainContent), 0644))

	srcFile := filepath.Join(subDir, "helper.go")
	srcContent := `package src
func helper() {
	println("helper")
}`
	require.NoError(t, os.WriteFile(srcFile, []byte(srcContent), 0644))

	utilFile := filepath.Join(nestedDir, "util.go")
	utilContent := `package utils
func utility() {
	println("utility")
}`
	require.NoError(t, os.WriteFile(utilFile, []byte(utilContent), 0644))

	stdout, stderr, err := runEddie(t, "search", tmpDir, "--tree-sitter-query", "(function_declaration name: (identifier) @func)")
	
	assert.NoError(t, err, "stderr: %s", stderr)
	assert.Contains(t, stdout, "main.go:", "should find main.go")
	assert.Contains(t, stdout, "helper.go:", "should find helper.go")
	assert.Contains(t, stdout, "util.go:", "should find util.go")
	assert.Contains(t, stdout, "@func: func main()", "should find main function")
	assert.Contains(t, stdout, "@func: func helper()", "should find helper function")
	assert.Contains(t, stdout, "@func: func utility()", "should find utility function")
}

func TestSearchCommand_ErrorConditions(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create a test Go file for invalid query test
	testFile := filepath.Join(tmpDir, "test.go")
	require.NoError(t, os.WriteFile(testFile, []byte("package main\nfunc test() {}"), 0644))

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "missing query flag",
			args:    []string{"search", tmpDir},
			wantErr: true,
			errMsg:  "required flag",
		},
		{
			name:    "nonexistent directory",
			args:    []string{"search", "/nonexistent/path", "--tree-sitter-query", "(function_declaration name: (identifier) @func)"},
			wantErr: true,
			errMsg:  "no such file or directory",
		},
		{
			name:    "invalid query syntax on go file",
			args:    []string{"search", testFile, "--tree-sitter-query", "invalid(query"},
			wantErr: true,
			errMsg:  "invalid query",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runEddie(t, tt.args...)
			
			if tt.wantErr {
				assert.Error(t, err)
				output := stdout + stderr
				assert.Contains(t, strings.ToLower(output), strings.ToLower(tt.errMsg), 
					"error output should contain: %s, got: %s", tt.errMsg, output)
			} else {
				assert.NoError(t, err, "stderr: %s", stderr)
			}
		})
	}
}

func TestSearchCommand_TypeScript(t *testing.T) {
	tmpDir := t.TempDir()

	tsFile := filepath.Join(tmpDir, "script.ts")
	tsContent := `interface User {
    name: string;
    age: number;
}

function greet(user: User): string {
    return "Hello, " + user.name + "!";
}

class Person implements User {
    name: string;
    age: number;
    
    constructor(name: string, age: number) {
        this.name = name;
        this.age = age;
    }
}`
	require.NoError(t, os.WriteFile(tsFile, []byte(tsContent), 0644))

	tests := []struct {
		name         string
		args         []string
		wantContains []string
		wantErr      bool
	}{
		{
			name: "interface declaration search",
			args: []string{"search", tmpDir, "--tree-sitter-query", "(interface_declaration name: (type_identifier) @interface)"},
			wantContains: []string{
				"@interface: interface User",
			},
			wantErr: false,
		},
		{
			name: "function declaration search",
			args: []string{"search", tmpDir, "--tree-sitter-query", "(function_declaration name: (identifier) @func)"},
			wantContains: []string{
				"@func: function greet(user: User): string",
			},
			wantErr: false,
		},
		{
			name: "class declaration search",
			args: []string{"search", tmpDir, "--tree-sitter-query", "(class_declaration name: (type_identifier) @class)"},
			wantContains: []string{
				"@class: class Person implements User",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runEddie(t, tt.args...)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, "stderr: %s", stderr)
				for _, want := range tt.wantContains {
					assert.Contains(t, stdout, want, "stdout should contain: %s", want)
				}
			}
		})
	}
}

func TestSearchCommand_Rust(t *testing.T) {
	tmpDir := t.TempDir()

	rsFile := filepath.Join(tmpDir, "main.rs")
	rsContent := `fn main() {
    println!("Hello, world!");
    greet("Rust");
}

fn greet(name: &str) {
    println!("Hello, {}!", name);
}

struct Person {
    name: String,
    age: u32,
}

impl Person {
    fn new(name: String, age: u32) -> Self {
        Person { name, age }
    }
    
    fn greet(&self) {
        println!("Hi, I'm {}", self.name);
    }
}`
	require.NoError(t, os.WriteFile(rsFile, []byte(rsContent), 0644))

	tests := []struct {
		name         string
		args         []string
		wantContains []string
		wantErr      bool
	}{
		{
			name: "function item search",
			args: []string{"search", tmpDir, "--tree-sitter-query", "(function_item name: (identifier) @func)"},
			wantContains: []string{
				"@func: fn main()",
				"@func: fn greet(name: &str)",
				"@func: fn new(name: String, age: u32) -> Self",
				"@func: fn greet(&self)",
			},
			wantErr: false,
		},
		{
			name: "struct item search",
			args: []string{"search", tmpDir, "--tree-sitter-query", "(struct_item name: (type_identifier) @struct)"},
			wantContains: []string{
				"@struct: struct Person",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runEddie(t, tt.args...)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, "stderr: %s", stderr)
				for _, want := range tt.wantContains {
					assert.Contains(t, stdout, want, "stdout should contain: %s", want)
				}
			}
		})
	}
}