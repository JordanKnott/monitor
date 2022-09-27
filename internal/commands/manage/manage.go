package manage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"strings"

	"github.com/jordanknott/monitor/internal/commands"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

var (
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
)

type Version struct {
	Version   string `json:"version"`
	BuildTime string `json:"buildTime"`
}

var manageCmd = &cobra.Command{
	Use:     "manage",
	Long:    "a tool for interacting with a WordPress install",
	Version: commands.VersionTemplate(),
	RunE:    newManageFn(),
}

func Execute() {
	manageCmd.AddCommand(newManageSnapshotCmd())
	manageCmd.Execute()
}

func newManageFn() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		raw, err := json.Marshal(Version{Version: sha1ver, BuildTime: buildTime})
		if err != nil {
			return err
		}
		fmt.Println(raw)
		old, err := ioutil.ReadFile("../wordpress-5.9.4/wp-includes/option.php")
		if err != nil {
			return err
		}

		new, err := ioutil.ReadFile("../option.php")
		if err != nil {
			return err
		}

		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(old), string(new), false)
		fmt.Println(diffs)

		var buff bytes.Buffer

		for _, diff := range diffs {
			text := strings.Replace(html.EscapeString(diff.Text), "\n", "<br>", -1)
			switch diff.Type {
			case diffmatchpatch.DiffInsert:
				_, _ = buff.WriteString("<ins class='diff-insert diff' style='background: #e6ffe6;'>")
				_, _ = buff.WriteString(text)
				_, _ = buff.WriteString("</ins>")
			case diffmatchpatch.DiffDelete:
				_, _ = buff.WriteString("<del class='diff-delete diff' style='background: #ffe6e6'>")
				_, _ = buff.WriteString(text)
				_, _ = buff.WriteString("</del>")
			case diffmatchpatch.DiffEqual:
				_, _ = buff.WriteString("<span class='diff'>")
				_, _ = buff.WriteString(text)
				_, _ = buff.WriteString("</span>")
			}

		}

		err = ioutil.WriteFile("test.html", buff.Bytes(), 0755)
		return err
	}
}
