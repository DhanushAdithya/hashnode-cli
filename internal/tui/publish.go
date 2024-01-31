package tui

import (
	"github.com/DhanushAdithya/hashnode-cli/internal/fetch"
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/charmbracelet/huh"
)

var (
	title       string
	tagsString  string
	cover       string
	content     string
	publication string
)

func Publish() {
	me := fetch.MeResponse()
	if len(me.Errors) > 0 {
		utils.RenderAPIErrors(me.Errors)
	}

	var publications []huh.Option[string]

	for _, pub := range me.Data.Me.Publications.Edges {
		publications = append(publications, huh.NewOption(pub.Node.Title, pub.Node.ID))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title").
				Placeholder("Article title...").
				Value(&title),
			huh.NewInput().
				Title("Tags").
				Description("Separate tags with a comma").
				Value(&tagsString),
			huh.NewInput().
				Title("Cover Image").
				Description("Enter url to the cover image").
				Value(&cover),
			huh.NewSelect[string]().
				Title("Publication").
				Description("Choose a publication to publish to").
				Options(publications...).
				Value(&publication),
		),
	)

	if err := form.Run(); err != nil {
		utils.Exit(err)
	}
}
