package gitviz

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	go_git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const DefaultDays = 183
const DefaultColorTheme = "green"
const outOfRange = 99999

const (
	ansiReset = "\033[0m"
	ansiMuted = "\033[38;5;240m"
	ansiToday = "\033[1;38;5;213m"
)

type column []int

type StatsOptions struct {
	Days    int
	Numbers bool
	Color   string
}

func Stats(email string, options StatsOptions) {
	if options.Days <= 0 {
		options.Days = DefaultDays
	}
	options.Color = normalizeColorTheme(options.Color)
	commits := processRepositories(email, options.Days)
	printCommitsStats(commits, options)
}

func processRepositories(email string, days int) map[int]int {
	filePath := getDotFilePath()
	repos := parseFileLinesToSlice(filePath)

	commits := make(map[int]int, days+1)

	for i := days; i >= 0; i-- {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = fillCommits(email, path, commits, days)
	}
	return commits
}

// 指定したリポジトリパスからコミットを取得
func fillCommits(email string, path string, commits map[int]int, days int) map[int]int {
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
		daysAgo := CountDaysSinceDate(c.Author.When, days) + offset

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

func CountDaysSinceDate(date time.Time, maxDays int) int {
	days := 0
	now := GetBeginningOfDay(time.Now())

	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++

		if days > maxDays {
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
func printCommitsStats(commits map[int]int, options StatsOptions) {
	keys := sortMapIntoSlice(commits)
	weeks := weeksForDays(options.Days)
	cols := buildCols(keys, commits, weeks)
	printCells(cols, weeks, options)
}

func sortMapIntoSlice(m map[int]int) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func buildCols(keys []int, commits map[int]int, weeks int) map[int]column {
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

		if week >= 0 && week <= weeks {
			if cols[week] == nil {
				cols[week] = make(column, 7)
			}
			cols[week][dayinweek] = commits[k]
		}
	}

	// 空の週も初期化
	for i := 0; i <= weeks; i++ {
		if cols[i] == nil {
			cols[i] = make(column, 7)
		}
	}

	return cols
}

func weeksForDays(days int) int {
	if days <= 0 {
		days = DefaultDays
	}
	return (days + 6) / 7
}

func printCells(cols map[int]column, weeks int, options StatsOptions) {
	numbers := options.Numbers
	printMonths(weeks, numbers)
	todayWeekdayIndex := int(time.Now().Weekday())

	for j := 0; j < 7; j++ {
		printDayCol(j, numbers)

		// 左から右へ（古い日付から新しい日付へ）表示
		for i := weeks; i >= 0; i-- {
			// 週の列にデータが存在するかチェック
			if col, ok := cols[i]; ok {
				if i == 0 && j == todayWeekdayIndex {
					printCell(col[j], true, options)
					continue
				} else {
					if len(col) > j {
						printCell(col[j], false, options)
						continue
					}
				}
			}
			printCell(0, false, options)
		}
		fmt.Printf("\n")
	}
	if numbers {
		printNumbersLegend(options.Color)
		return
	}
	printBlockLegend(options.Color)
	printSummary(cols)
}

// グラフの最初の行に月を表示させる
func printMonths(weeks int, numbers bool) {
	var monthLine strings.Builder
	monthLine.WriteString(dayLabelPadding(numbers))

	now := time.Now()
	today := GetBeginningOfDay(now)
	currentSunday := today.Add(time.Duration(-int(today.Weekday())) * 24 * time.Hour)

	monthMarks := make(map[int]string)

	// 各週の開始日を確認して月の境界を見つける
	for i := weeks; i >= 0; i-- {
		weekStart := currentSunday.Add(time.Duration(-i*7) * 24 * time.Hour)
		prevWeekStart := currentSunday.Add(time.Duration(-(i+1)*7) * 24 * time.Hour)

		// 月が変わった場合、または最も古い週の場合
		if i == weeks || weekStart.Month() != prevWeekStart.Month() {
			monthMarks[i] = weekStart.Month().String()[:3]
		}
	}

	// 左から右へ（古い日付から新しい日付へ）表示
	width := cellWidth(numbers)
	for i := weeks; i >= 0; i-- {
		if month, ok := monthMarks[i]; ok {
			monthLine.WriteString(fmt.Sprintf("%-*s", width, month))
		} else {
			monthLine.WriteString(strings.Repeat(" ", width))
		}
	}
	fmt.Println(ansiMuted + monthLine.String() + ansiReset)
}

func printDayCol(day int, numbers bool) {
	out := dayLabelPadding(numbers)
	switch day {
	case 0:
		out = formatDayLabel("Sun", numbers)
	case 1:
		out = formatDayLabel("Mon", numbers)
	case 2:
		out = formatDayLabel("Tue", numbers)
	case 3:
		out = formatDayLabel("Wed", numbers)
	case 4:
		out = formatDayLabel("Thu", numbers)
	case 5:
		out = formatDayLabel("Fri", numbers)
	case 6:
		out = formatDayLabel("Sat", numbers)
	}
	fmt.Printf("%s%s%s", ansiMuted, out, ansiReset)
}

func printCell(val int, today bool, options StatsOptions) {
	style := cellColor(val, options.Numbers, options.Color)
	if today {
		style += ansiToday
	}
	if options.Numbers {
		fmt.Printf("%s%s%s", style, numberCellLabel(val), ansiReset)
		return
	}
	fmt.Printf("%s%s%s", style, blockCellLabel(val, today), ansiReset)
}

func numberCellLabel(val int) string {
	switch {
	case val == 0:
		return " ·  "
	case val > 99:
		return "99+ "
	default:
		return fmt.Sprintf("%3d ", val)
	}
}

func blockCellLabel(val int, today bool) string {
	if today {
		return "◆ "
	}
	if val == 0 {
		return "· "
	}
	return "■ "
}

func cellColor(val int, numbers bool, theme string) string {
	if !numbers {
		colors := blockColorTheme(theme)
		switch {
		case val == 0:
			return "\033[38;5;238m"
		case val < 5:
			return ansiFg(colors[0])
		case val < 10:
			return ansiFg(colors[1])
		case val < 15:
			return ansiFg(colors[2])
		default:
			return ansiFg(colors[3])
		}
	}

	colors := numberColorTheme(theme)
	switch {
	case val == 0:
		return "\033[38;5;245m"
	case val < 5:
		return ansiFgBg(colors.fgDark, colors.backgrounds[0])
	case val < 10:
		return ansiFgBg(colors.fgDark, colors.backgrounds[1])
	case val < 15:
		return ansiFgBg(colors.fgLight, colors.backgrounds[2])
	default:
		return ansiFgBg(colors.fgLight, colors.backgrounds[3])
	}
}

func normalizeColorTheme(theme string) string {
	switch strings.ToLower(strings.TrimSpace(theme)) {
	case "blue", "purple", "orange", "gray":
		return strings.ToLower(strings.TrimSpace(theme))
	default:
		return DefaultColorTheme
	}
}

func blockColorTheme(theme string) []int {
	switch theme {
	case "blue":
		return []int{153, 75, 33, 27}
	case "purple":
		return []int{183, 141, 99, 57}
	case "orange":
		return []int{222, 215, 208, 166}
	case "gray":
		return []int{250, 246, 242, 238}
	default:
		return []int{155, 119, 77, 40}
	}
}

type numberColors struct {
	fgDark      int
	fgLight     int
	backgrounds []int
}

func numberColorTheme(theme string) numberColors {
	switch theme {
	case "blue":
		return numberColors{fgDark: 17, fgLight: 255, backgrounds: []int{153, 75, 33, 27}}
	case "purple":
		return numberColors{fgDark: 17, fgLight: 255, backgrounds: []int{183, 141, 99, 57}}
	case "orange":
		return numberColors{fgDark: 94, fgLight: 255, backgrounds: []int{222, 215, 208, 166}}
	case "gray":
		return numberColors{fgDark: 232, fgLight: 255, backgrounds: []int{250, 246, 242, 238}}
	default:
		return numberColors{fgDark: 232, fgLight: 255, backgrounds: []int{115, 79, 38, 24}}
	}
}

func ansiFg(color int) string {
	return fmt.Sprintf("\033[38;5;%dm", color)
}

func ansiFgBg(fg int, bg int) string {
	return fmt.Sprintf("\033[38;5;%d;48;5;%dm", fg, bg)
}

func cellWidth(numbers bool) int {
	if numbers {
		return 4
	}
	return 2
}

func dayLabelPadding(numbers bool) string {
	if numbers {
		return "    "
	}
	return "   "
}

func formatDayLabel(label string, numbers bool) string {
	if numbers {
		return fmt.Sprintf("%-4s", label)
	}
	return fmt.Sprintf("%-3s", label)
}

func printBlockLegend(theme string) {
	fmt.Printf("\n%s   Less %s·%s %s■%s %s■%s %s■%s %s■%s More   %s◆%s Today\n",
		ansiMuted,
		cellColor(0, false, theme), ansiReset,
		cellColor(1, false, theme), ansiReset,
		cellColor(5, false, theme), ansiReset,
		cellColor(10, false, theme), ansiReset,
		cellColor(15, false, theme), ansiReset,
		ansiToday, ansiReset,
	)
}

func printNumbersLegend(theme string) {
	fmt.Printf("\n%s    Less %s ·  %s %s  1 %s %s  5 %s %s 10 %s %s 15 %s More   %s  0 %s Today\n",
		ansiMuted,
		cellColor(0, true, theme), ansiReset,
		cellColor(1, true, theme), ansiReset,
		cellColor(5, true, theme), ansiReset,
		cellColor(10, true, theme), ansiReset,
		cellColor(15, true, theme), ansiReset,
		ansiToday, ansiReset,
	)
}

func printSummary(cols map[int]column) {
	total := 0
	activeDays := 0
	maxCommits := 0

	for _, col := range cols {
		for _, commits := range col {
			total += commits
			if commits > 0 {
				activeDays++
			}
			if commits > maxCommits {
				maxCommits = commits
			}
		}
	}

	fmt.Printf("%s   Total %d commits   Active %d days   Max %d/day%s\n", ansiMuted, total, activeDays, maxCommits, ansiReset)
}
