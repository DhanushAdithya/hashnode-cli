package tui

import (
	"os"
	"strings"

	"github.com/DhanushAdithya/hashnode-cli/internal/fetch"
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
)

var (
	title       string
	tagsString  string
	cover       string
	publication string
	file        string
	path        string
	quit        bool
)

func loadFlags() {
	title = viper.GetString("title")
	tagsString = strings.Join(viper.GetStringSlice("tags"), ",")
	cover = viper.GetString("cover-image")
	publication = viper.GetString("publicationId")
	file = viper.GetString("file")
	path = viper.GetString("path")
}

type fileModel struct {
	filepicker.Model
	selected bool
}

func (m fileModel) Init() tea.Cmd {
	return m.Model.Init()
}

func (m fileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)
	if didSelect, path := m.Model.DidSelectFile(msg); didSelect {
		file = path
		m.selected = true
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			quit = true
			m.selected = true
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m fileModel) View() string {
	if m.selected {
		return ""
	}
	return m.Model.View()
}

func renderForm(me fetch.Me) {
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

func postArticle() {
	wg.Add(1)
	response := make(chan struct{})
	go func() {
		defer wg.Done()
		RenderLoad(response)
	}()

	content, err := os.ReadFile(file)
	if err != nil {
		utils.Exit("Unable to read file:", file, "\nError:", err)
	}
	tags := strings.Split(tagsString, ",")

	publish := fetch.PublishResponse(title, strings.TrimSpace(string(content)), cover, publication, tags)
	if len(publish.Errors) > 0 {
		close(response)
		wg.Wait()
		utils.RenderAPIErrors(publish.Errors)
	}
	close(response)
	wg.Wait()

	utils.RenderSuccess("Post published successfully!")
}

func renderFilePicker() {
	fp := filepicker.New()
	fp.AllowedTypes = []string{"md"}
	fp.CurrentDirectory = path

	model := fileModel{Model: fp}
	if _, err := tea.NewProgram(model).Run(); err != nil {
		utils.Exit(err)
	}
}

func Publish() {
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

	loadFlags()
	renderForm(me)
	if file == "" {
		renderFilePicker()
	}
	if quit {
		return
	}
	postArticle()
}
