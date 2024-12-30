// internal/render/render.go
package render

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gookit/color"
	"github.com/ksauraj/k8au-shell-analyzer/internal/analyzer"
	"github.com/ksauraj/k8au-shell-analyzer/internal/types"
)

type WrappedResponse struct {
	Text string
}

// RenderLoading renders the loading screen
func RenderLoading() string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Render("Analyzing your shell history... 🔍")
}

// RenderTabs renders the tab bar
func RenderTabs(tabs []string, active int) string {
	var tabsDisplay strings.Builder

	for i, tab := range tabs {
		style := lipgloss.NewStyle().
			Padding(0, 2)

		if i == active {
			style = style.
				Bold(true).
				Background(lipgloss.Color("4")).
				Foreground(lipgloss.Color("15"))
		}

		tabsDisplay.WriteString(style.Render(tab))
	}

	return tabsDisplay.String()
}

func RenderOverview(data analyzer.ShellData) string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1)

	var content strings.Builder
	content.WriteString(color.Green.Sprintf("📊 Shell Usage Overview\n\n"))

	for shell, history := range data.Histories {
		content.WriteString(fmt.Sprintf("Shell: %s\n", color.Cyan.Sprint(shell)))
		content.WriteString(fmt.Sprintf("Commands: %d\n", len(history)))

		// Add shell configuration information
		if config, exists := data.ShellConfigs[shell]; exists {
			content.WriteString("\nConfiguration:\n")
			content.WriteString(fmt.Sprintf("• Aliases: %d\n", len(config.Aliases)))
			content.WriteString(fmt.Sprintf("• Plugins: %d\n", len(config.Plugins)))
			content.WriteString(fmt.Sprintf("• Environment Variables: %d\n", len(config.Environment)))

			// List up to 3 plugins
			if len(config.Plugins) > 0 {
				content.WriteString("\nInstalled Plugins:\n")
				for i, plugin := range config.Plugins {
					if i >= 3 { // Show only the first 3 plugins
						break
					}
					content.WriteString(fmt.Sprintf("• %s (from %s)\n",
						color.Yellow.Sprint(plugin.Name),
						plugin.Source))
				}
				if len(config.Plugins) > 3 {
					content.WriteString(fmt.Sprintf("• And %d more...\n", len(config.Plugins)-3))
				}
			}

			// List some aliases if any
			if len(config.Aliases) > 0 {
				content.WriteString("\nSome Aliases:\n")
				count := 0
				for alias, command := range config.Aliases {
					if count >= 5 { // Show only first 5 aliases
						break
					}
					content.WriteString(fmt.Sprintf("• %s → %s\n",
						color.Yellow.Sprint(alias),
						command))
					count++
				}
			}
		}
		content.WriteString("\n")
	}

	return style.Render(content.String())
}

// RenderTechProfile renders the tech profile tab
func RenderTechProfile(profile analyzer.TechProfile) string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1)

	var content strings.Builder
	content.WriteString(color.Green.Sprintf("💻 Technical Profile\n\n"))

	// Primary Role
	if profile.PrimaryRole != "" {
		content.WriteString(fmt.Sprintf("🎯 Primary Role: %s\n\n",
			color.Cyan.Sprint(profile.PrimaryRole)))
	} else {
		content.WriteString("🎯 Primary Role: Not enough data\n\n")
	}

	// Tech Stack
	content.WriteString("💻 Tech Stack:\n")
	if len(profile.TechStack) > 0 {
		for _, tech := range profile.TechStack {
			content.WriteString(fmt.Sprintf("• %s\n", tech))
		}
	} else {
		content.WriteString("No tech stack data available\n")
	}
	content.WriteString("\n")

	// Secondary Skills
	content.WriteString("🛠️  Secondary Skills:\n")
	if len(profile.SecondarySkills) > 0 {
		for _, skill := range profile.SecondarySkills {
			content.WriteString(fmt.Sprintf("• %s\n", skill))
		}
	} else {
		content.WriteString("No secondary skills data available\n")
	}
	content.WriteString("\n")

	// Proficiency Levels
	content.WriteString("📊 Proficiency Levels:\n")
	if len(profile.Proficiency) > 0 {
		// Sort proficiencies for consistent display
		var items []struct {
			Name  string
			Level float64
		}
		for tech, level := range profile.Proficiency {
			items = append(items, struct {
				Name  string
				Level float64
			}{tech, level})
		}
		// Sort by proficiency level in descending order
		sort.Slice(items, func(i, j int) bool {
			return items[i].Level > items[j].Level
		})

		for _, item := range items {
			bars := int(item.Level * 20)
			if bars < 0 {
				bars = 0
			}
			barStr := strings.Repeat("█", bars) + strings.Repeat("░", 20-bars)
			content.WriteString(fmt.Sprintf("%-15s %s %.1f%%\n",
				item.Name, barStr, item.Level*100))
		}
	} else {
		content.WriteString("No proficiency data available\n")
	}

	return style.Render(content.String())
}

