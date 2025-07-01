package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	go_git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const daysInLastSixMonths = 183
const outOfRange = 99999
const weeksInLastSixMonths = 26

type column []int

func stats(email string) {
	commits := processRepositories(email)
	printCommitsStats(commits)
}

func processRepositories(email string) map[int]int {
	filePath := getDotFilePath()
	repos := parseFileLinesToSlice(filePath)

	daysInMap := daysInLastSixMonths
	commits := make(map[int]int, daysInMap)

	for i := daysInMap; i > 0; i-- {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = fillCommits(email, path, commits)
	}
	return commits
}

// 指定したリポジトリパスからコミットを取得
func fillCommits(email string, path string, commits map[int]int) map[int]int {
	repo, err := go_git.PlainOpen(path)
	if err != nil {
		log.Fatal("リポジトリの取得に失敗しました :", path, err)
	}

	// HEADリファレンスを取得
	ref, err := repo.Head()
	if err != nil {
		log.Fatalf("最新情報の取得に失敗しました (%s): %v", path, err)
	}

	iterator, err := repo.Log(&go_git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatalf("コミットログの取得に失敗しました (%s):%v", path, err)
	}

	// GitHub風の表示
	offset := calcOffset()
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := CountDaysSinceDate(c.Author.When) + offset

		if c.Author.Email != email {
			return nil
		}
		if daysAgo != outOfRange {
			commits[daysAgo]++
		}
		return nil
	})
	if err != nil {
		log.Fatalf("コミット履歴のイテレーション中にエラーが発生しました (%s):%v", path, err)
	}
	return commits
}

func GetBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}

func CountDaysSinceDate(date time.Time) int {
	days := 0
	now := GetBeginningOfDay(time.Now())

	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++

		if days > daysInLastSixMonths {
			return outOfRange
		}
	}
	return days
}

// 必要な日数のオフセットを計算して返す
func calcOffset() int {
	return 0
}

// 集計されたコミットをグラフ形式で出力
func printCommitsStats(commits map[int]int) {
	keys := sortMapIntoSlice(commits)
	cols := buildCols(keys, commits)
	printCells(cols)
}

func sortMapIntoSlice(m map[int]int) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func buildCols(keys []int, commits map[int]int) map[int]column {
	cols := make(map[int]column)

	now := time.Now()
	// 今週の日曜日を基準点として設定
	today := GetBeginningOfDay(now)
	currentSunday := today.Add(time.Duration(-int(today.Weekday())) * 24 * time.Hour)

	// 各日について処理
	for _, k := range keys {
		commitDate := GetBeginningOfDay(now).Add(time.Duration(-k) * 24 * time.Hour)

		// その日が属する週の日曜日を計算
		commitSunday := commitDate.Add(time.Duration(-int(commitDate.Weekday())) * 24 * time.Hour)

		// 現在の日曜日からの週数を計算
		weeksDiff := int(currentSunday.Sub(commitSunday).Hours() / (24 * 7))

		// 週番号を設定（0が今週、1が先週、...）
		week := weeksDiff
		dayinweek := int(commitDate.Weekday())

		if week >= 0 && week <= weeksInLastSixMonths {
			if cols[week] == nil {
				cols[week] = make(column, 7)
			}
			cols[week][dayinweek] = commits[k]
		}
	}

	// 空の週も初期化
	for i := 0; i <= weeksInLastSixMonths; i++ {
		if cols[i] == nil {
			cols[i] = make(column, 7)
		}
	}

	return cols
}

func printCells(cols map[int]column) {
	printMonths()
	todayWeekdayIndex := int(time.Now().Weekday())

	for j := 0; j < 7; j++ {
		printDayCol(j)

		// 左から右へ（古い日付から新しい日付へ）表示
		for i := weeksInLastSixMonths; i >= 0; i-- {
			// 週の列にデータが存在するかチェック
			if col, ok := cols[i]; ok {
				if i == 0 && j == todayWeekdayIndex {
					printCell(col[j], true)
					continue
				} else {
					if len(col) > j {
						printCell(col[j], false)
						continue
					}
				}
			}
			printCell(0, false)
		}
		fmt.Printf("\n")
	}
}

// グラフの最初の行に月を表示させる
func printMonths() {
	var monthLine strings.Builder
	monthLine.WriteString("     ")

	now := time.Now()
	today := GetBeginningOfDay(now)
	currentSunday := today.Add(time.Duration(-int(today.Weekday())) * 24 * time.Hour)

	monthMarks := make(map[int]string)

	// 各週の開始日を確認して月の境界を見つける
	for i := weeksInLastSixMonths; i >= 0; i-- {
		weekStart := currentSunday.Add(time.Duration(-i*7) * 24 * time.Hour)
		prevWeekStart := currentSunday.Add(time.Duration(-(i+1)*7) * 24 * time.Hour)

		// 月が変わった場合、または最も古い週の場合
		if i == weeksInLastSixMonths || weekStart.Month() != prevWeekStart.Month() {
			monthMarks[i] = weekStart.Month().String()[:3]
		}
	}

	// 左から右へ（古い日付から新しい日付へ）表示
	for i := weeksInLastSixMonths; i >= 0; i-- {
		if month, ok := monthMarks[i]; ok {
			monthLine.WriteString(fmt.Sprintf("%-4s", month))
		} else {
			monthLine.WriteString("    ")
		}
	}
	fmt.Println(monthLine.String())
}

func printDayCol(day int) {
	out := "     "
	switch day {
	case 0:
		out = " Sun "
	case 1:
		out = " Mon "
	case 2:
		out = " Tue "
	case 3:
		out = " Wed "
	case 4:
		out = " Thu "
	case 5:
		out = " Fri "
	case 6:
		out = " Sat "
	}
	fmt.Printf("%s", out)
}

func printCell(val int, today bool) {
	escape := "\033[0;37;30m"

	if today {
		escape = "\033[1;37;45m"
	} else {
		switch {
		case val > 0 && val < 5:
			escape = "\033[38;5;17;48;5;153m"
		case val >= 5 && val < 10:
			escape = "\033[38;5;17;48;5;75m"
		case val >= 10 && val < 15:
			escape = "\033[38;5;18;48;5;33m"
		case val >= 15:
			escape = "\033[38;5;17;104m"
		}
	}

	if val == 0 {
		fmt.Printf(escape+"%-4s\033[0m", " - ")
		return
	}

	// コミット数に応じた数値の表示フォーマットを設定
	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	// エスケープシーケンスとフォーマットされた数値を標準出力に表示し、色をリセット
	fmt.Printf(escape+str+"\033[0m", val)
}
