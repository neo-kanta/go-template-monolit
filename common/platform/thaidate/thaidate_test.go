package thaidate_test

import (
	"testing"
	"time"

	"go-transfer-agent/common/platform/thaidate"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToBE(t *testing.T) {
	// 1. Standard Gregorian Date
	gregorian := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	beStr := thaidate.ToBE(gregorian)
	assert.Equal(t, "2569-01-02", beStr)

	// 2. Zero Time
	zeroStr := thaidate.ToBE(time.Time{})
	assert.Equal(t, "", zeroStr)
}

func TestToBERFC3339(t *testing.T) {
	gregorian := time.Date(2026, 1, 2, 15, 30, 0, 0, time.UTC)
	beRFC := thaidate.ToBERFC3339(gregorian)
	assert.Equal(t, "2569-01-02T15:30:00Z", beRFC)

	zeroStr := thaidate.ToBERFC3339(time.Time{})
	assert.Equal(t, "", zeroStr)
}

func TestParseBE(t *testing.T) {
	// 1. Simple YYYY-MM-DD
	parsed, err := thaidate.ParseBE("2569-01-02")
	require.NoError(t, err)
	assert.Equal(t, 2026, parsed.Year())
	assert.Equal(t, time.January, parsed.Month())
	assert.Equal(t, 2, parsed.Day())

	// 2. RFC3339 format
	parsedRFC, err := thaidate.ParseBE("2569-01-02T15:30:00Z")
	require.NoError(t, err)
	assert.Equal(t, 2026, parsedRFC.Year())
	assert.Equal(t, 15, parsedRFC.Hour())

	// 3. Empty string
	emptyParsed, err := thaidate.ParseBE("")
	require.NoError(t, err)
	assert.True(t, emptyParsed.IsZero())

	// 4. Invalid formats
	_, err = thaidate.ParseBE("Invalid")
	assert.Error(t, err)

	_, err = thaidate.ParseBE("99999-99-99")
	assert.Error(t, err)
}
