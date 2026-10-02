package language

import "fmt"

type Language struct {
	ID             string
	Image          string
	SourceFilename string
	CompileCommand []string
	RunCommand     []string
}

var supported = map[string]Language{
	"cpp": {
		ID:             "cpp",
		Image:          "code-runner-cpp:2026.09",
		SourceFilename: "main.cpp",
		CompileCommand: []string{"g++", "main.cpp", "-std=c++23", "-O2", "-pipe", "-DONLINE_JUDGE", "-o", "main"},
		RunCommand:     []string{"./main"},
	},
	"python": {
		ID:             "python",
		Image:          "code-runner-pypy:2026.09",
		SourceFilename: "main.py",
		RunCommand:     []string{"pypy3", "main.py"},
	},
}

func Lookup(id string) (Language, error) {
	language, ok := supported[id]
	if !ok {
		return Language{}, fmt.Errorf("unsupported language %q", id)
	}
	return language, nil
}
