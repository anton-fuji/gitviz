package main

import (
	"fmt"
	"log"
	"sort"
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
		daysAgo := countDaysSinceDate(c.Author.When) + offset

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

func getBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}

func countDaysSinceDate(date time.Time) int {
	days := 0
	now := getBeginningOfDay(time.Now())

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
	var offset int
	weekday := time.Now().Weekday()
	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}
	return offset
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
	// 処理中の週のコミット数を一時的に保存
	col := column{}
	for _, k := range keys {
		week := int(k / 7)
		dayinweek := k % 7
		if dayinweek == 0 {
			col = column{}
		}
		col = append(col, commits[k])
		if dayinweek == 6 {
			cols[week] = col
		}
	}
	return cols
}

func printCells(cols map[int]column) {
	printMonths()

	// 日曜から金曜まで行を処理
	for j := 6; j >= 0; j-- {
		for i := weeksInLastSixMonths + 1; i >= 0; i-- {
			if i == weeksInLastSixMonths+1 {
				printDayCol(j)
			}

			// 週の列(i)にデータが存在するかチェック
			if col, ok := cols[i]; ok {
				if i == 0 && j == calcOffset()-1 {
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
	week := getBeginningOfDay(time.Now()).Add(-(time.Duration(daysInLastSixMonths) * time.Hour * 24))
	month := week.Month()
	fmt.Printf("     ")

	for {
		if week.Month() != month {
			fmt.Printf("%s", week.Month().String()[:3])
			fmt.Printf("    ")
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}
		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}

func printDayCol(day int) {
	out := "     "
	switch day {
	case 0: // 日曜日
		out = " Sun "
	case 1: // 月曜日
		out = " Mon "
	case 2: // 火曜日
		out = " Tue "
	case 3: // 水曜日
		out = " Wed "
	case 4: // 木曜日
		out = " Thu "
	case 5: // 金曜日
		out = " Fri "
	case 6: // 土曜日
		out = " Sat "
	}
	fmt.Printf("%s", out)
}

func printCell(val int, today bool) {
	escape := "\033[0;37;30m"

	// コミット数のカラー
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;46m"
	case val >= 5 && val < 10:
		escape = "\033[1;37;44m"
	case val >= 10 && val < 20:
		escape = "\033[1;97;44m"
	case val >= 20:
		escape = "\033[1;97;104m"
	}

	if today {
		escape = "\033[1;37;45m"
	}

	if val == 0 {
		fmt.Printf("\033[0;37;40m%-4s\033[0m", " - ")
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
