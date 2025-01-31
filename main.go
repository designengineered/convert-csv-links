package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// OutputFormat defines which output formats should be generated
type OutputFormat struct {
	HTML bool
	TXT  bool
}

// generateLinks creates the output files based on the CSV input
func generateLinks(csvPath string, outputPath string, format OutputFormat) error {
	// Open CSV file
	csvFile, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %w", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	// Skip header row
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("error reading header: %w", err)
	}

	// Create output files based on selected formats
	var htmlFile, txtFile *os.File

	if format.HTML {
		htmlPath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".html"
		htmlFile, err = os.Create(htmlPath)
		if err != nil {
			return fmt.Errorf("error creating HTML file: %w", err)
		}
		defer htmlFile.Close()

		// Write HTML header
		htmlFile.WriteString("<!DOCTYPE html>\n<html>\n<body>\n")
	}

	if format.TXT {
		txtPath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".txt"
		txtFile, err = os.Create(txtPath)
		if err != nil {
			return fmt.Errorf("error creating TXT file: %w", err)
		}
		defer txtFile.Close()
	}

	// Process each row
	for {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("error reading row: %w", err)
		}

		if len(row) >= 2 {
			text := strings.TrimSpace(row[0])
			url := strings.TrimSpace(row[1])
			link := fmt.Sprintf("<a href=\"%s\">%s</a>\n", url, text)

			if format.HTML {
				htmlFile.WriteString(link)
			}
			if format.TXT {
				txtFile.WriteString(link)
			}
		} else {
			fmt.Printf("Skipping invalid row: %v\n", row)
		}
	}

	// Close HTML file with footer
	if format.HTML {
		htmlFile.WriteString("</body>\n</html>")
	}

	return nil
}

// proceedWithConversion handles the actual conversion process
func proceedWithConversion(inputPath, outputPath string, format OutputFormat, window fyne.Window, status *widget.Label) {
	err := generateLinks(inputPath, outputPath, format)
	if err != nil {
		dialog.ShowError(err, window)
		status.SetText("")
		return
	}

	// Show success message
	extensions := []string{}
	if format.HTML {
		extensions = append(extensions, ".html")
	}
	if format.TXT {
		extensions = append(extensions, ".txt")
	}
	successMsg := fmt.Sprintf("Files generated successfully!\nSaved as: %s",
		strings.TrimSuffix(outputPath, filepath.Ext(outputPath))+
			strings.Join(extensions, " and "))

	dialog.ShowInformation("Success", successMsg, window)
	status.SetText("Conversion completed successfully!")
}

func main() {
	myApp := app.NewWithID("com.csvconverter.app")
	window := myApp.NewWindow("CSV to Links Converter")

	status := widget.NewLabel("")

	// Create format selection checkboxes
	format := OutputFormat{false, true} // Default to TXT only
	htmlCheck := widget.NewCheck("HTML", func(checked bool) {
		format.HTML = checked
	})
	txtCheck := widget.NewCheck("TXT", func(checked bool) {
		format.TXT = checked
	})
	txtCheck.SetChecked(true) // Default to TXT selected

	// Group checkboxes
	formatGroup := container.NewHBox(
		widget.NewLabel("Output Format:"),
		htmlCheck,
		txtCheck,
	)

	startButton := widget.NewButton("Convert CSV to Links", func() {
		if !format.HTML && !format.TXT {
			dialog.ShowError(fmt.Errorf("please select at least one output format"), window)
			return
		}

		status.SetText("Opening file selector...")

		openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				status.SetText("")
				return
			}
			if reader == nil {
				status.SetText("")
				return
			}

			inputPath := reader.URI().Path()
			reader.Close()

			status.SetText("Selecting output location...")

			saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, window)
					status.SetText("")
					return
				}
				if writer == nil {
					status.SetText("")
					return
				}

				// Get the path but don't use the writer
				outputPath := writer.URI().Path()
				writer.Close()

				// Remove the file that was automatically created
				os.Remove(outputPath)

				basePath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath))

				// Check for existing files
				var existingFiles []string
				if format.HTML {
					if _, err := os.Stat(basePath + ".html"); err == nil {
						existingFiles = append(existingFiles, basePath+".html")
					}
				}
				if format.TXT {
					if _, err := os.Stat(basePath + ".txt"); err == nil {
						existingFiles = append(existingFiles, basePath+".txt")
					}
				}

				if len(existingFiles) > 0 {
					message := fmt.Sprintf("The following files already exist:\n%s\n\nDo you want to overwrite them?",
						strings.Join(existingFiles, "\n"))
					dialog.ShowConfirm("Files Exist", message, func(overwrite bool) {
						if overwrite {
							status.SetText("Converting...")
							proceedWithConversion(inputPath, outputPath, format, window, status)
						} else {
							status.SetText("")
						}
					}, window)
				} else {
					status.SetText("Converting...")
					proceedWithConversion(inputPath, outputPath, format, window, status)
				}
			}, window)

			initialName := filepath.Base(inputPath)
			initialName = strings.TrimSuffix(initialName, filepath.Ext(initialName))
			saveDialog.SetFileName(initialName)
			saveDialog.Show()
		}, window)

		csvFilter := storage.NewExtensionFileFilter([]string{".csv"})
		openDialog.SetFilter(csvFilter)
		openDialog.Show()
	})

	instructions := widget.NewTextGridFromString(
		"This tool converts CSV files to link format.\n\n" +
			"1. Select your desired output format(s) using the checkboxes\n" +
			"2. Click the convert button to start\n" +
			"3. Select your input CSV file\n" +
			"4. Choose where to save the output file(s)")

	content := container.NewVBox(
		instructions,
		formatGroup,
		startButton,
		status,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(800, 600))
	window.CenterOnScreen()
	window.Show()
	myApp.Run()
}
