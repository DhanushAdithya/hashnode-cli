package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

type Error struct {
	Message    string   `json:"message"`
	Path       []string `json:"path"`
	Extensions struct {
		Code string `json:"code"`
	} `json:"extensions"`
}

var (
	ErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8080"))
	SuccessStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#80FF80"))
	InfoStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#8080FF"))
	RowStyle      = lipgloss.NewStyle().Padding(0, 1)
	UsernameStyle = RowStyle.Copy().Foreground(lipgloss.Color("244"))
	LabelColStyle = RowStyle.Copy().Foreground(lipgloss.Color("99"))
	TitleStyle    = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		// b.Left = "┤"
		b.Right = "├"
		return RowStyle.Copy().BorderStyle(b).Foreground(lipgloss.Color("99"))
	}()

	ErrorCodes = map[string]string{
		"GRAPHQL_VALIDATION_FAILED": "GraphQL query is invalid.",
		"UNAUTHENTICATED":           "No token set. Please run 'hashnode auth <token>' to set a token.",
		"FORBIDDEN":                 "Not allowed to access this resource.",
		"BAD_USER_INPUT":            "User input is invalid.",
		"NOT_FOUND":                 "Resource not found.",
	}
)

func RenderTitle(title string, width int) string {
	titleString := TitleStyle.Render(title)
	if (width < 0) || (width < lipgloss.Width(titleString)) {
		return titleString
	}
	// sideLength := (width - lipgloss.Width(titleString)) / 2
	sideLength := (width - lipgloss.Width(titleString))
	lines := strings.Repeat("─", sideLength)
	return lipgloss.JoinHorizontal(lipgloss.Center, titleString, lines)
}

func SetupConfig() {
	viper.SetConfigName("hashnode")
	viper.SetConfigType("yaml")
	homeDir, _ := os.UserHomeDir()
	configFile := filepath.Join(homeDir, "hashnode.yaml")
	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(configFile); err != nil {
				Exit("Unable to create config file:", err)
			}
		}
	}
	viper.AddConfigPath(homeDir)
	if err := viper.ReadInConfig(); err != nil {
		Exit("Unable to read config:", err)
	}
}

func Linkify(text, href string) string {
	return InfoStyle.Render(fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", href, text))
}

func Listify(list []string) string {
	for i, v := range list {
		list[i] = fmt.Sprintf(`"%s"`, v)
	}
	return fmt.Sprintf("[%s]", strings.Join(list, ","))
}

func Exit(message ...interface{}) {
	fmt.Println(ErrorStyle.Render(fmt.Sprint(message...)))
	os.Exit(1)
}

func CheckToken() {
	if viper.GetString("token") == "" {
		Exit(ErrorCodes["UNAUTHENTICATED"])
	}
}

func SetToken(token string) {
	viper.Set("token", token)
	if err := viper.WriteConfig(); err != nil {
		Exit("Unable to set Token to config file")
	}
	fmt.Println(SuccessStyle.Render("Token set successfully"))
}

func RenderAPIErrors(errors []Error) {
	for _, err := range errors {
		if msg, ok := ErrorCodes[err.Extensions.Code]; ok {
			fmt.Println(ErrorStyle.Render(msg))
		} else {
			fmt.Println(ErrorStyle.Render(err.Message))
		}
	}
	Exit()
}
