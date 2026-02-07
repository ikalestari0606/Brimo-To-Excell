package main

import (
	"regexp"
	"strconv"
	"strings"
)

func ParseLines(lines []string) []Mutasi {
	var out []Mutasi
	var cur *Mutasi
	var buf []string

	for _, raw := range lines {
		l := strings.TrimSpace(raw)
		if l == "" {
			continue
		}

		// 🔹 transaksi baru jika ketemu tanggal
		if isDate(l) {
			if cur != nil {
				finalize(cur, buf)

				// 🚫 filter BRISIM SAJA (sesuai permintaan sebelumnya)
				if !strings.Contains(strings.ToUpper(cur.Uraian), "BRISIM") {
					out = append(out, *cur)
				}
			}

			cur = &Mutasi{Tanggal: l}
			buf = nil
			continue
		}

		// 🔹 SEMUA BARIS SETELAH TANGGAL DIANGGAP BAGIAN TRANSAKSI
		if cur != nil {
			buf = append(buf, l)
		}
	}

	// 🔹 transaksi terakhir (WAJIB)
	if cur != nil {
		finalize(cur, buf)

		if !strings.Contains(strings.ToUpper(cur.Uraian), "BRISIM") {
			out = append(out, *cur)
		}
	}

	return out
}

// ======================
// FILTER HEADER
// ======================
func isIgnoredHeader(line string) bool {
	l := strings.ToUpper(strings.TrimSpace(line))

	return strings.HasPrefix(l, "KEPADA YTH") ||
		strings.HasPrefix(l, "TO :") ||
		strings.Contains(l, "TANGGAL LAPORAN") ||
		strings.Contains(l, "STATEMENT DATE") ||

		// 🚫 HEADER PERUSAHAAN & PERIODE
		strings.Contains(l, "PT MIN GOOK INDONESIA") ||
		strings.Contains(l, "PERIODE TRANSAKSI") ||
		strings.Contains(l, "07/01/26") ||
		strings.Contains(l, "TRANSACTION PERIOD :")
}

// ======================
// FINALIZE TRANSAKSI
// ======================
func finalize(m *Mutasi, parts []string) {
	var desc []string
	var nums []string

	reTeller := regexp.MustCompile(`^0\d{6}$`)
	reAmount := regexp.MustCompile(`\d{1,3}(?:[.,]\d{3})*(?:[.,]\d{2})`)

	isOpening := false

	for _, p := range parts {
		up := strings.ToUpper(p)

		// deteksi opening balance / saldo awal
		if strings.Contains(up, "OPENING BALANCE") ||
			strings.Contains(up, "SALDO AWAL") {
			isOpening = true
			desc = append(desc, p)
			continue
		}

		// teller dari kolom teller PDF
		if m.Teller == "" && reTeller.MatchString(p) {
			m.Teller = p
			continue
		}

		// ambil semua angka nominal
		matches := reAmount.FindAllString(p, -1)
		if len(matches) > 0 {
			nums = append(nums, matches...)
			continue
		}

		// sisanya uraian
		desc = append(desc, p)
	}

	// rapikan uraian (maks 3 baris)
	if len(desc) > 3 {
		desc = desc[:3]
	}
	m.Uraian = strings.Join(desc, "\n")

	// KHUSUS OPENING BALANCE
	if isOpening {
		// debit & kredit memang tidak ada → biarkan kosong
		if len(nums) > 0 {
			m.Saldo = nums[len(nums)-1]
		}
		return
	}

	// TRANSAKSI NORMAL
	if len(nums) >= 3 {
		m.Debit = nums[len(nums)-3]
		m.Kredit = nums[len(nums)-2]
		m.Saldo = nums[len(nums)-1]
	}
}

// ======================
// SUMMARY (OPSIONAL)
// ======================
func ParseSummary(lines []string) Summary {
	var nums []float64
	re := regexp.MustCompile(`\d{1,3}(,\d{3})*\.\d{2}`)

	for _, l := range lines {
		matches := re.FindAllString(l, -1)
		for _, m := range matches {
			nums = append(nums, parseAmount(m))
		}
	}

	n := len(nums)
	if n < 4 {
		return Summary{}
	}

	return Summary{
		SaldoAwal:   nums[n-4],
		TotalDebit:  nums[n-3],
		TotalKredit: nums[n-2],
		SaldoAkhir:  nums[n-1],
	}
}

// ======================
// UTIL ANGKA
// ======================
func parseAmount(s string) float64 {
	s = strings.ReplaceAll(s, ",", "")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func toNumber(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	s = strings.ReplaceAll(s, ",", "")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func isAmount(s string) bool {
	return regexp.MustCompile(`^[\d.,]+$`).MatchString(s)
}

func isDate(s string) bool {
	return regexp.MustCompile(`\d{2}/\d{2}/\d{2}`).MatchString(s)
}
