# gitviz 👽
![Go Version](https://img.shields.io/badge/Go-1.24.2-00ADD8.svg?logo=go&logoColor=white)
![GitHub Release](https://img.shields.io/github/v/release/anton-fuji/gitviz)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`gitviz` is a command-line utility written in Go that visualizes your Git commit activity from local repositories, mimicking GitHub's contribution graph directly in your terminal.

# gitviz image 
![](demo/demo.png)

# Setup
If you have Homebrew installed, you can easily set up `gitviz` by following these steps.
### 1. Tap the `gitviz` Homebrew repository
```sh
brew tap anton-fuji/gitviz
```

### 2. Install gitviz
```sh
brew install gitviz
```
## Using `go install`  
If you prefer to install `gitviz` directly from source using Go.
### 1. Clone this Repository
```sh
git clone https://github.com/anton-fuji/gitviz.git 
cd gitviz
```

### 2. Initialize Go modules and download dependencies
```sh
go mod tidy
```

### 3. Install the executable
```sh
go install .
```
This command compiles gitviz and places the executable in your $GOPATH/bin, making it accessible from any directory in your terminal.

# Usage
`gitviz` provides two primary CLI options.
## 1. Scanning and Regisering Repositories
Before visualizing, you need to tell gitviz which repositories to monitor. This command scans a parent directory for Git repositories and saves their paths for future analysis.
```sh
gitviz -add /path/to/your/git/projects/directory
```

> [!NOTE]
> Registered repository paths are saved in a hidden file named `.gitlocalstats` in your home directory

## 2. Displaying the Contribution Graph
Once repositories are registered, you can generate and display your contribution graph by specifying your Git commit email address.
```sh
gitviz -graph your.email@example.com
```
- Replace `your-address@example.com` with the email address you use for your Git Commits.
