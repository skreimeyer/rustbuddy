package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/skreimeyer/rustbuddy/crates"
	"github.com/spf13/cobra"
)

// treeCmd represents the tree command
var treeCmd = &cobra.Command{
	Use:   "tree [CRATE]",
	Short: "Print a crate's dependency tree",
	Long: `Tree takes the name of a crate as an argument, calls the crates.io
api and then prints the dependency tree. The tree will include useful
information such as the most recent version, release cadence and date of last
publication.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		makeTree(args)
	},
}

func init() {
	rootCmd.AddCommand(treeCmd)
}

func makeTree(args []string) {
	// data, err := crates.FetchCrate(args[0])
	// if err != nil {
	// 	fmt.Println("Could not fetch crate:", args[0])
	// 	fmt.Println(err)
	// }
	const header = `
+----------------------+------------+------------+------------+----------------+
|lvl|       Name       |  Version   |  Last Ver  | Pub Cadence|    Last Pub    |
|---|------------------|------------|------------|------------|----------------|
`
	fmt.Print(header)
	err := makeAllRows(args[0], "", 0)
	if err != nil {
		fmt.Println("Returned error:", err)
	}
	// var output string
	// output += header
	// top, err := makeRow(data, "", 0)
	// if err != nil {
	// 	fmt.Println("cannot extract row data. exiting...")
	// 	return
	// }
	// output += top + "\n"
	// fmt.Print(output)
	fmt.Println(`+----------------------+------------+------------+------------+----------------+`)
	return
}

func makeAllRows(crate, ver string, depth int) error {
	data, err := crates.FetchCrate(crate)
	if err != nil {
		fmt.Println(err)
		return err
	}
	row, err := makeRow(data, ver, depth)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(row)
	depth++
	if ver == "" {
		ver = data.Versions[0].Num
	}
	deps, err := crates.FetchDeps(crate, ver)
	if len(deps) == 0 {
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, d := range deps {
		cName := d.CrateId
		cVer := toSemVer(d.Req)
		err = makeAllRows(cName, cVer, depth)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	return nil
}

func makeRow(data crates.CrateData, ver string, depth int) (string, error) {
	var row string
	name := data.Crate.Name
	if len(data.Versions) == 0 {
		return row, errors.New("no versions")
	}
	if ver == "" {
		ver = data.Versions[0].Num
	}
	lastVer := data.Versions[0].Num
	cad := cadence(data.Versions)
	lastPub, err := time.Parse(layout, data.Versions[0].UpdatedAt)
	if err != nil {
		return row, errors.New("cannot parse pub time")
	}
	y, m, d := lastPub.Date()
	latest := fmt.Sprintf("%v %v %v", y, m, d)

	row += fmt.Sprintf("|%03d", depth)
	row += "|"
	row += center(name, 18)
	row += "|"
	row += center(ver, 12)
	row += "|"
	row += center(lastVer, 12)
	row += "|"
	row += center(cad, 12)
	row += "|"
	row += center(latest, 16)
	row += "|"
	return row, nil
}

func center(s string, n int) string {
	var result string
	c := utf8.RuneCountInString(s)
	if c > n {
		for i, r := range s {
			if i == n {
				break
			}
			result += string(r)
		}
		return result
	}
	frontPad := (n - c) / 2
	backPad := n - frontPad - c
	result = strings.Repeat(" ", frontPad) + s + strings.Repeat(" ", backPad)
	return result
}

const layout = "2006-01-02T15:04:05.000000+00:00" // for time parser

func cadence(vers []crates.Versions) string {
	if len(vers) < 2 {
		return "N/A"
	}
	count := 0
	var sumDelta time.Duration
	last := time.Time{}
	for _, v := range vers {
		t, err := time.Parse(layout, v.UpdatedAt)
		if err != nil {
			continue
		}
		count++
		if last.IsZero() {
			last = t
			continue
		}
		sumDelta += last.Sub(t)
		last = t
	}
	avg := sumDelta.Hours() / float64(count)
	switch {
	case avg < 48.0:
		return "daily"
	case avg < 336.0:
		return "weekly"
	case avg < 1440.0:
		return "monthly"
	case avg < 4380.0:
		return "quarterly"
	case avg < 17520.0:
		return "annually"
	default:
		return "rarely"

	}
}

// This isn't technically correct, but finding the highest patch version would
// require an API call to the base crate, finding the highest version (so a
// strconv and sort or loop), THEN making an API call to the version we want.
// The information almost certainly isn't important enough to warrant the added
// latency. Conventionally, selecting the highest patch version is typical,
// which *shouldn't* entail a modification of dependencies in most cases.
func toSemVer(s string) string {
	s = strings.TrimPrefix(s, "^")
	components := strings.Split(s, ".")
	for {
		if len(components) < 3 {
			components = append(components, "0")
		}
		break
	}
	return strings.Join(components, ".")

}
