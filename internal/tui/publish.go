package tui

import (
	"os"
	"strings"
	"sync"

	"github.com/DhanushAdithya/hashnode-cli/internal/fetch"
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
)

func Publish() {
	var (
		title       string = viper.GetString("title")
		tagsString  string = strings.Join(viper.GetStringSlice("tags"), ",")
		cover       string = viper.GetString("cover-image")
		publication string = viper.GetString("publicationId")
	)

	var wg sync.WaitGroup
	wg.Add(1)
	response := make(chan struct{})
	go func() {
		defer wg.Done()
		RenderLoad(response)
	}()

	me := fetch.MeResponse()
	if len(me.Errors) > 0 {
		close(response)
		wg.Wait()
		utils.RenderAPIErrors(me.Errors)
	}
	close(response)
	wg.Wait()

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

	wg.Add(1)
	response2 := make(chan struct{})
	go func() {
		defer wg.Done()
		RenderLoad(response2)
	}()

	content, err := os.ReadFile(viper.GetString("file"))
	if err != nil {
		utils.Exit("Unable to read file:", viper.GetString("file"), "\nError:", err)
	}
	tags := strings.Split(tagsString, ",")

	publish := fetch.PublishResponse(title, string(content), cover, publication, tags)
	if len(publish.Errors) > 0 {
		close(response2)
		wg.Wait()
		utils.RenderAPIErrors(publish.Errors)
	}
	close(response2)
	wg.Wait()

	utils.RenderSuccess("Post published successfully!")
}
