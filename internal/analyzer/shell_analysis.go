// internal/analyzer/shell_analysis.go
package analyzer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func AnalyzeShells() tea.Msg {
	data := InitShellData()

	// Read shell histories
	shellPaths := map[string]string{
		"bash": "~/.bash_history",
		"zsh":  "~/.zsh_history",
		"fish": "~/.local/share/fish/fish_history",
	}

	for shell, path := range shellPaths {
		expandedPath := expandPath(path)
		if history, err := readHistory(expandedPath); err == nil {
			data.Histories[shell] = history
			analyzeCommands(history, &data)
			data.ShellConfigs[shell] = analyzeShellConfigs(shell)
		}
	}

	// Analyze tool usage separately
	var allEntries []CommandEntry
	for _, history := range data.Histories {
		allEntries = append(allEntries, history...)
	}
	data.Insights.ToolUsage = analyzeToolUsage(allEntries)

	return data
}

func readHistory(path string) ([]CommandEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []CommandEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if cmd := cleanHistoryLine(line); cmd != "" {
			entries = append(entries, CommandEntry{
				Command:    cmd,
				Timestamp:  time.Now(), // For simplicity
				Categories: categorizeCommand(cmd),
			})
		}
	}

	return entries, scanner.Err()
}

func cleanHistoryLine(line string) string {
	parts := strings.Fields(line)
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func categorizeCommand(cmd string) []string {
	categories := []string{}
	patterns := map[string][]string{
		"development": {"git", "docker", "npm", "go", "python"},
		"system":      {"sudo", "systemctl", "ps", "top"},
		"file":        {"ls", "cd", "cp", "mv", "rm"},
	}

	for category, patterns := range patterns {
		for _, pattern := range patterns {
			if strings.HasPrefix(cmd, pattern) {
				categories = append(categories, category)
				break
			}
		}
	}

	return categories
}

func analyzeCommands(entries []CommandEntry, data *ShellData) {
	// Initialize maps for analysis
	langUsage := make(map[string]int)
	toolUsage := make(map[string]int)
	timeOfDay := make(map[int]int)
	commandPatterns := make(map[string]int)

	// Get installed languages
	installedLangs := getInstalledLanguages()

	// Analyze each command
	for _, entry := range entries {
		cmd := entry.Command
		hour := entry.Timestamp.Hour()
		timeOfDay[hour]++

		// Language usage analysis
		for lang := range installedLangs {
			if strings.Contains(cmd, lang) ||
				strings.Contains(cmd, getPackageManager(lang)) {
				langUsage[lang]++
			}
		}

		// Development tool analysis
		tools := []string{"git", "docker", "kubectl", "terraform", "ansible", "make"}
		for _, tool := range tools {
			if strings.HasPrefix(cmd, tool) && checkToolInstalled(tool) {
				toolUsage[tool]++
			}
		}

		// Analyze command patterns
		analyzeCommandPattern(cmd, commandPatterns)
	}

	// Update TechnicalProfile
	techProfile := &data.Insights.TechnicalProfile

	// Calculate primary role based on most used language/tool
	if primaryLang, ok := getMostUsed(langUsage); ok {
		techProfile.PrimaryRole = fmt.Sprintf("%s Developer", strings.Title(primaryLang))
	}

	// Calculate tech stack
	techProfile.TechStack = make([]string, 0)
	for lang := range installedLangs {
		if langUsage[lang] > 0 {
			techProfile.TechStack = append(techProfile.TechStack, lang)
		}
	}

	// Calculate proficiency
	totalCommands := len(entries)
	if totalCommands > 0 {
		for lang, count := range langUsage {
			techProfile.Proficiency[lang] = float64(count) / float64(totalCommands)
		}
		for tool, count := range toolUsage {
			techProfile.Proficiency[tool] = float64(count) / float64(totalCommands)
		}
	}

	// Update WorkPatterns
	patterns := &data.Insights.WorkPatterns
	patterns.PeakHours = getPeakHours(timeOfDay)

	// Calculate productivity metrics based on command complexity and variety
	patterns.Productivity = calculateProductivityMetrics(entries, commandPatterns)
}

// internal/analyzer/shell_analysis.go
func analyzeToolUsage(entries []CommandEntry) ToolUsage {
	toolUsage := ToolUsage{
		Editors:    make(map[string]int),
		Languages:  make(map[string]int),
		BuildTools: make(map[string]int),
	}

	// Get installed languages
	installedLangs := getInstalledLanguages()

	// Analyze each command
	for _, entry := range entries {
		cmd := entry.Command

		// Language usage analysis
		for lang := range installedLangs {
			if strings.Contains(cmd, lang) ||
				strings.Contains(cmd, getPackageManager(lang)) {
				toolUsage.Languages[lang]++
			}
		}

		// Editor usage analysis
		editors := []string{"vim", "nvim", "emacs", "code", "nano"}
		for _, editor := range editors {
			if strings.HasPrefix(cmd, editor) && checkToolInstalled(editor) {
				toolUsage.Editors[editor]++
			}
		}

		// Build tool usage analysis
		buildTools := []string{"make", "maven", "gradle", "npm", "yarn", "pip", "cargo", "composer", "bundler"}
		for _, tool := range buildTools {
			if strings.HasPrefix(cmd, tool) && checkToolInstalled(tool) {
				toolUsage.BuildTools[tool]++
			}
		}
	}

	return toolUsage
}

func getPackageManager(lang string) string {
	managers := map[string]string{
		"python": "pip",
		"node":   "npm",
		"go":     "go get",
		"rust":   "cargo",
		"ruby":   "gem",
		"php":    "composer",
	}
	return managers[lang]
}

func analyzeCommandPattern(cmd string, patterns map[string]int) {
	// Define common command patterns
	patternMap := map[string]*regexp.Regexp{
		"git_workflow": regexp.MustCompile(`git (commit|push|pull|merge)`),
		"build":        regexp.MustCompile(`(make|build|compile)`),
		"deploy":       regexp.MustCompile(`(deploy|kubectl|docker)`),
		"test":         regexp.MustCompile(`test|spec|pytest`),
	}

	for pattern, regex := range patternMap {
		if regex.MatchString(cmd) {
			patterns[pattern]++
		}
	}
}

func getMostUsed(usage map[string]int) (string, bool) {
	var maxKey string
	var maxVal int
	for k, v := range usage {
		if v > maxVal {
			maxKey = k
			maxVal = v
		}
	}
	return maxKey, maxVal > 0
}

func getPeakHours(timeOfDay map[int]int) []int {
	type hourCount struct {
		hour  int
		count int
	}

	var hours []hourCount
	for h, c := range timeOfDay {
		hours = append(hours, hourCount{h, c})
	}

	sort.Slice(hours, func(i, j int) bool {
		return hours[i].count > hours[j].count
	})

	// Return top 3 peak hours
	var peaks []int
	for i := 0; i < len(hours) && i < 3; i++ {
		peaks = append(peaks, hours[i].hour)
	}
	return peaks
}

func calculateProductivityMetrics(entries []CommandEntry, patterns map[string]int) map[string]float64 {
	metrics := make(map[string]float64)
	totalCommands := len(entries)

	if totalCommands == 0 {
		return metrics
	}

	// Command variety score
	uniqueCommands := make(map[string]bool)
	for _, entry := range entries {
		uniqueCommands[entry.Command] = true
	}
	metrics["Command Variety"] = float64(len(uniqueCommands)) / float64(totalCommands)

	// Workflow complexity score
	workflowScore := float64(patterns["git_workflow"]+patterns["build"]+
		patterns["deploy"]+patterns["test"]) / float64(totalCommands)
	metrics["Workflow Complexity"] = workflowScore

	return metrics
}

func checkToolInstalled(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

func getInstalledLanguages() map[string]string {
	languages := map[string]string{
		// Programming Languages
		"python":  "python --version",
		"python3": "python3 --version",
		"node":    "node --version",
		"go":      "go version",
		"java":    "java -version",
		"ruby":    "ruby --version",
		"php":     "php --version",
		"rust":    "rustc --version",
		"perl":    "perl --version",
		"scala":   "scala -version",
		"kotlin":  "kotlin -version",
		"swift":   "swift --version",
		"r":       "R --version",
		"julia":   "julia --version",
		"haskell": "ghc --version",
		"elixir":  "elixir --version",
		"erlang":  "erl -version",
		"clang":   "clang --version",
		"gcc":     "gcc --version",
		"dotnet":  "dotnet --version",
		"lua":     "lua -v",
		"ocaml":   "ocaml -version",
		"dart":    "dart --version",
		"zig":     "zig version",
		"nim":     "nim --version",

		// Build Tools & Package Managers
		"maven":    "mvn --version",
		"gradle":   "gradle --version",
		"npm":      "npm --version",
		"yarn":     "yarn --version",
		"pnpm":     "pnpm --version",
		"pip":      "pip --version",
		"cargo":    "cargo --version",
		"composer": "composer --version",
		"bundler":  "bundle --version",

		// DevOps & Cloud Tools
		"docker":    "docker --version",
		"kubectl":   "kubectl version --client",
		"terraform": "terraform version",
		"ansible":   "ansible --version",
		"vagrant":   "vagrant --version",
		"helm":      "helm version",
		"aws":       "aws --version",
		"gcloud":    "gcloud --version",
		"azure":     "az --version",

		// Version Control
		"git":       "git --version",
		"svn":       "svn --version",
		"mercurial": "hg --version",

		// Databases
		"mysql":   "mysql --version",
		"psql":    "psql --version",
		"mongodb": "mongod --version",
		"redis":   "redis-cli --version",

		// Web Servers & Tools
		"nginx":   "nginx -v",
		"apache2": "apache2 -v",
		"curl":    "curl --version",
		"wget":    "wget --version",

		// Text Editors & IDEs
		"vim":   "vim --version",
		"nvim":  "nvim --version",
		"emacs": "emacs --version",
		"code":  "code --version",

		// Shell & Terminal Tools
		"zsh":  "zsh --version",
		"bash": "bash --version",
		"fish": "fish --version",
		"tmux": "tmux -V",
	}

	installed := make(map[string]string)
	for lang, cmd := range languages {
		if out, err := exec.Command("sh", "-c", cmd).Output(); err == nil {
			installed[lang] = string(out)
		}
	}

	// Sort and keep only top 10 most used
	type usageEntry struct {
		name  string
		count int
	}
	var usageList []usageEntry
	for name := range installed {
		count := 0
		// Count occurrences in command history (you'll need to pass this data somehow)
		// For now, we'll just store all installed ones
		usageList = append(usageList, usageEntry{name, count})
	}

	// Sort by usage count
	sort.Slice(usageList, func(i, j int) bool {
		return usageList[i].count > usageList[j].count
	})

	// Keep only top 10
	result := make(map[string]string)
	for i := 0; i < len(usageList) && i < 10; i++ {
		name := usageList[i].name
		result[name] = installed[name]
	}

	return result
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func analyzeShellConfigs(shell string) ShellConfig {
	configPaths := map[string][]string{
		"bash": {
			"~/.bashrc",
			"~/.bash_profile",
			"~/.bash_aliases",
		},
		"zsh": {
			"~/.zshrc",
			"~/.zsh_plugins",
			"~/.zprofile",
		},
		"fish": {
			"~/.config/fish/config.fish",
			"~/.config/fish/functions",
			"~/.config/fish/conf.d",
		},
	}

	config := ShellConfig{
		ConfigFiles: make(map[string]ConfigInfo),
		Aliases:     make(map[string]string),
		Environment: make(map[string]string),
		Plugins:     make([]PluginInfo, 0),
	}

	// Read and analyze config files
	for _, paths := range configPaths[shell] {
		expandedPath := expandPath(paths)
		if info, err := os.Stat(expandedPath); err == nil {
			content, _ := os.ReadFile(expandedPath)
			config.ConfigFiles[paths] = ConfigInfo{
				Path:     expandedPath,
				Modified: info.ModTime(),
				Content:  string(content),
			}

			// Parse the config file
			parseShellConfig(string(content), &config)
		}
	}

	// Detect plugins based on shell type
	detectPlugins(shell, &config)

	return config
}

func parseShellConfig(content string, config *ShellConfig) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		// Parse aliases
		if strings.HasPrefix(line, "alias ") {
			parts := strings.SplitN(strings.TrimPrefix(line, "alias "), "=", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				value := strings.Trim(strings.TrimSpace(parts[1]), "'\"")
				config.Aliases[name] = value
			}
		}

		// Parse environment variables
		if strings.HasPrefix(line, "export ") {
			parts := strings.SplitN(strings.TrimPrefix(line, "export "), "=", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				value := strings.Trim(strings.TrimSpace(parts[1]), "'\"")
				config.Environment[name] = value
			}
		}
	}
}

func detectPlugins(shell string, config *ShellConfig) {
	switch shell {
	case "zsh":
		detectZshPlugins(config)
	case "fish":
		detectFishPlugins(config)
	case "bash":
		detectBashPlugins(config)
	}
}

func detectZshPlugins(config *ShellConfig) {
	// Check for Oh My Zsh plugins
	omzPath := expandPath("~/.oh-my-zsh")
	if info, err := os.Stat(omzPath); err == nil && info.IsDir() {
		pluginsPath := filepath.Join(omzPath, "plugins")
		if pluginsDir, err := os.ReadDir(pluginsPath); err == nil {
			for _, pluginDir := range pluginsDir {
				if pluginDir.IsDir() {
					config.Plugins = append(config.Plugins, PluginInfo{
						Name:        pluginDir.Name(),
						Source:      filepath.Join(pluginsPath, pluginDir.Name()),
						LastUpdated: info.ModTime(),
					})
				}
			}
		}
	}

	// Check for other plugin managers (Antigen, Zinit, Zplug, etc.)
	pluginManagers := []string{
		"~/.antigen",
		"~/.zinit",
		"~/.zplug",
	}

	for _, manager := range pluginManagers {
		path := expandPath(manager)
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			config.Plugins = append(config.Plugins, PluginInfo{
				Name:        filepath.Base(manager),
				Source:      path,
				LastUpdated: info.ModTime(),
			})
		}
	}
}

