# gitviz 👽
![Go Version](https://img.shields.io/badge/Go-1.24.2-00ADD8.svg?logo=go&logoColor=white)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`gitviz` is a command-line utility written in Go that visualizes your Git commit activity from local repositories, mimicking GitHub's contribution graph directly in your terminal.

# Demo
![](demo/demo-gitviz.gif)

# Setup
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
> Registered repository paths are saved in a hidden file named `.gitlocalstatsin` in your home directory

## 2. Displaying the Contribution Graph
Once repositories are registered, you can generate and display your contribution graph by specifying your Git commit email address.
```sh
gitviz -email your.email@example.com
```
- Replace `your-address@example.com` with the email address you use for your Git Commits.
