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
	flag.StringVar(&outputFileName, "output", "codebase.md", "Output file name")
	flag.Parse()

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Initialize .cpignore file
	if err := initCPIgnore(currentDir); err != nil {
		fmt.Println("Error initializing .cpignore:", err)
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

	// Load .gitignore and .cpignore rules
	ignoreMatcher, _ := loadIgnoreRules(currentDir)
	if ignoreMatcher == nil {
		fmt.Println("Error loading ignore rules")
		return
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

		if ignoreMatcher.MatchesPath(relPath) {
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

func initCPIgnore(dir string) error {
	cpignorePath := filepath.Join(dir, ".cpignore")
	if _, err := os.Stat(cpignorePath); os.IsNotExist(err) {
		file, err := os.Create(cpignorePath)
		if err != nil {
			return err
		}
		defer file.Close()
	} else if err != nil {
		return err
	}
	return nil
}

		defaultRules := []string{
			// Version control
			".git", ".svn", ".hg", ".bzr", "CVS", ".hcl",

			// Dependencies and package managers
			"node_modules", "bower_components", "jspm_packages", "packages",
			"vendor", ".vendor", "*.sum", // Added *.sum here
			".venv", "venv", "env", ".env",
			"__pycache__", ".pytest_cache", ".mypy_cache", ".ruff_cache",
			"pip-wheel-metadata", ".pnpm-store",

			// Build outputs and caches
			"build", "dist", "out", "target", "bin", "obj",
			".gradle", ".maven", ".m2", ".ivy2", ".sbt",
			".next", ".nuxt", ".vuepress", ".docusaurus",
			".cache", ".parcel-cache", ".webpack",

			// IDE and editor files
			".idea", ".vscode", ".vs", "*.swp", "*.swo", "*~",
			".project", ".classpath", ".settings",

			// OS-specific files
			".DS_Store", "Thumbs.db", "desktop.ini",

			// Logs and temporary files
			"*.log", "npm-debug.log*", "yarn-debug.log*", "yarn-error.log*",
			"*.tmp", "*.temp", "*.bak", "*.swp",

			// Configuration files
			"*.lock", "*.conf", "*.config", "*.ini",
			"package-lock.json", "yarn.lock", "composer.lock", "Gemfile.lock",
			"Cargo.lock", "Pipfile.lock", "poetry.lock", "mix.lock",

			// Documentation and metadata
			"README*", "CHANGELOG*", "CONTRIBUTING*", "AUTHORS", "CONTRIBUTORS", "CODE_OF_CONDUCT*",
			"LICENSE", "COPYING",

			// Testing and coverage
			"coverage", ".nyc_output", ".coveralls.yml", ".istanbul.yml",

			// Mobile development
			"Pods", ".cocoapods",

			// Database files
			"*.sqlite", "*.db",

			// Compiled source and binaries
			"*.com", "*.class", "*.dll", "*.exe", "*.o", "*.so",

			// Compressed files
			"*.7z", "*.dmg", "*.gz", "*.iso", "*.jar", "*.rar", "*.tar", "*.zip",

			// Media files (if not part of the project)
			"*.mp3", "*.mp4", "*.avi", "*.mov", "*.wav",

			// Other common ignore patterns
			"*.pid", "*.seed", "*.pid.lock",
			".env.local", ".env.development.local", ".env.test.local", ".env.production.local",
			".eslintcache", ".stylelintcache",
			".terraform", "*.tfstate", "*.tfstate.*",

			// Specific files to ignore
			outputFileName, // The name of the output file
			".cpignore",    // The .cpignore file itself
		}

		for _, rule := range defaultRules {
			if _, err := file.WriteString(rule + "\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

func loadIgnoreRules(dir string) (*ignore.GitIgnore, error) {
	gitignorePath := filepath.Join(dir, ".gitignore")
	cpignorePath := filepath.Join(dir, ".cpignore")

	var lines []string

	// Read .gitignore if it exists
	if gitignoreContent, err := os.ReadFile(gitignorePath); err == nil {
		lines = append(lines, strings.Split(string(gitignoreContent), "\n")...)
	}

	// Read .cpignore if it exists
	if cpignoreContent, err := os.ReadFile(cpignorePath); err == nil {
		lines = append(lines, strings.Split(string(cpignoreContent), "\n")...)
	}

	// Compile ignore rules
	ignoreRules := ignore.CompileIgnoreLines(lines...)
	return ignoreRules, nil
}

func processFile(filePath, currentDir string, writer *bufio.Writer) error {
	if isLikelyTextFile(filePath) {
		relPath, err := filepath.Rel(currentDir, filePath)
		if err != nil {
			return err
		}
		fullPath := filepath.Join(currentDir, relPath)
		fmt.Fprintf(writer, "### %s\n", fullPath)

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

		writer.WriteString("```\n")
		writer.WriteString(cleanedContent)
		writer.WriteString("\n```\n\n")
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
	writer.WriteString("### Project Structure\n")

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
	// Convert tree output to Markdown format
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			writer.WriteString("- " + line + "\n")
		}
	}
	return nil
}
