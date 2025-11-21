package pdf

import (
    "bytes"
    "github.com/phpdave11/gofpdf"
)

// Generate формирует PDF с помощью библиотеки gofpdf.
// Возвращает байты PDF-файла.
func Generate(title string, lines []string) []byte {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.SetTitle(title, false)
    pdf.AddPage()

    pdf.SetFont("Helvetica", "", 16)
    pdf.Cell(0, 10, title)
    pdf.Ln(12)

    pdf.SetFont("Helvetica", "", 12)
    for _, line := range lines {
        pdf.CellFormat(0, 8, line, "", 1, "", false, 0, "")
    }

    var buf bytes.Buffer
    _ = pdf.Output(&buf)
    return buf.Bytes()
}