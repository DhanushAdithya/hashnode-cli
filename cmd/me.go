package cmd

import (
	"fmt"
	"sync"

	"github.com/DhanushAdithya/hashnode-cli/internal/fetch"
	"github.com/DhanushAdithya/hashnode-cli/internal/tui"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8080"))

func printTable(me fetch.Me) {
	baseStyle := lipgloss.NewStyle().Padding(0, 1)
	labelStyle := baseStyle.Copy().Foreground(lipgloss.Color("99"))
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
				return baseStyle.Copy().Foreground(lipgloss.Color("244"))
			}
			switch col {
			case 0:
				return labelStyle
			}
			return baseStyle
		})
	fmt.Println(t)
}

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetString("token") == "" {
			fmt.Println(
				errorStyle.Render("No token set. Please run 'hashnode auth <token>' to set a token."),
			)
			return
		}
		var wg sync.WaitGroup
		wg.Add(1)
		response := make(chan struct{})
		go func() {
			defer wg.Done()
			tui.RenderLoad(response)
		}()
		data := fetch.MeResponse()
		if len(data.Errors) != 0 {
			close(response)
			wg.Wait()
			fmt.Println(errorStyle.Render(data.Errors[0].Message))
			return
		}
		close(response)
		wg.Wait()
		printTable(data)
	},
}

func init() {
	rootCmd.AddCommand(meCmd)
}
