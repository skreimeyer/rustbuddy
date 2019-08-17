package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   "bump [flags]",
	Short: "Bump the current version",
	Long: `Bump locates the Cargo.toml and increments the version with options
for what type of change is desired in typical semantic-versioning. Major
version changes typically mean non-backward-compatible API changes. Minor 
changes are appropriate for new or changed functionality with backward
compatibility. Patch changes are appropriate for small fixes. Bump will default
to a patch upgrade.

Example:

Current version:0.1.2
-------------------------------
bump		0.1.3
bump --minor	0.2.0
bump --major	1.0.0
`,
	Run: func(cmd *cobra.Command, args []string) {
		bumpMain()
	},
}

var major bool
var minor bool
var dir string

func init() {
	rootCmd.AddCommand(bumpCmd)
	bumpCmd.Flags().BoolVar(&major, "major", false, "Increment by major version number.")
	bumpCmd.Flags().BoolVar(&minor, "minor", false, "Increment minor version number.")
	bumpCmd.Flags().StringVar(&dir, "dir", "./", "path to directory with Cargo.toml.")
}

func bumpMain() error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot read files in %s", dir)

	}
	var file os.FileInfo
	for _, f := range files {
		if strings.Contains(f.Name(), "Cargo.toml") {
			file = f
			break
		}
	}
	if len(file.Name()) == 0 {
		return errors.New("cannot find Cargo.toml in this directory")
	}
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return errors.New("cannot read Cargo.toml")
	}
	text := string(content)
	re := regexp.MustCompile(`version\s?=\s?"?\d+\.\d+\.\d+"?`)
	version := re.FindString(text)
	newVersion, err := bumper(version)
	if err != nil {
		return err
	}
	newText := strings.Replace(text, version, newVersion, 1)
	err = ioutil.WriteFile(file.Name(), []byte(newText), 0655)
	if err != nil {
		return err
	}
	return nil

}

func bumper(version string) (string, error) {
	fmt.Println("VERSION:", version)
	v := strings.TrimSpace(strings.Split(version, "=")[1])
	v = strings.Trim(v, `"`)
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return "", errors.New("non-standard versioning")
	}
	mj, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return "", errors.New("non-standard versioning")
	}
	mn, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return "", errors.New("non-standard versioning")
	}
	p, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return "", errors.New("non-standard versioning")
	}
	if major == true {
		return fmt.Sprintf(`version = "%d.%d.%d"`, mj+1, 0, 0), nil
	}
	if minor == true {
		return fmt.Sprintf(`version = "%d.%d.%d"`, mj, mn+1, 0), nil
	}
	return fmt.Sprintf(`version = "%d.%d.%d"`, mj, mn, p+1), nil
}
