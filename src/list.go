package src

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	todos    []string
	cursor   int
	selected map[int]string
	Download func(url string)
	newName  string
	show     bool
}

// ANSI转义码示例，设置前景色（文本颜色）
func colorText(text string, foregroundColor int) string {
	return fmt.Sprintf("\033[38;5;%dm%s\033[0m", foregroundColor, text)
}

// Init 初始化
func (m model) Init() tea.Cmd {
	return nil
}

// Update 更新
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.todos)-1 {
				m.cursor++
			}
		case " ", "enter":
			m.Download(m.selected[m.cursor])
			m.show = true
			return m, tea.Quit
		}
	}

	return m, nil
}

// View 显示
func (m model) View() string {
	if m.show {
		var s strings.Builder
		s.WriteString(fmt.Sprintf("%s %s\n", colorText(m.todos[m.cursor], 91), "✔️"))
		s.WriteString("输入下列命令启动项目: \n\n")
		s.WriteString(colorText(fmt.Sprintf("\t cd %s\n", m.newName), 35))
		s.WriteString(colorText("\t npm install\n", 35))
		s.WriteString(colorText("\t npm run serve\n", 35))
		s.WriteString("\n")
		return s.String()
	} else {
		s := "选择版本:\n\n"
		for i, choice := range m.todos {
			checked := " "
			if m.cursor == i {
				checked = ">"
			}

			if i == m.cursor {
				choice = colorText(choice, 92)
			}

			s += fmt.Sprintf("\t[%s] %s\n", colorText(checked, 92), choice)
		}
		s += "\n回车空格确认,q或ctrl+c退出\n"
		return s
	}
}

// json 结构体
type configItem struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// 模板数组
var config []configItem

// 读取 json 文件
func rend(filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取文件错误:", err)
		return err
	}

	// 解码JSON数组到切片
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("解析JSON错误:", err)
		return err
	}

	return nil
}

var todos []string
var selected = map[int]string{}

func init() {

	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("获取可执行文件路径错误:", err)
		return
	}

	// 获取当前可执行文件所在的目录
	exeDir := filepath.Dir(exePath)

	configPath := filepath.Join(exeDir, "config.json")

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("文件 config.json 不存在于目录:", exeDir)
		return
	}

	err = rend(configPath)
	if err != nil {
		fmt.Println(err)
	}

	for i, item := range config {
		todos = append(todos, item.Name)
		selected[i] = item.Url
	}
}

func List(Download func(url string), newName string) {
	var initModel = model{
		todos:    todos,
		selected: selected,
		Download: Download,
		newName:  newName,
		show:     false,
	}

	cmd := tea.NewProgram(initModel)
	if _, err := cmd.Run(); err != nil {
		fmt.Println("start failed:", err)
		os.Exit(1)
	}
}
