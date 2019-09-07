package record

import (
	"fmt"
	"github.com/metakeule/fmtdate"
	"github.com/shakinm/xlsReader/helpers"
	"github.com/shakinm/xlsReader/xls/structure"
	"strings"
)

//FORMAT: Number Format

var FormatRecord = []byte{0x1E, 0x04} //(41Eh)

/*
The FORMAT record describes a number format in the workbook.
All the FORMAT records should appear together in a BIFF file. The order of FORMAT
records in an existing BIFF file should not be changed. It is possible to write custom
number formats in a file, but they should be added at the end of the existing FORMAT
records.

Record Data
Offset		Field Name		Size		Contents
------------------------------------------------
4			ifmt			2			Format index code (for internal use only)
6			cch				2			Length of the string
7			grbit			1			Option Flags (described in Unicode Strings in BIFF8 section)
8			rgb				var			Array of string characters

Excel uses the ifmt structure to identify built-in formats when it reads a file that was
created by a different localized version. For more information about built-in formats,
see "XF".

*/

type Format struct {
	ifmt     [2]byte
	stFormat structure.XLUnicodeRichExtendedString
}

func (r *Format) Read(stream []byte) {
	copy(r.ifmt[:], stream[0:2])
	r.stFormat.Read(stream[2:])

}
func (r *Format) GetIndex() int {
	return int(helpers.BytesToUint16(r.ifmt[:]))
}

func (r *Format) GetFormatString(data structure.CellData) string {
	if r.GetIndex() > 164 {
		if data.GetType() == "*record.BoolErr" {

		return data.GetString()
		}
		if data.GetType() == "*record.Number" {
			if strings.Contains(r.stFormat.String(), "0.00") {
				return fmt.Sprintf("%.2f", data.GetFloat64()*100) + "%"
			}
		}
		if r.stFormat.String() == "General" {
			return data.GetString()
		}
		t := helpers.TimeFromExcelTime(data.GetFloat64(), false)
		dateFormat := strings.ReplaceAll(r.stFormat.String(), "HH:MM:SS", "hh:mm:ss")
		dateFormat = strings.ReplaceAll(dateFormat, "\\", "")
		return fmtdate.Format(dateFormat, t)
	}
	return data.GetString()
}