func detectFishPlugins(config *ShellConfig) {
	fishPluginPath := expandPath("~/.config/fish/conf.d")
	if files, err := os.ReadDir(fishPluginPath); err == nil {
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".fish") {
				info, _ := file.Info()
				config.Plugins = append(config.Plugins, PluginInfo{
					Name:        strings.TrimSuffix(file.Name(), ".fish"),
					Source:      filepath.Join(fishPluginPath, file.Name()),
					LastUpdated: info.ModTime(),
				})
			}
		}
	}
}

func detectBashPlugins(config *ShellConfig) {
	// Check for common bash plugin managers and extensions
	bashPluginPaths := []string{
		"~/.bash_it",
		"~/.local/share/bash-completion",
	}

	for _, path := range bashPluginPaths {
		expandedPath := expandPath(path)
		if info, err := os.Stat(expandedPath); err == nil && info.IsDir() {
			config.Plugins = append(config.Plugins, PluginInfo{
				Name:        filepath.Base(path),
				Source:      expandedPath,
				LastUpdated: info.ModTime(),
			})
		}
	}
}

func analyzeCommandComplexity(data *ShellData) float64 {
	var totalCommands, complexCommands float64

	for _, history := range data.Histories {
		for _, entry := range history {
			totalCommands++

			// Count pipes and redirections
			if strings.Contains(entry.Command, "|") ||
				strings.Contains(entry.Command, ">") ||
				strings.Contains(entry.Command, "<") {
				complexCommands++
			}

			// Count commands with multiple arguments
			if len(strings.Fields(entry.Command)) > 2 {
				complexCommands += 0.5
			}
		}
	}

	if totalCommands == 0 {
		return 0
	}

	return (complexCommands / totalCommands) * 100
}

