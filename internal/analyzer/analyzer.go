// internal/analyzer/analyzer.go
package analyzer

import (
	"fmt"
	"strings"
	"time"

	"github.com/ksauraj/k8au-shell-analyzer/internal/types"
)

// ShellData contains all the analyzed shell data
type ShellData struct {
	Histories    map[string][]CommandEntry
	CommonCmds   map[string]int
	TimePatterns map[string]int
	Insights     DetailedInsights
	ShellConfigs map[string]ShellConfig
}

// CommandEntry represents a single command entry in the shell history
type CommandEntry struct {
	Command    string
	Timestamp  time.Time
	Count      int
	Categories []string
}

// DetailedInsights contains detailed insights about the user's shell usage
type DetailedInsights struct {
	TechnicalProfile TechProfile
	WorkPatterns     WorkPatterns
	ToolUsage        ToolUsage
}

// TechProfile contains technical profile information
type TechProfile struct {
	PrimaryRole     string
	SecondarySkills []string
	TechStack       []string
	Proficiency     map[string]float64
}

// WorkPatterns contains work pattern information
type WorkPatterns struct {
	PeakHours       []int
	CommonWorkflows []string
	Productivity    map[string]float64
}

// ToolUsage contains tool usage statistics
type ToolUsage struct {
	Editors    map[string]int
	Languages  map[string]int
	BuildTools map[string]int
}

// ShellConfig contains shell configuration information
type ShellConfig struct {
	ConfigFiles map[string]ConfigInfo
	Plugins     []PluginInfo
	Aliases     map[string]string
	Environment map[string]string
}

// ConfigInfo contains information about a configuration file
type ConfigInfo struct {
	Path     string
	Modified time.Time
	Content  string
}

// PluginInfo contains information about a plugin
type PluginInfo struct {
	Name        string
	Source      string
	LastUpdated time.Time
}

// InitShellData initializes an empty ShellData structure
func InitShellData() ShellData {
	return ShellData{
		Histories:    make(map[string][]CommandEntry),
		CommonCmds:   make(map[string]int),
		TimePatterns: make(map[string]int),
		Insights: DetailedInsights{
			TechnicalProfile: TechProfile{
				Proficiency: make(map[string]float64),
			},
			WorkPatterns: WorkPatterns{
				Productivity: make(map[string]float64),
			},
			ToolUsage: ToolUsage{
				Editors:    make(map[string]int),
				Languages:  make(map[string]int),
				BuildTools: make(map[string]int),
			},
		},
		ShellConfigs: make(map[string]ShellConfig),
	}
}

// ShellDataToString converts ShellData into a concise string representation
func ShellDataToString(data ShellData) string {
	var result strings.Builder

	// Add shell usage summary
	for shell, history := range data.Histories {
		result.WriteString(fmt.Sprintf("Shell: %s, Commands: %d\n", shell, len(history)))
	}

	// Add tech stack
	if len(data.Insights.TechnicalProfile.TechStack) > 0 {
		result.WriteString("Tech Stack: " + strings.Join(data.Insights.TechnicalProfile.TechStack, ", ") + "\n")
	}

	// Add peak hours
	if len(data.Insights.WorkPatterns.PeakHours) > 0 {
		result.WriteString("Peak Hours: ")
		for _, hour := range data.Insights.WorkPatterns.PeakHours {
			result.WriteString(fmt.Sprintf("%02d:00 ", hour))
		}
		result.WriteString("\n")
	}

	// Add productivity metrics
	if len(data.Insights.WorkPatterns.Productivity) > 0 {
		result.WriteString("Productivity Metrics:\n")
		for metric, value := range data.Insights.WorkPatterns.Productivity {
			result.WriteString(fmt.Sprintf("- %s: %.1f%%\n", metric, value*100))
		}
	}

	// Add tool usage
	if len(data.Insights.ToolUsage.Editors) > 0 {
		result.WriteString("Editors:\n")
		for editor, count := range data.Insights.ToolUsage.Editors {
			result.WriteString(fmt.Sprintf("- %s: %d uses\n", editor, count))
		}
	}

	return result.String()
}

func GenerateTimelineData(data ShellData) []types.TimelineEntry {
	var timelineData []types.TimelineEntry

	// Track unique commands to avoid duplicates
	uniqueCommands := make(map[string]bool)

	// Iterate through shell histories
	for shell, history := range data.Histories {
		for _, entry := range history {
			// Skip if we already have this command
			if uniqueCommands[entry.Command] {
				continue
			}

			// Add interesting commands to the timeline
			if isInterestingCommand(entry.Command) {
				timelineData = append(timelineData, types.TimelineEntry{
					Timestamp: entry.Timestamp,
					Command:   entry.Command,
					Shell:     shell,
				})
				uniqueCommands[entry.Command] = true
			}

			// Stop after collecting 15 commands
			if len(timelineData) >= 15 {
				return timelineData
			}
		}
	}

	return timelineData
}

// isInterestingCommand checks if a command is worth showing in the timeline
func isInterestingCommand(command string) bool {
	// List of interesting commands
	interestingCommands := []string{"git", "docker", "kubectl", "terraform", "ansible", "make", "npm", "go", "python", "java", "ssh", "scp", "curl", "wget", "vim", "nvim", "emacs", "code"}

	// Check if the command contains special characters
	hasSpecialChars := strings.ContainsAny(command, "|><&;")

	// Check if the command is a typo
	isTypo := isTypoCommand(command)

	// Check if the command is in the interesting list or has special characters or is a typo
	for _, interesting := range interestingCommands {
		if strings.HasPrefix(command, interesting) {
			return true
		}
	}

	return hasSpecialChars || isTypo
}

// isTypoCommand checks if a command is a common typo
func isTypoCommand(command string) bool {
	commonTypos := []string{"sl", "cd..", "pythoon", "gti", "vmi", "nivm", "emasc", "clea", "exot"}
	for _, typo := range commonTypos {
		if command == typo {
			return true
		}
	}
	return false
}
