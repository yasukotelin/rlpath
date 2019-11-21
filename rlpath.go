package rlpath

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/yasukotelin/tabwriter"

	"github.com/mattn/go-rl"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

const (
	padding       = 4
	defaultPrompt = "$ "
)

// Scanner can scan standard input with completetion path.
type Scanner struct {
	// Prompt is left edge text, like a $.
	Prompt string
	// RootDir is root directory to start scanning.
	// If this is empty, start scanning from execution path.
	RootDir string
	// OnlyDir is flag that shows only directory.
	OnlyDir bool
}

func (sc *Scanner) getRootDir() (string, error) {
	if sc.RootDir == "" {
		return "./", nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	newPath := strings.Replace(sc.RootDir, "~", usr.HomeDir, 1)
	return newPath, nil
}

// Scan scans path from standard input with completion.
func (sc *Scanner) Scan() (string, error) {
	// Change the current directory to specified root directory
	// for changing complete target.
	rootDir, err := sc.getRootDir()
	if err != nil {
		return "", err
	}
	curDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if err := os.Chdir(rootDir); err != nil {
		return "", err
	}
	defer os.Chdir(curDir)

	r := rl.NewRl()
	if sc.Prompt == "" {
		r.Prompt = defaultPrompt
	} else {
		r.Prompt = sc.Prompt
	}
	r.CompleteFunc = func(text string, pos int) (int, []string) {
		rs := []rune(text)
		start := pos
		for pos >= 0 {
			if pos == 0 || pos > 0 && rs[pos-1] == ' ' && (pos == 1 || rs[pos-2] != '\\') {
				v := string(rs[pos:])
				if runtime.GOOS == "windows" {
					v = strings.Replace(v, `/`, `\`, -1)
				}

				completePath, _ := getForCompleteFiles(v, sc.RootDir, sc.OnlyDir)

				if len(completePath.Paths) > 0 {
					for _, v := range completePath.Paths {
						if runtime.GOOS == "windows" {
							v = strings.Replace(v, `\`, `/`, -1)
						}
					}

					// Display a list if there are more completions.
					if len(completePath.Paths) != 1 {
						printCompleteFiles(completePath)
					}
					return pos, completePath.Paths
				} else {
					// Nothing completions.
					return start, []string{}
				}
			}
			pos--
		}
		return -1, nil
	}

	for {
		b, err := r.ReadLine()
		if err != nil {
			return "", err
		}
		s := string(b)
		return s, nil
	}
}

type completePath struct {
	Names []string
	Paths []string
}

func getForCompleteFiles(text string, rootDir string, onlyDir bool) (*completePath, error) {
	cmpPath := new(completePath)
	files, err := filepath.Glob(text + "*")
	if err != nil {
		return cmpPath, err
	}

	cmpPath.Paths = make([]string, 0, len(files))
	cmpPath.Names = make([]string, 0, len(files))
	for _, f := range files {
		finfo, err := os.Stat(f)
		if err != nil {
			return nil, err
		}
		if finfo.IsDir() {
			cmpPath.Names = append(cmpPath.Names, finfo.Name()+"/")
			cmpPath.Paths = append(cmpPath.Paths, f+"/")
		} else {
			if !onlyDir {
				cmpPath.Names = append(cmpPath.Names, finfo.Name())
				cmpPath.Paths = append(cmpPath.Paths, f)
			}
		}
	}

	return cmpPath, err
}

func printCompleteFiles(cmpPath *completePath) error {
	maxItemNum, err := computeMaxItemByTerminalWidth(cmpPath)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, padding, ' ', 0)
	fmt.Fprintln(w)
	for i, name := range cmpPath.Names {
		if i != 0 && i%maxItemNum == 0 {
			fmt.Fprintln(w, "\t")
		}
		fmt.Fprint(w, name+"\t")
	}
	fmt.Fprintln(w)

	return w.Flush()
}

// computeItemWidth returns max item number that can disp in the terminal on a row.
func computeMaxItemByTerminalWidth(cmpPath *completePath) (int, error) {
	termWidth, err := terminal.Width()
	if err != nil {
		return 0, err
	}

	// Gets a sorted slice by filename descend.
	names := make([]string, len(cmpPath.Names))
	copy(names, cmpPath.Names)
	sort.SliceStable(names, func(i, j int) bool {
		return len(names[i]) > len(names[j])
	})

	// Disps about 80% of Terminal width
	width := int(float64(termWidth) * 0.8)
	itemNum := 0
	for _, name := range names {
		nameLen := len(name)
		width = width - nameLen - padding
		if width <= 0 {
			break
		}
		itemNum++
	}

	return itemNum, nil
}
