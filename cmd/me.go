package cmd

import (
	"fmt"
	"sync"

	"github.com/DhanushAdithya/hashnode-cli/internal/fetch"
	"github.com/DhanushAdithya/hashnode-cli/internal/tui"
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

func printTable(me fetch.Me) {
	rows := [][]string{
		{"NAME", me.Data.Me.Name},
		{"USERNAME", "@" + me.Data.Me.Username},
		{"FOLLOWING", fmt.Sprintf("%d", me.Data.Me.FollowingsCount)},
		{"FOLLOWERS", fmt.Sprintf("%d", me.Data.Me.FollowersCount)},
	}
	if len(me.Data.Me.NewPost.Nodes) == 0 {
		rows = append(rows, []string{"NEWEST POST", "No posts yet"})
	} else {
		rows = append(rows, []string{"NEWEST POST", me.Data.Me.NewPost.Nodes[0].Title})
	}
	if len(me.Data.Me.OldPost.Nodes) == 0 {
		rows = append(rows, []string{"OLDEST POST", "No posts yet"})
	} else {
		rows = append(rows, []string{"OLDEST POST", me.Data.Me.OldPost.Nodes[0].Title})
	}
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		BorderRow(true).
		Width(70).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 2 && col == 1 {
				return utils.UsernameStyle
			}
			switch col {
			case 0:
				return utils.LabelColStyle
			}
			return utils.RowStyle
		})
	fmt.Println(t)
}

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Look up your Hashnode profile",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckToken()
		var wg sync.WaitGroup
		wg.Add(1)
		response := make(chan struct{})
		go func() {
			defer wg.Done()
			tui.RenderLoad(response)
		}()
		data := fetch.MeResponse()
		if len(data.Errors) > 0 {
			close(response)
			wg.Wait()
			utils.RenderAPIErrors(data.Errors)
		}
		close(response)
		wg.Wait()
		printTable(data)
	},
}

func init() {
	rootCmd.AddCommand(meCmd)
}
