package main

import (
	"flag"
	"log"
	"os/exec"
	"strings"

	gitviz "github.com/anton-fuji/gitviz/internal"
)

func main() {
	var folder string
	var graphEmail string
	var days int
	var numbers bool
	var color string
	flag.StringVar(&folder, "add", "", "add a new folder to scan for Git repositories")
	flag.StringVar(&graphEmail, "graph", "", "the email to scan")
	flag.IntVar(&days, "days", gitviz.DefaultDays, "the number of days to include in the graph")
	flag.BoolVar(&numbers, "numbers", false, "show commit counts in each graph cell")
	flag.StringVar(&color, "color", gitviz.DefaultColorTheme, "graph color theme: green, blue, purple, orange, gray")
	flag.Parse()

	if folder != "" {
		gitviz.Scan(folder)
		return
	}
	if days <= 0 {
		log.Fatal("days must be greater than 0")
	}
	if graphEmail == "" {
		graphEmail = inferGitEmail()
	}
	gitviz.Stats(graphEmail, gitviz.StatsOptions{
		Days:    days,
		Numbers: numbers,
		Color:   color,
	})
}

func inferGitEmail() string {
	out, err := exec.Command("git", "config", "user.email").Output()
	email := strings.TrimSpace(string(out))
	if err != nil || email == "" {
		log.Fatal("email is required: pass -graph or set git config user.email")
	}
	return email
}
