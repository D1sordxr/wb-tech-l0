package tools

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToText(val string) pgtype.Text {
	if val == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: val, Valid: true}
}

func ToInt4(val int32) pgtype.Int4 {
	if val == 0 {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: val, Valid: true}
}

func ToInt8(val int64) pgtype.Int8 {
	if val == 0 {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: val, Valid: true}
}

func ToTimestamp(t time.Time) pgtype.Timestamp {
	if t.IsZero() {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{Time: t, Valid: true}
}
