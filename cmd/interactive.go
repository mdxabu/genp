/*
Copyright © 2025 - github.com/mdxabu
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/mdxabu/genp/internal"
	"github.com/mdxabu/genp/internal/store"
	"github.com/spf13/cobra"
)

// Command represents an available command in the interactive mode
type Command struct {
	Name        string
	Description string
	Execute     func() tea.Msg
}

// getCommands returns the available commands in interactive mode
func getCommands() []Command {
	return []Command{
		{
			Name:        "create",
			Description: "Generate a new password",
			Execute:     executeCreate,
		},
		{
			Name:        "show",
			Description: "Display stored passwords",
			Execute:     executeShow,
		},
		{
			Name:        "help",
			Description: "Show available commands",
			Execute:     executeHelp,
		},
		{
			Name:        "exit",
			Description: "Exit interactive mode",
			Execute:     executeExit,
		},
	}
}

type model struct {
	input            string
	showCommandList  bool
	selectedCmd      int
	output           string
	shouldExit       bool
	cursor           int
}

type executeResult struct {
	output string
	exit   bool
}

func initialModel() model {
	return model{
		input:           "",
		showCommandList: false,
		selectedCmd:     0,
		output:          "",
		shouldExit:      false,
		cursor:          0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.shouldExit = true
			return m, tea.Quit

		case tea.KeyEnter:
			commands := getCommands()
			if m.showCommandList && m.selectedCmd < len(commands) {
				// Execute selected command
				cmd := commands[m.selectedCmd]
				m.showCommandList = false
				m.input = ""
				m.selectedCmd = 0
				return m, cmd.Execute
			} else if strings.TrimSpace(m.input) != "" {
				// Parse and execute typed command
				cmdName := strings.TrimSpace(strings.TrimPrefix(m.input, "/"))
				m.input = ""
				m.showCommandList = false
				m.selectedCmd = 0
				
				for _, cmd := range commands {
					if cmd.Name == cmdName {
						return m, cmd.Execute
					}
				}
				m.output = fmt.Sprintf("Unknown command: %s", cmdName)
			}
			return m, nil

		case tea.KeyUp:
			if m.showCommandList && m.selectedCmd > 0 {
				m.selectedCmd--
			}
			return m, nil

		case tea.KeyDown:
			commands := getCommands()
			if m.showCommandList && m.selectedCmd < len(commands)-1 {
				m.selectedCmd++
			}
			return m, nil

		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				if m.input != "/" {
					m.showCommandList = false
				}
			}
			return m, nil

		case tea.KeyEsc:
			m.showCommandList = false
			m.input = ""
			m.selectedCmd = 0
			return m, nil

		default:
			if msg.Type == tea.KeyRunes {
				m.input += string(msg.Runes)
				if m.input == "/" {
					m.showCommandList = true
					m.selectedCmd = 0
				} else if strings.HasPrefix(m.input, "/") {
					m.showCommandList = true
					// Filter commands based on input
				}
			}
			return m, nil
		}

	case executeResult:
		m.output = msg.output
		if msg.exit {
			m.shouldExit = true
			return m, tea.Quit
		}
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FFFF")).
		MarginBottom(1)
	
	b.WriteString(headerStyle.Render("GenP Interactive Mode"))
	b.WriteString("\n")
	
	// Instructions
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Faint(true)
	
	b.WriteString(instructionStyle.Render("Type / to see available commands or type a command directly"))
	b.WriteString("\n\n")

	// Input field
	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00"))
	
	b.WriteString(inputStyle.Render("> " + m.input + "█"))
	b.WriteString("\n\n")

	// Command list (if showing)
	if m.showCommandList {
		commands := getCommands()
		commandListStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00FFFF")).
			Padding(1)
		
		var cmdList strings.Builder
		cmdList.WriteString(lipgloss.NewStyle().Bold(true).Render("Available Commands:"))
		cmdList.WriteString("\n\n")
		
		for i, cmd := range commands {
			if i == m.selectedCmd {
				selectedStyle := lipgloss.NewStyle().
					Background(lipgloss.Color("#00FFFF")).
					Foreground(lipgloss.Color("#000000")).
					Bold(true)
				cmdList.WriteString(selectedStyle.Render(fmt.Sprintf("  ▶ /%s - %s", cmd.Name, cmd.Description)))
			} else {
				cmdList.WriteString(fmt.Sprintf("    /%s - %s", cmd.Name, cmd.Description))
			}
			cmdList.WriteString("\n")
		}
		
		b.WriteString(commandListStyle.Render(cmdList.String()))
		b.WriteString("\n\n")
		
		// Navigation hint
		hintStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Faint(true)
		b.WriteString(hintStyle.Render("Use ↑↓ to navigate, Enter to select, Esc to cancel"))
		b.WriteString("\n")
	}

	// Output area
	if m.output != "" {
		outputStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#888888")).
			Padding(1).
			MarginTop(1)
		
		b.WriteString(outputStyle.Render(m.output))
		b.WriteString("\n")
	}

	return b.String()
}

// Command execution functions
func executeCreate() tea.Msg {
	// Generate password with default options
	password := internal.GeneratePassword(12, true, true, true)
	output := fmt.Sprintf("Generated Password: %s\n\nDo you want to store this password? (Run 'genp create' command with flags for custom options)", password)
	return executeResult{output: output, exit: false}
}

func executeShow() tea.Msg {
	passwords, err := store.GetAllPasswords()
	if err != nil {
		return executeResult{output: fmt.Sprintf("Error: %v", err), exit: false}
	}

	if len(passwords) == 0 {
		return executeResult{output: "No passwords stored yet.", exit: false}
	}

	// For interactive mode, we'll just show that passwords exist
	// Full decryption would require password input which is complex in TUI
	output := fmt.Sprintf("You have %d stored password(s).\nUse 'genp show' command for full access with master password.", len(passwords))
	return executeResult{output: output, exit: false}
}

func executeHelp() tea.Msg {
	commands := getCommands()
	var help strings.Builder
	help.WriteString("Available Commands:\n\n")
	for _, cmd := range commands {
		help.WriteString(fmt.Sprintf("/%s - %s\n", cmd.Name, cmd.Description))
	}
	help.WriteString("\nFor advanced options, use the CLI commands directly (e.g., 'genp create --help')")
	return executeResult{output: help.String(), exit: false}
}

func executeExit() tea.Msg {
	return executeResult{output: "Goodbye!", exit: true}
}

// interactiveCmd represents the interactive command
var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start interactive mode with command palette",
	Long: `Start an interactive mode where you can use commands with a visual interface.
	
Type / to see all available commands in a palette view, similar to Claude Code.
Navigate with arrow keys and press Enter to execute.`,
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("Starting GenP Interactive Mode...\n")
		color.Yellow("Press Ctrl+C to exit\n\n")
		
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			color.Red("Error running interactive mode: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}
