package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var badTitleCharsRe = regexp.MustCompile(`[^a-zA-Z0-9]`)
var dashesRe = regexp.MustCompile(`-{2,}`)

// FatalOnErr prints error and exits if errir is not nil
func FatalOnErr(err error) {
	if err != nil {
		Fatal(err)
	}
}

// Fatal prints error and exits if errir
func Fatal(i ...interface{}) {
	fmt.Fprint(os.Stderr, "overmind: ")
	fmt.Fprintln(os.Stderr, i...)
	os.Exit(1)
}

// EscapeTitle makes title usable for tmux session name
func EscapeTitle(title string) string {
	return strings.ToLower(
		dashesRe.ReplaceAllString(badTitleCharsRe.ReplaceAllString(title, "-"), "-"),
	)
}

// RunCmd runs shell command and returns running error
func RunCmd(cmd string, args ...string) error {
	return exec.Command(cmd, args...).Run()
}

// SplitAndTrim splits string, trims every entry and removes blank entries
func SplitAndTrim(str string) (res []string) {
	split := strings.Split(str, ",")
	for _, s := range split {
		s = strings.Trim(s, " ")
		if len(s) > 0 {
			res = append(res, s)
		}
	}
	return
}

// StringsContain returns true if provided string slice contains provided string
func StringsContain(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

// WildcardMatch returns true if provided string matches provided wildcard
func WildcardMatch(pattern, str string) bool {
	re := regexp.MustCompile(
		fmt.Sprintf(
			"^%s$",
			strings.Replace(regexp.QuoteMeta(pattern), "\\*", ".*", -1),
		),
	)

	return re.MatchString(str)
}

// Max finds max integer
func Max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

// ScanLines reads line by line from reader. Doesn't throw "token too long" error like bufio.Scanner
func ScanLines(r io.Reader, callback func([]byte) bool) error {
	var (
		err  error
		line []byte
	)

	reader := bufio.NewReader(r)

	for {
		line, _, err = reader.ReadBytes('\n')
		if err != nil {
			break
		}
		if !callback(line) {
			return nil
		}
	}
	if err != io.EOF {
		return err
	}
	return nil
}

// ConvertError converts specific errors to the standard error type
func ConvertError(err error) error {
	if exErr, ok := err.(*exec.ExitError); ok {
		return errors.New(string(exErr.Stderr))
	}

	return err
}
