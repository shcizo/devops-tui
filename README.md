# devops-tui

A terminal user interface (TUI) for Azure DevOps Boards, inspired by [jira-cli](https://github.com/ankitpokhrel/jira-cli) and [JiraTUI](https://jiratui.sh/).

## Features

- View Azure DevOps work items in a clean terminal interface
- Filter by Sprint, State, and Assigned To
- Vim-style navigation (j/k/g/G)
- Fullscreen detail view
- Open work items in browser
- Cross-platform (Windows, macOS, Linux)

## Installation

### Build from source

```bash
go build -o devops-tui .
```

### Move to PATH

```bash
mv devops-tui /usr/local/bin/
```

## Configuration

Create a config file at `~/.config/devops-tui/config.yaml`:

```yaml
# Azure DevOps connection
organization: "my-organization"
project: "my-project"
team: "my-team"

# Authentication (or use AZURE_DEVOPS_PAT env variable)
pat: "your-personal-access-token"

# UI settings
theme: "default"

# Default filters at startup
defaults:
  sprint: "current"
  state: "all"
  assigned: "me"
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `AZURE_DEVOPS_PAT` | Personal Access Token (recommended) |
| `AZURE_DEVOPS_ORG` | Organization (overrides config) |
| `AZURE_DEVOPS_PROJECT` | Project (overrides config) |
| `AZURE_DEVOPS_TEAM` | Team (overrides config) |

### PAT Permissions

Your Personal Access Token needs these scopes:
- `Work Items (Read)` - Read work items
- `Project and Team (Read)` - List sprints/iterations

## Keyboard Shortcuts

### Global

| Key | Description |
|-----|-------------|
| `Tab` | Switch to next panel |
| `Shift+Tab` | Switch to previous panel |
| `?` | Show/hide help |
| `Ctrl+r` | Reload data |
| `q` / `Ctrl+c` | Quit |

### Navigation

| Key | Description |
|-----|-------------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `g` | Go to first item |
| `G` | Go to last item |

### Actions

| Key | Description |
|-----|-------------|
| `Enter` / `Space` | Select filter / Open in browser |
| `v` | View fullscreen details |

### Detail View

| Key | Description |
|-----|-------------|
| `Esc` / `q` | Back to main view |
| `Enter` | Open in browser |
| `j` / `k` | Scroll description |

## Tech Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Viper](https://github.com/spf13/viper) - Configuration

## License

MIT
