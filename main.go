package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

var outputFileName string

func main() {
	flag.StringVar(&outputFileName, "output", "output.txt", "Output file name")
	flag.Parse()

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	// Write the folder hierarchy tree
	if err := writeFolderTree(currentDir, writer); err != nil {
		fmt.Println("Error writing folder tree:", err)
		return
	}

	writer.WriteString("\n\n")

	// Load .gitignore rules
	gitignoreMatcher, err := loadGitIgnore(currentDir)
	if err != nil {
		fmt.Println("Error loading .gitignore:", err)
		// Not returning here, because we might want to continue without .gitignore rules
	}

	// Process all files using filepath.Walk
	if err := filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(currentDir, path)
		if err != nil {
			return err
		}

		if gitignoreMatcher != nil && gitignoreMatcher.MatchesPath(relPath) {
			fmt.Printf("Ignoring (gitignore): %s\n", relPath)
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldIgnoreFile(info.Name()) || isHiddenFileOrDir(info.Name()) || info.Name() == outputFileName {
			fmt.Printf("Ignoring: %s\n", relPath)
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			return processFile(path, currentDir, writer)
		}
		return nil
	}); err != nil {
		fmt.Println("Error walking directory:", err)
	}

	fmt.Printf("Output written to %s\n", outputFileName)
}

func shouldIgnoreFile(filename string) bool {
	ignoreSuffixes := []string{
		".lock", ".conf", ".config", ".ini", ".log", ".tmp", ".cache", ".mod", ".sum",
	}

	lowerFilename := strings.ToLower(filename)
	for _, suffix := range ignoreSuffixes {
		if strings.HasSuffix(lowerFilename, suffix) {
			return true
		}
	}

	ignoreNames := []string{
		"package-lock.json", "yarn.lock", "composer.lock", "Gemfile.lock",
		"Cargo.lock", "Pipfile.lock", "poetry.lock", "mix.lock", "LICENSE",
		"README", "CHANGELOG", "CONTRIBUTING", "AUTHORS", "CONTRIBUTORS", "CODE_OF_CONDUCT",
	}

	for _, name := range ignoreNames {
		if lowerFilename == name {
			return true
		}
	}

	return strings.Contains(lowerFilename, "conf")
}

func isHiddenFileOrDir(name string) bool {
	return strings.HasPrefix(name, ".")
}

func processFile(filePath, currentDir string, writer *bufio.Writer) error {
	if isLikelyTextFile(filePath) {
		relPath, err := filepath.Rel(currentDir, filePath)
		if err != nil {
			return err
		}
		fullPath := filepath.Join(currentDir, relPath)
		fmt.Fprintf(writer, "File: %s\n", fullPath)

		fileContent, err := os.ReadFile(fullPath)
		if err != nil {
			return err
		}

		cleanedContent := strings.Map(func(r rune) rune {
			if r >= 32 && r <= 126 {
				return r
			}
			return -1
		}, string(fileContent))

		writer.WriteString(cleanedContent)
		writer.WriteString("\n\n")
	}
	return nil
}

func isLikelyTextFile(filename string) bool {
	textFileSuffixes := []string{
		".asm", ".asp", ".aspx", ".awk", ".bat", ".c", ".cfg", ".cfm", ".cgi", ".clj", ".cls", ".cmd",
		".coffee", ".conf", ".cpp", ".cs", ".css", ".csv", ".dart", ".diff", ".dockerfile", ".elm",
		".erl", ".ex", ".exs", ".f", ".f90", ".f95", ".fs", ".go", ".gradle", ".groovy", ".h", ".haml",
		".handlebars", ".hbs", ".hpp", ".hs", ".htm", ".html", ".ini", ".java", ".jl", ".js", ".json",
		".jsp", ".jsx", ".kt", ".kts", ".less", ".log", ".lua", ".m", ".makefile", ".md", ".ml", ".mm",
		".mod", ".nix", ".pas", ".patch", ".php", ".pl", ".pm", ".po", ".properties", ".ps", ".ps1",
		".py", ".r", ".rb", ".rdf", ".rs", ".rst", ".rtf", ".s", ".sass", ".scala", ".sch", ".scss",
		".sh", ".shtml", ".sql", ".svg", ".swift", ".tcl", ".tex", ".tf", ".ts", ".tsx", ".ttl",
		".txt", ".vb", ".vba", ".vbs", ".vhdl", ".vim", ".vue", ".wasm", ".wiki", ".xhtml", ".xml",
		".xsl", ".yaml", ".yml", ".zsh",
	}

	lowerFilename := strings.ToLower(filename)
	for _, suffix := range textFileSuffixes {
		if strings.HasSuffix(lowerFilename, suffix) {
			return true
		}
	}

	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err = file.Read(buffer); err != nil {
		return false
	}

	contentType := http.DetectContentType(buffer)
	return strings.HasPrefix(contentType, "text/")
}

func writeFolderTree(dir string, writer *bufio.Writer) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "tree", "/F", dir)
	} else {
		cmd = exec.Command("tree", "-I", ".*", dir)
		if _, err := exec.LookPath("tree"); err != nil {
			cmd = exec.Command("find", ".", "-not", "-path", "*/.*", "-print")
		}
	}
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	writer.WriteString(string(output))
	return nil
}

func loadGitIgnore(dir string) (*ignore.GitIgnore, error) {
	gitignorePath := filepath.Join(dir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		return nil, nil // No .gitignore file
	}

	return ignore.CompileIgnoreFile(gitignorePath)
}
