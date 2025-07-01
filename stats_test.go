package main

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

	assert.Equal(t, 0, CountDaysSinceDate(now), "CountDaysSinceDate(Today) failed: Expected 0")

	yesterday := now.Add(-24 * time.Hour)
	assert.Equal(t, 1, CountDaysSinceDate(yesterday), "CountDaysSinceDate(Yesterday) failed: Expected 1")

	oneWeekAgo := now.Add(-7 * 24 * time.Hour)
	assert.Equal(t, 7, CountDaysSinceDate(oneWeekAgo), "CountDaysSinceDate(OneWeekAgo) failed: Expected 7")

	farBeyondSixMonths := now.Add(-(daysInLastSixMonths + 10) * 24 * time.Hour)
	assert.Equal(t, outOfRange, CountDaysSinceDate(farBeyondSixMonths), "CountDaysSinceDate(FarBeyondSixMonths) failed: Expected outOfRange")

	tomorrow := now.Add(24 * time.Hour)
	assert.Equal(t, 0, CountDaysSinceDate(tomorrow), "CountDaysSinceDate(Tomorrow) failed: Expected 0 (for future date relative to current logic)")
}
