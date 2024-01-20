package cmd

import (
	"fmt"
	"os"

	"github.com/DhanushAdithya/hashnode-cli/internal/fetch"
	"github.com/aquasecurity/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func printTable(me fetch.Me) {
	t := table.New(os.Stdout)
	t.AddRow("Name", me.Data.Me.Name)
	t.AddRow("Username", me.Data.Me.Username)
	t.AddRow("Following", fmt.Sprintf("%d", me.Data.Me.FollowingsCount))
	t.AddRow("Followers", fmt.Sprintf("%d", me.Data.Me.FollowersCount))
	if len(me.Data.Me.NewPost.Nodes) == 0 {
		t.AddRow("Newest Post", "No posts yet")
	} else {
		t.AddRow("Newest Post", me.Data.Me.NewPost.Nodes[0].Title)
	}
	if len(me.Data.Me.OldPost.Nodes) == 0 {
		t.AddRow("Oldest Post", "No posts yet")
	} else {
		t.AddRow("Oldest Post", me.Data.Me.OldPost.Nodes[0].Title)
	}
	t.Render()
}

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetString("token") == "" {
			fmt.Println(
				lipgloss.
					NewStyle().
					Foreground(lipgloss.Color("#FF8080")).
					Render("No token set. Please run 'hashnode auth <token>' to set a token."),
			)
			return
		}
		data := fetch.MeResponse()
		if len(data.Errors) != 0 {
			fmt.Println(
				lipgloss.
					NewStyle().
					Foreground(lipgloss.Color("#FF8080")).
					Render(data.Errors[0].Message),
			)
			return
		}
		printTable(data)
	},
}

func init() {
	rootCmd.AddCommand(meCmd)
}
