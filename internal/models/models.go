// internal/models/models.go
package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ksauraj/k8au-shell-analyzer/internal/analyzer"
	"github.com/ksauraj/k8au-shell-analyzer/internal/gemini"
	"github.com/ksauraj/k8au-shell-analyzer/internal/render"
	"github.com/ksauraj/k8au-shell-analyzer/internal/types"
)

type Model struct {
	viewport              viewport.Model
	loading               bool
	err                   error
	shellData             analyzer.ShellData
	currentView           string
	tabs                  []string
	activeTab             int
	logger                *log.Logger
	sections              []gemini.Section
	currentSectionIndex   int
	currentAnimationFrame int
	animationTicker       *time.Ticker
	sectionSwitchTicker   *time.Ticker
	timelineData          []types.TimelineEntry
}

func InitialModel() Model {
	logFile, err := os.OpenFile("shell_analyzer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	tabs := []string{"Overview", "Tech Profile", "Work Patterns", "Tool Usage", "Wrapped", "Timeline"}

	animationTicker := time.NewTicker(500 * time.Millisecond)
	sectionSwitchTicker := time.NewTicker(10 * time.Second)

	return Model{
		viewport:            viewport.New(80, 24),
		loading:             true,
		currentView:         "main",
		tabs:                tabs,
		activeTab:           0,
		logger:              logger,
		animationTicker:     animationTicker,
		sectionSwitchTicker: sectionSwitchTicker,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		analyzer.AnalyzeShells,
		tea.EnterAltScreen,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
			return m, nil
		case "right", "l", "n":
			if len(m.sections) > 0 {
				m.currentSectionIndex = (m.currentSectionIndex + 1) % len(m.sections)
			}
			return m, nil
		case "left", "h", "p":
			if len(m.sections) > 0 {
				m.currentSectionIndex--
				if m.currentSectionIndex < 0 {
					m.currentSectionIndex = len(m.sections) - 1
				}
			}
			return m, nil
		}

	case analyzer.ShellData:
		m.loading = false
		m.shellData = msg
		m.timelineData = analyzer.GenerateTimelineData(msg)

		wrappedResp, err := gemini.GenerateWrapped(analyzer.ShellDataToString(msg))
		if err != nil {
			m.err = err
			m.logger.Printf("Error generating wrapped response: %v", err)
			return m, nil
		}

		// Debug log
		m.logger.Printf("Generated %d sections", len(wrappedResp.Sections))

		// Remove animation data and store sections
		m.sections = make([]gemini.Section, len(wrappedResp.Sections))
		for i := range wrappedResp.Sections {
			m.sections[i] = wrappedResp.Sections[i]
			m.sections[i].Animation = nil
		}

		m.currentSectionIndex = 0

		// Debug log
		m.logger.Printf("Stored %d sections, starting at index %d",
			len(m.sections), m.currentSectionIndex)

		// Start the section switch ticker if we have sections
		if len(m.sections) > 0 {
			m.sectionSwitchTicker = time.NewTicker(10 * time.Second)
		}

		return m, nil

	case time.Time:
		if len(m.sections) > 0 {
			switch msg {
			case <-m.sectionSwitchTicker.C:
				m.currentSectionIndex = (m.currentSectionIndex + 1) % len(m.sections)
				return m, nil
			}
		}
		m.viewport, _ = m.viewport.Update(msg)
		return m, nil

	default:
		m.viewport, _ = m.viewport.Update(msg)
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	if m.loading {
		return render.RenderLoading()
	}

	// Header with title and version
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Padding(0, 1).
		Render("🚀 K8au Shell Analyzer v1.0.1-beta")

	// Render tabs
	tabBar := render.RenderTabs(m.tabs, m.activeTab)

	// Content (existing switch case)
	var content string
	switch m.tabs[m.activeTab] {
	case "Overview":
		content = render.RenderOverview(m.shellData)
	case "Tech Profile":
		content = render.RenderTechProfile(m.shellData.Insights.TechnicalProfile)
	case "Work Patterns":
		content = render.RenderWorkPatterns(m.shellData.Insights.WorkPatterns)
	case "Tool Usage":
		content = render.RenderToolUsage(m.shellData.Insights.ToolUsage)
	case "Timeline":
		content = render.RenderTimeline(m.timelineData)
	case "Wrapped":
		if len(m.sections) == 0 {
			content = lipgloss.NewStyle().
				Width(50).
				BorderStyle(lipgloss.RoundedBorder()).
				Padding(1).
				Render("Generating wrapped view...")
		} else {
			currentSection := m.sections[m.currentSectionIndex]
			content = lipgloss.NewStyle().
				Width(50).
				BorderStyle(lipgloss.RoundedBorder()).
				Padding(1).
				Render(fmt.Sprintf(
					"📺 Slide %d/%d\n\n%s\n\n%s\n\n%s",
					m.currentSectionIndex+1,
					len(m.sections),
					lipgloss.NewStyle().Bold(true).Render(currentSection.Title),
					lipgloss.NewStyle().Width(48).Render(currentSection.Description),
					render.RenderQuotes(currentSection.Quotes),
				))
		}
	}
	// Footer with controls
	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(0, 1).
		Render("↑/↓: Navigate • Tab: Switch Views • q: Quit • Left/Right: Change Slides • By Ksauraj")

	// Join all components vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"\n",
		tabBar,
		"\n",
		content,
		"\n",
		footer,
	)
}

func (m Model) Cleanup() {
	m.animationTicker.Stop()
	m.sectionSwitchTicker.Stop()
	tea.ExitAltScreen()
}
