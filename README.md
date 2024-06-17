# Code Printer

Code Printer is a utility designed to help developers quickly and efficiently copy and paste code in AI chat environments. This tool traverses directories, processes files, and outputs their content in a clean and readable format, making it easier to share code snippets during conversations.

## Features

- Traverses directories and processes files, respecting `.gitignore` rules.
- Ignores files based on configurable suffixes and names.
- Builds for multiple platforms (macOS M1, Windows x64, Linux x64).
- Outputs the directory structure and file content in a readable format.

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) installed on your machine.
- A Unix-like shell (e.g., bash) for running the build script.

### Installing Dependencies

The project uses Go modules for dependency management, and the dependencies are specified in the `go.mod` and `go.sum` files.

To install the dependencies, run the following command:

```sh
go mod tidy
```

This will download and install the necessary dependencies as specified in the `go.mod` file.

### Building the Project

The `build.sh` script is provided to build the project for different platforms. Follow these steps to build the project:

1. Ensure the `build.sh` script has execute permissions. If not, you can set it using:

    ```sh
    chmod +x build.sh
    ```

2. Run the build script to generate executables for macOS M1, Windows x64, and Linux x64:

    ```sh
    ./build.sh
    ```

### Running the Program

To run the Code Printer, use the following command:

```sh
./bin/codeprinter-mac-arm64 -output output.txt
```

Replace `codeprinter-mac-arm64` with the appropriate binary for your platform.

### Command-line Options

- `-output`: Specify the name of the output file. Default is `output.txt`.

### Example

```sh
./bin/codeprinter-mac-arm64 -output output.txt
```

The program will traverse the current directory, process the files, and write the output to `output.txt`.

## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a pull request.

## License

This project is licensed under the MIT License.

## Acknowledgements

- [go-gitignore](https://github.com/sabhiram/go-gitignore) for handling `.gitignore` rules.
