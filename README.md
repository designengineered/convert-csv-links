# CSV to Links Converter

A simple desktop application that converts CSV files containing text and URLs into HTML and/or TXT link formats.

## Features

- Convert CSV files to HTML and/or TXT formats
- Choose output format(s) (HTML, TXT, or both)
- Simple and intuitive GUI
- File overwrite protection
- Cross-platform support

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [Releases](https://github.com/yourusername/csv-to-links/releases) page:

- Windows: `csv-to-links-windows.exe`
- macOS: `csv-to-links-macos.app`
- Linux: `csv-to-links-linux`

### Building from Source

1. Install Go 1.21 or later
2. Install required dependencies:
   ```bash
   go get fyne.io/fyne/v2
   ```
3. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/csv-to-links.git
   cd csv-to-links
   ```
4. Build the application:
   ```bash
   go build
   ```

## Usage

1. Launch the application
2. Select your desired output format(s) using the checkboxes
3. Click "Convert CSV to Links" button
4. Select your input CSV file
5. Choose where to save the output file(s)

### CSV Format

Your CSV file should have at least two columns:
- First column: Link text
- Second column: URL

Example:
```csv
Text,URL
Google,https://www.google.com
GitHub,https://github.com
```

## Development

### Prerequisites

- Go 1.21 or later
- Fyne toolkit dependencies:
  - Windows: gcc
  - macOS: Xcode
  - Linux: gcc, libgl1-mesa-dev, xorg-dev

### Building for Different Platforms

#### Windows
```bash
GOOS=windows GOARCH=amd64 go build -o csv-to-links-windows.exe
```

#### macOS
```bash
GOOS=darwin GOARCH=amd64 go build -o csv-to-links-macos
```

#### Linux
```bash
GOOS=linux GOARCH=amd64 go build -o csv-to-links-linux
```

## License

MIT License - See [LICENSE](LICENSE) for details