func generateRecommendations(data *ShellData) []string {
	recommendations := []string{}

	// Analyze shell configuration
	for shell, config := range data.ShellConfigs {
		if len(config.Aliases) < 5 {
			recommendations = append(recommendations,
				fmt.Sprintf("Consider adding more aliases to your %s configuration to improve productivity", shell))
		}

		if len(config.Plugins) < 3 {
			recommendations = append(recommendations,
				fmt.Sprintf("Explore popular %s plugins to enhance your shell experience", shell))
		}
	}

	return recommendations
}

func generateWorkflowTips(data *ShellData) []string {
	tips := []string{}

	// Analyze command patterns
	commonPatterns := analyzeCommandPatterns(data)
	for pattern, count := range commonPatterns {
		if count > 10 {
			tips = append(tips, fmt.Sprintf(
				"You frequently use '%s'. Consider creating an alias for this pattern", pattern))
		}
	}

	return tips
}

func analyzeCommandPatterns(data *ShellData) map[string]int {
	patterns := make(map[string]int)

	for _, history := range data.Histories {
		for _, entry := range history {
			// Look for common command sequences
			parts := strings.Fields(entry.Command)
			if len(parts) > 1 {
				pattern := strings.Join(parts[:2], " ")
				patterns[pattern]++
			}
		}
	}

	return patterns
}
