package sdk

import (
	"strconv"
	"time"
)

// API Time for clean API nanoseconds to time.Time Objects
type APITime time.Time

func (tt *APITime) UnmarshalJSON(b []byte) (err error) {
	r := string(b)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(tt) = time.Unix(q/1000, 0)
	return nil
}

func (tt APITime) Time() time.Time {
	return time.Time(tt).UTC()
}

func (tt APITime) String() string {
	return tt.Time().String()
}

// Chunk big strings for log debug output
func Chunk(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	len := 0
	for _, r := range s {
		chunk[len] = r
		len++
		if len == chunkSize {
			chunks = append(chunks, string(chunk))
			len = 0
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return chunks
}