// RenderWorkPatterns renders the work patterns tab
func RenderWorkPatterns(patterns analyzer.WorkPatterns) string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1)

	var content strings.Builder
	content.WriteString(color.Yellow.Sprintf("⏰ Work Patterns\n\n"))

	// Daily Activity
	content.WriteString("📅 Daily Activity:\n")
	for _, hour := range patterns.PeakHours {
		content.WriteString(fmt.Sprintf("Peak hour: %02d:00\n", hour))
	}
	content.WriteString("\n")

	// Productivity Metrics
	content.WriteString("📈 Productivity Metrics:\n")
	for metric, value := range patterns.Productivity {
		bars := int(value * 20)
		barStr := strings.Repeat("█", bars) + strings.Repeat("░", 20-bars)
		content.WriteString(fmt.Sprintf("%-20s %s %.1f%%\n", metric, barStr, value*100))
	}
	content.WriteString("\n")

	// Common Workflows
	content.WriteString("🔄 Common Workflows:\n")
	for _, workflow := range patterns.CommonWorkflows {
		content.WriteString(fmt.Sprintf("• %s\n", workflow))
	}

	return style.Render(content.String())
}

func RenderToolUsage(usage analyzer.ToolUsage) string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1)

	var content strings.Builder
	content.WriteString(color.Magenta.Sprintf("🔧 Tool Usage Statistics\n\n"))

	// Editors Section
	content.WriteString("📝 Editors:\n")
	if len(usage.Editors) > 0 {
		for editor, count := range usage.Editors {
			content.WriteString(fmt.Sprintf("• %s: %d uses\n", editor, count))
		}
	} else {
		content.WriteString("No editor usage data available\n")
	}
	content.WriteString("\n")

	// Languages Section
	content.WriteString("💻 Programming Languages:\n")
	if len(usage.Languages) > 0 {
		for lang, count := range usage.Languages {
			content.WriteString(fmt.Sprintf("• %s: %d uses\n", lang, count))
		}
	} else {
		content.WriteString("No language usage data available\n")
	}
	content.WriteString("\n")

	// Build Tools Section
	content.WriteString("🛠️  Build Tools:\n")
	if len(usage.BuildTools) > 0 {
		for tool, count := range usage.BuildTools {
			content.WriteString(fmt.Sprintf("• %s: %d uses\n", tool, count))
		}
	} else {
		content.WriteString("No build tool usage data available\n")
	}

	return style.Render(content.String())
}

func RenderWrapped(content string) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1)
	return style.Render(content)
}

// removeMarkdownPlaceholders removes markdown placeholders from the text
func removeMarkdownPlaceholders(text string) string {
	// Remove (Text animation: ...) placeholders
	for strings.Contains(text, "(Text animation:") {
		start := strings.Index(text, "(Text animation:")
		end := strings.Index(text[start:], ")") + start + 1
		text = text[:start] + text[end:]
	}

	// Remove **bold** markdown
	text = strings.ReplaceAll(text, "**", "")

	// Remove *italic* markdown
	text = strings.ReplaceAll(text, "*", "")

	return text
}

func RenderTimeline(entries []types.TimelineEntry) string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1)

	var content strings.Builder
	content.WriteString(color.Green.Sprintf("⏳ Interesting Commands Timeline\n\n"))

	for _, entry := range entries {
		content.WriteString(fmt.Sprintf("📅 %s - %s (%s)\n",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			color.Cyan.Sprint(entry.Command),
			color.Yellow.Sprint(entry.Shell)))
	}

	return style.Render(content.String())
}

func RenderQuotes(quotes []string) string {
	var content strings.Builder

	// Add a header for the quotes section
	content.WriteString(color.Green.Sprintf("📜 Quotes\n\n"))

	// Render each quote
	for _, quote := range quotes {
		// Remove any unwanted markdown or formatting
		quote = strings.ReplaceAll(quote, "**", "") // Remove bold markdown
		quote = strings.ReplaceAll(quote, "*", "")  // Remove italic markdown

		// Add the quote with proper indentation
		content.WriteString(fmt.Sprintf("• \"%s\"\n", quote))
	}

	return content.String()
}
