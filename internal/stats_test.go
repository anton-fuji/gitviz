package gitviz

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetBeginningOfDay(t *testing.T) {
	inputTime1 := time.Date(2025, time.July, 1, 12, 30, 45, 0, time.Local)
	expectedOutput1 := time.Date(2025, time.July, 1, 0, 0, 0, 0, time.Local)
	actualOutput1 := GetBeginningOfDay(inputTime1)

	assert.Equal(t, expectedOutput1, actualOutput1, "Test Case 1 failed: Times should be equal at beginning of day")

	// すでに0時0分0秒の時刻を与えた場合
	inputTime2 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	expectedOutput2 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	actualOutput2 := GetBeginningOfDay(inputTime2)

	assert.Equal(t, expectedOutput2, actualOutput2, "Test Case 2 failed: Already beginning of day")

	// タイムゾーンが異なる場合
	jstLoc, _ := time.LoadLocation("Asia/Tokyo")
	inputTime3 := time.Date(2025, time.July, 1, 15, 0, 0, 0, jstLoc)
	expectedOutput3 := time.Date(2025, time.July, 1, 0, 0, 0, 0, jstLoc)
	actualOutput3 := GetBeginningOfDay(inputTime3)

	assert.Equal(t, expectedOutput3, actualOutput3, "Test Case 3 failed: Timezone should be preserved")
}

func TestCountDaysSinceDate(t *testing.T) {
	now := GetBeginningOfDay(time.Now())

	assert.Equal(t, 0, CountDaysSinceDate(now, DefaultDays), "CountDaysSinceDate(Today) failed: Expected 0")

	yesterday := now.Add(-24 * time.Hour)
	assert.Equal(t, 1, CountDaysSinceDate(yesterday, DefaultDays), "CountDaysSinceDate(Yesterday) failed: Expected 1")

	oneWeekAgo := now.Add(-7 * 24 * time.Hour)
	assert.Equal(t, 7, CountDaysSinceDate(oneWeekAgo, DefaultDays), "CountDaysSinceDate(OneWeekAgo) failed: Expected 7")

	farBeyondSixMonths := now.Add(-(DefaultDays + 10) * 24 * time.Hour)
	assert.Equal(t, outOfRange, CountDaysSinceDate(farBeyondSixMonths, DefaultDays), "CountDaysSinceDate(FarBeyondSixMonths) failed: Expected outOfRange")

	tomorrow := now.Add(24 * time.Hour)
	assert.Equal(t, 0, CountDaysSinceDate(tomorrow, DefaultDays), "CountDaysSinceDate(Tomorrow) failed: Expected 0 (for future date relative to current logic)")
}

func TestCountDaysSinceDateUsesMaxDays(t *testing.T) {
	now := GetBeginningOfDay(time.Now())
	threeDaysAgo := now.Add(-3 * 24 * time.Hour)

	assert.Equal(t, outOfRange, CountDaysSinceDate(threeDaysAgo, 2))
	assert.Equal(t, 3, CountDaysSinceDate(threeDaysAgo, 3))
}

func TestWeeksForDays(t *testing.T) {
	assert.Equal(t, 1, weeksForDays(1))
	assert.Equal(t, 1, weeksForDays(7))
	assert.Equal(t, 2, weeksForDays(8))
	assert.Equal(t, 27, weeksForDays(DefaultDays))
}

func TestNormalizeColorTheme(t *testing.T) {
	assert.Equal(t, DefaultColorTheme, normalizeColorTheme(""))
	assert.Equal(t, DefaultColorTheme, normalizeColorTheme("unknown"))
	assert.Equal(t, "blue", normalizeColorTheme(" blue "))
	assert.Equal(t, "purple", normalizeColorTheme("PURPLE"))
	assert.Equal(t, "orange", normalizeColorTheme("orange"))
	assert.Equal(t, "gray", normalizeColorTheme("gray"))
}

func TestCellColorSupportsThemesInBothModes(t *testing.T) {
	themes := []string{"green", "blue", "purple", "orange", "gray"}

	for _, theme := range themes {
		assert.NotEmpty(t, cellColor(1, false, theme), "block mode color should exist for %s", theme)
		assert.NotEmpty(t, cellColor(1, true, theme), "numbers mode color should exist for %s", theme)
	}
}
