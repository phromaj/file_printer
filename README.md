# Code Printer

Code Printer is a command-line tool that helps you copy and paste code from your project files into AI chat conversations. It generates a text file containing the folder hierarchy and the contents of text-based files in your project directory.

## Features

- Generates a folder hierarchy tree of your project directory
- Includes the contents of text-based files in the output
- Ignores files and directories based on predefined rules and .gitignore
- Supports cross-platform builds for M1 Mac, x64 Windows, and x64 Linux

## Prerequisites

- Go programming language (version 1.16 or later)

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/your-username/code-printer.git
   ```

2. Change to the project directory:

   ```
   cd code-printer
   ```

3. Initialize the Go module and download dependencies:

   ```
   go mod init github.com/your-username/code-printer
   go mod tidy
   ```

   Make sure to replace `github.com/your-username/code-printer` with your actual repository path.

4. Review the `main.go` file to understand how the project is structured and how to customize it if needed.

5. Run the build script:

   ```
   ./build.sh
   ```

   This will generate the executable files for different platforms in the `bin` directory.

## Usage

1. Navigate to your project directory:

   ```
   cd /path/to/your/project
   ```

2. Run the Code Printer tool:

   ```
   /path/to/code-printer/bin/codeprinter
   ```

   By default, the output will be written to a file named `output.txt` in the current directory.

   You can specify a different output file name using the `-output` flag:

   ```
   /path/to/code-printer/bin/codeprinter -output myoutput.txt
   ```

3. The generated output file will contain the folder hierarchy tree and the contents of text-based files in your project directory.

## Configuration

Code Printer automatically ignores certain files and directories based on predefined rules and the `.gitignore` file in your project directory.

- Predefined ignored file suffixes: `.lock`, `.conf`, `.config`, `.ini`, `.log`, `.tmp`, `.cache`, `.mod`, `.sum`
- Predefined ignored file names: `package-lock.json`, `yarn.lock`, `composer.lock`, `Gemfile.lock`, `Cargo.lock`, `Pipfile.lock`, `poetry.lock`, `mix.lock`
- Hidden files and directories (starting with `.`) are also ignored
- Files and directories specified in the `.gitignore` file are ignored

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is open-source and available under the [MIT License](LICENSE).