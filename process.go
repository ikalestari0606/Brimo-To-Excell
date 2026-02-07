package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

func Process() {
	log.Println("🚀 Proses PDF (1 PDF = 1 Excel) dimulai...")

	err := filepath.Walk("input", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(info.Name()), ".pdf") {
			return nil
		}

		log.Println("📄 Membaca:", path)

		lines := ExtractPDFText(path)
		if len(lines) == 0 {
			log.Println("⚠️ PDF kosong:", path)
			return nil
		}

		data := ParseLines(lines)
		sum := ParseSummary(lines)

		if len(data) == 0 {
			log.Println("⚠️ Tidak ada transaksi:", path)
			return nil
		}

		// 🔹 nama output = nama pdf
		base := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
		outPath := filepath.Join("output", base+".xlsx")

		WriteExcelPerFile(data, sum, outPath)

		log.Println("✅ Selesai:", outPath)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("🎉 Semua PDF selesai diproses")
}
func WriteExcelPerFile(data []Mutasi, sum Summary, outputPath string) {
	f := excelize.NewFile()
	sheet := "Mutasi"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// STYLE
	styleNumber, _ := f.NewStyle(&excelize.Style{NumFmt: 4})
	wrapStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
	})

	// HEADER
	headers := []string{
		"Tanggal",
		"Uraian Transaksi",
		"Teller",
		"Debit",
		"Kredit",
		"Saldo",
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// DATA
	for i, d := range data {
		r := i + 2
		f.SetCellValue(sheet, "A"+itoa(r), d.Tanggal)
		f.SetCellValue(sheet, "B"+itoa(r), d.Uraian)
		f.SetCellValue(sheet, "C"+itoa(r), d.Teller)
		f.SetCellValue(sheet, "D"+itoa(r), toNumber(d.Debit))
		f.SetCellValue(sheet, "E"+itoa(r), toNumber(d.Kredit))
		f.SetCellValue(sheet, "F"+itoa(r), toNumber(d.Saldo))
	}

	// FORMAT
	f.SetColStyle(sheet, "D:F", styleNumber)
	f.SetColStyle(sheet, "B", wrapStyle)

	f.SetColWidth(sheet, "A", "A", 14)
	f.SetColWidth(sheet, "B", "B", 40)
	f.SetColWidth(sheet, "C", "C", 14)
	f.SetColWidth(sheet, "D", "F", 18)

	// SUMMARY
	start := len(data) + 3
	f.SetCellValue(sheet, "B"+itoa(start), "Saldo Awal")
	f.SetCellValue(sheet, "F"+itoa(start), sum.SaldoAwal)

	f.SetCellValue(sheet, "B"+itoa(start+1), "Total Debit")
	f.SetCellValue(sheet, "F"+itoa(start+1), sum.TotalDebit)

	f.SetCellValue(sheet, "B"+itoa(start+2), "Total Kredit")
	f.SetCellValue(sheet, "F"+itoa(start+2), sum.TotalKredit)

	f.SetCellValue(sheet, "B"+itoa(start+3), "Saldo Akhir")
	f.SetCellValue(sheet, "F"+itoa(start+3), sum.SaldoAkhir)

	os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err := f.SaveAs(outputPath); err != nil {
		log.Fatal(err)
	}
}

func WriteExcel(data []Mutasi, sum Summary) {
	f := excelize.NewFile()
	sheet := "Mutasi"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// ======================
	// STYLE
	// ======================

	// angka (numeric)
	styleNumber, _ := f.NewStyle(&excelize.Style{
		NumFmt: 4, // 2 desimal
	})

	// wrap text untuk uraian
	wrapStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
	})

	// ======================
	// HEADER
	// ======================
	headers := []string{
		"Tanggal",
		"Uraian Transaksi",
		"Teller",
		"Debit",
		"Kredit",
		"Saldo",
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// ======================
	// DATA
	// ======================
	for i, d := range data {
		r := i + 2
		f.SetCellValue(sheet, "A"+itoa(r), d.Tanggal)
		f.SetCellValue(sheet, "B"+itoa(r), d.Uraian)
		f.SetCellValue(sheet, "C"+itoa(r), d.Teller)
		f.SetCellValue(sheet, "D"+itoa(r), toNumber(d.Debit))
		f.SetCellValue(sheet, "E"+itoa(r), toNumber(d.Kredit))
		f.SetCellValue(sheet, "F"+itoa(r), toNumber(d.Saldo))
	}

	// ======================
	// FORMAT KOLUMN
	// ======================
	f.SetColStyle(sheet, "D:F", styleNumber)
	f.SetColStyle(sheet, "B", wrapStyle)

	f.SetColWidth(sheet, "A", "A", 14)
	f.SetColWidth(sheet, "B", "B", 40)
	f.SetColWidth(sheet, "C", "C", 14)
	f.SetColWidth(sheet, "D", "F", 18)

	// ======================
	// SUMMARY
	// ======================
	start := len(data) + 3

	f.SetCellValue(sheet, "B"+itoa(start), "Saldo Awal")
	f.SetCellValue(sheet, "F"+itoa(start), sum.SaldoAwal)

	f.SetCellValue(sheet, "B"+itoa(start+1), "Total Debit")
	f.SetCellValue(sheet, "F"+itoa(start+1), sum.TotalDebit)

	f.SetCellValue(sheet, "B"+itoa(start+2), "Total Kredit")
	f.SetCellValue(sheet, "F"+itoa(start+2), sum.TotalKredit)

	f.SetCellValue(sheet, "B"+itoa(start+3), "Saldo Akhir")
	f.SetCellValue(sheet, "F"+itoa(start+3), sum.SaldoAkhir)

	// ======================
	// SAVE
	// ======================
	os.MkdirAll("output", 0755)
	if err := f.SaveAs("output/output.xlsx"); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Excel selesai & rapi")
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
