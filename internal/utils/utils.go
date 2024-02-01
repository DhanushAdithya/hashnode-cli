package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

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
	ErrorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8080"))
	SuccessStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#80FF80"))
	InfoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#8080FF"))
	RowStyle       = lipgloss.NewStyle().Padding(0, 1)
	UsernameStyle  = RowStyle.Copy().Foreground(lipgloss.Color("244"))
	LabelColStyle  = RowStyle.Copy().Foreground(lipgloss.Color("99"))
	ActiveTabStyle = RowStyle.Copy().
			Background(lipgloss.Color("99")).
			Foreground(lipgloss.Color("#fff")).
			Margin(1, 0, 0, 2)
	InactiveTabStyle = RowStyle.Copy().Margin(1, 0, 0, 2).Background(lipgloss.Color("#353533"))
	TitleStyle       = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		// b.Left = "┤"
		b.Right = "├"
		return RowStyle.Copy().BorderStyle(b).Foreground(lipgloss.Color("99"))
	}()
	AuthorStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("99")).
			Foreground(lipgloss.Color("#fff")).
			Padding(0, 1).
			MarginLeft(2)

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

func Listify[T any](list []T) string {
	strList := make([]string, len(list))
	for i, v := range list {
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			strList[i] = fmt.Sprintf(`"%v"`, v)
		case reflect.Map:
			mapBytes, _ := json.Marshal(v)
			strList[i] = string(mapBytes)
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(strList, ","))
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
	RenderSuccess("Token set successfully")
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

func RenderSuccess(message string) {
	fmt.Println(SuccessStyle.Render(message))
}

func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func FormatDate(date string) string {
	parsedTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		Exit("Error parsing time:", err)
	}
	formattedDate := parsedTime.Format("Jan 2, 2006")
	return formattedDate
}

func RenderStatusBar(width int, author string, readTime int, published string, scrollPercent float64) string {
	name := AuthorStyle.Render(author)
	readTimeStr := UsernameStyle.Background(lipgloss.Color("#353533")).Render(fmt.Sprintf("%d min read", readTime))
	date := UsernameStyle.Background(lipgloss.Color("#353533")).Render(FormatDate(published))
	scrollPercentStr := UsernameStyle.Background(lipgloss.Color("#353533")).Render(fmt.Sprintf("%3.f%%", scrollPercent*100))
	spacer := lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.Color("#353533")).
		Render(strings.Repeat(" ", width-lipgloss.Width(name)-lipgloss.Width(readTimeStr)-lipgloss.Width(date)-lipgloss.Width(scrollPercentStr)-2))
	return "\n" + lipgloss.JoinHorizontal(
		lipgloss.Left,
		name,
		readTimeStr,
		spacer,
		date,
		scrollPercentStr,
	)
}

func FindIndex(list []string, item string) int {
	for i, v := range list {
		if strings.ToUpper(v) == strings.ToUpper(item) {
			return i
		}
	}
	return -1
}
