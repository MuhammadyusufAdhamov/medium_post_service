package utils

import "database/sql"

func NUllString(s string) (ns sql.NullString) {
	if s != "" {
		ns.String = s
		ns.Valid = true
	}
	return ns
}

func NUllFloat64(s float64) (ns sql.NullFloat64) {
	if s != 0 {
		ns.Float64 = s
		ns.Valid = true
	}
	return ns
}

func NUllInt32(s int32) (ns sql.NullInt32) {
	if s != 0 {
		ns.Int32 = s
		ns.Valid = true
	}
	return ns
}

func FormatNUllTime(nt sql.NullTime, format string) string {
	if nt.Valid {
		return nt.Time.Format(format)
	}
	return ""
}