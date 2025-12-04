# devops-tui - MVP Specifikation

> En terminal user interface (TUI) för Azure DevOps Boards, inspirerad av [jira-cli](https://github.com/ankitpokhrel/jira-cli) och [JiraTUI](https://jiratui.sh/).

## Översikt

**devops-tui** är ett terminalbaserat verktyg för att navigera och granska Azure DevOps work items direkt från kommandoraden. MVP:n fokuserar på read-only boards-funktionalitet med en panel-baserad layout inspirerad av Lazygit.

### Mål

- Snabb överblick av work items utan att lämna terminalen
- Vim-inspirerad navigation för effektivt arbetsflöde
- Enkel installation via single binary (Go)

### Avgränsningar (MVP)

| Inkluderat | Exkluderat |
|------------|------------|
| Boards (work items) | Pipelines/Builds |
| Read-only visning | Skapa/redigera items |
| Fasta filter (Sprint/State/Assigned) | WIQL-queries |
| Ett konfigurerat projekt | Projekt-switcher |
| Personal Access Token auth | OAuth/SSO |

---

## Tech Stack

```
┌─────────────────────────────────────────┐
│              devops-tui                 │
├─────────────────────────────────────────┤
│  UI Layer                               │
│  ├─ Bubble Tea (TUI framework)          │
│  ├─ Bubbles (komponenter)               │
│  └─ Lip Gloss (styling)                 │
├─────────────────────────────────────────┤
│  Data Layer                             │
│  ├─ Azure DevOps REST API v7.1          │
│  └─ HTTP client (net/http)              │
├─────────────────────────────────────────┤
│  Config                                 │
│  ├─ Viper (konfiguration)               │
│  └─ YAML config file                    │
└─────────────────────────────────────────┘
```

### Varför Go + Bubble Tea?

1. **Single binary** - Enkel distribution, inga runtime-beroenden
2. **Cross-platform** - Windows, macOS, Linux utan extra arbete
3. **Beprövat** - Används av GitHub CLI, lazygit, och jira-cli
4. **Bra ekosystem** - Charm.sh har kompletterande bibliotek

### Dependencies

```go
require (
    github.com/charmbracelet/bubbletea  // TUI framework
    github.com/charmbracelet/bubbles    // UI komponenter
    github.com/charmbracelet/lipgloss   // Styling
    github.com/spf13/viper              // Konfiguration
)
```

---

## UI Design

### Layout (3-panel)

```
┌─ devops-tui ──────────────────────────────────────────────────┐
│                                                               │
│ ┌─ FILTER ──────────┐ ┌─ WORK ITEMS ────────────────────────┐ │
│ │                   │ │                                     │ │
│ │ Sprint            │ │  ID     TYPE   STATE    TITLE       │ │
│ │ ─────────────     │ │  ─────────────────────────────────  │ │
│ │ ▸ Sprint 42       │ │  #1234  Story  Active   Implement   │ │
│ │   Sprint 41       │ │▸ #1235  Task   Active   Create lo   │ │
│ │   Sprint 40       │ │  #1236  Task   New      Add JWT m   │ │
│ │   Backlog         │ │  #1237  Bug    Active   Fix token   │ │
│ │                   │ │  #1238  Task   Resolved Setup CI    │ │
│ ├───────────────────┤ │                                     │ │
│ │ State             │ │                                     │ │
│ │ ─────────────     │ │                                     │ │
│ │ ● All             │ │                                     │ │
│ │ ○ New             │ └─────────────────────────────────────┘ │
│ │ ○ Active          │ ┌─ DETAILS ───────────────────────────┐ │
│ │ ○ Resolved        │ │                                     │ │
│ │ ○ Closed          │ │ #1235 Create login component        │ │
│ ├───────────────────┤ │                                     │ │
│ │ Assigned          │ │ Type: Task         State: Active    │ │
│ │ ─────────────     │ │ Assigned: Samuel   Sprint: 42       │ │
│ │ ● All             │ │ Area: Frontend     Priority: 2      │ │
│ │ ○ Me              │ │                                     │ │
│ │                   │ │ Parent: #1234 Implement auth flow   │ │
│ └───────────────────┘ │                                     │ │
│                       │ ─── Description ───                 │ │
│                       │ Create a React login component      │ │
│                       │ with email/password fields...       │ │
│                       │                                     │ │
│                       │ ─── Tags ───                        │ │
│                       │ frontend · react · auth             │ │
│                       │                                     │ │
│                       └─────────────────────────────────────┘ │
│                                                               │
├───────────────────────────────────────────────────────────────┤
│ Tab Panel  j/k Navigate  g/G Top/Bottom  Enter Open  ? Help   │
└───────────────────────────────────────────────────────────────┘
```

### Paneler

| Panel | Bredd | Innehåll | Interaktion |
|-------|-------|----------|-------------|
| **Filter** | ~20% | Sprint, State, Assigned filter | Välj filter med Enter/Space |
| **Work Items** | ~80% (topp) | Tabell med work items | Navigera med j/k, öppna med Enter |
| **Details** | ~80% (botten) | Detaljer för valt item | Automatiskt uppdaterad |

### Fullskärms-detaljvy

När användaren trycker `v` på ett work item:

```
┌─ #1235 Create login component ────────────────────────────────┐
│                                                               │
│  ┌─ METADATA ───────────────────────────────────────────────┐ │
│  │                                                          │ │
│  │  Type:      Task              ID:        #1235           │ │
│  │  State:     Active            Created:   2024-01-15      │ │
│  │  Assigned:  Samuel Enocsson   Updated:   2024-01-18      │ │
│  │  Sprint:    Sprint 42         Priority:  2               │ │
│  │  Area:      Frontend                                     │ │
│  │                                                          │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                               │
│  ┌─ PARENT ─────────────────────────────────────────────────┐ │
│  │  #1234 User Story: Implement authentication flow         │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                               │
│  ┌─ DESCRIPTION ────────────────────────────────────────────┐ │
│  │                                                          │ │
│  │  Create a React login component with the following:      │ │
│  │                                                          │ │
│  │  - Email input field with validation                     │ │
│  │  - Password input field                                  │ │
│  │  - "Remember me" checkbox                                │ │
│  │  - Submit button with loading state                      │ │
│  │  - Error message display                                 │ │
│  │                                                          │ │
│  │  The component should follow our design system and       │ │
│  │  integrate with the existing auth context.               │ │
│  │                                                          │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                               │
│  ┌─ TAGS ───────────────────────────────────────────────────┐ │
│  │  frontend · react · auth · ui                            │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                               │
├───────────────────────────────────────────────────────────────┤
│ Esc Back  Enter Open in browser  j/k Scroll                   │
└───────────────────────────────────────────────────────────────┘
```

---

## Keyboard Shortcuts

### Globala

| Key | Beskrivning |
|-----|-------------|
| `Tab` | Byt till nästa panel |
| `Shift+Tab` | Byt till föregående panel |
| `?` | Visa/dölj hjälp |
| `Ctrl+r` | Ladda om data |
| `q` / `Ctrl+c` | Avsluta |

### Filter-panel

| Key | Beskrivning |
|-----|-------------|
| `j` / `↓` | Nästa filter/val |
| `k` / `↑` | Föregående filter/val |
| `Enter` / `Space` | Välj filter |
| `g` | Gå till första |
| `G` | Gå till sista |

### Work Items-panel

| Key | Beskrivning |
|-----|-------------|
| `j` / `↓` | Nästa work item |
| `k` / `↑` | Föregående work item |
| `g` | Gå till första |
| `G` | Gå till sista |
| `Enter` | Öppna i webbläsare |
| `v` | Visa fullskärms-detaljer |
| `/` | Sök (filter by text) |

### Detaljvy (fullskärm)

| Key | Beskrivning |
|-----|-------------|
| `Esc` / `q` | Tillbaka till huvudvy |
| `Enter` | Öppna i webbläsare |
| `j` / `k` | Scrolla beskrivning |

---

## Konfiguration

### Config-fil

Placering: `~/.config/devops-tui/config.yaml`

```yaml
# Azure DevOps-anslutning
organization: "my-organization"
project: "my-project"

# Autentisering
# PAT kan anges här eller via miljövariabel AZURE_DEVOPS_PAT
pat: ""

# UI-inställningar
theme: "default"  # default, dark, light

# Standardfilter vid uppstart
defaults:
  sprint: "current"      # "current", "all", eller specifikt namn
  state: "all"           # "all", "new", "active", "resolved", "closed"
  assigned: "me"         # "all", "me"
```

### Miljövariabler

| Variabel | Beskrivning |
|----------|-------------|
| `AZURE_DEVOPS_PAT` | Personal Access Token (rekommenderat) |
| `AZURE_DEVOPS_ORG` | Organisation (override config) |
| `AZURE_DEVOPS_PROJECT` | Projekt (override config) |

### PAT-behörigheter

Personal Access Token behöver följande scope:
- `Work Items (Read)` - Läsa work items
- `Project and Team (Read)` - Lista sprints/iterationer

---

## Azure DevOps API

### Endpoints som används

```
Base URL: https://dev.azure.com/{organization}/{project}/_apis

Work Items:
  GET /wit/wiql                    # Kör WIQL-query
  GET /wit/workitems/{id}          # Hämta enskilt work item
  GET /wit/workitems?ids={ids}     # Hämta flera work items

Iterationer (Sprints):
  GET /work/teamsettings/iterations  # Lista iterationer

Team:
  GET /teams                       # Lista teams
```

### WIQL-queries

För att hämta work items använder vi WIQL (Work Item Query Language):

```sql
-- Alla items i current sprint, tilldelade mig
SELECT [System.Id], [System.Title], [System.State], [System.WorkItemType]
FROM WorkItems
WHERE [System.TeamProject] = @project
  AND [System.IterationPath] = @currentIteration
  AND [System.AssignedTo] = @me
ORDER BY [System.ChangedDate] DESC

-- Alla aktiva items
SELECT [System.Id], [System.Title], [System.State], [System.WorkItemType]
FROM WorkItems
WHERE [System.TeamProject] = @project
  AND [System.State] = 'Active'
ORDER BY [System.ChangedDate] DESC
```

### Response-struktur (Work Item)

```json
{
  "id": 1235,
  "rev": 5,
  "fields": {
    "System.Id": 1235,
    "System.Title": "Create login component",
    "System.State": "Active",
    "System.WorkItemType": "Task",
    "System.AssignedTo": {
      "displayName": "Samuel Enocsson",
      "uniqueName": "samuel@example.com"
    },
    "System.IterationPath": "MyProject\\Sprint 42",
    "System.AreaPath": "MyProject\\Frontend",
    "System.Description": "<div>Create a React login...</div>",
    "System.Tags": "frontend; react; auth",
    "System.Parent": 1234,
    "Microsoft.VSTS.Common.Priority": 2,
    "System.CreatedDate": "2024-01-15T10:00:00Z",
    "System.ChangedDate": "2024-01-18T14:30:00Z"
  },
  "url": "https://dev.azure.com/org/project/_apis/wit/workItems/1235"
}
```

---

## Projektstruktur

```
devops-tui/
├── main.go                 # Entry point
├── go.mod
├── go.sum
├── README.md
├── SPEC.md                 # Denna fil
│
├── cmd/
│   └── root.go             # CLI setup (cobra om vi vill ha subcommands)
│
├── internal/
│   ├── config/
│   │   └── config.go       # Viper config loading
│   │
│   ├── api/
│   │   ├── client.go       # Azure DevOps HTTP client
│   │   ├── workitems.go    # Work item queries
│   │   └── iterations.go   # Sprint/iteration queries
│   │
│   ├── ui/
│   │   ├── app.go          # Bubble Tea main model
│   │   ├── styles.go       # Lip Gloss styles
│   │   ├── keys.go         # Keyboard bindings
│   │   │
│   │   ├── components/
│   │   │   ├── filter.go       # Filter panel
│   │   │   ├── workitems.go    # Work items list
│   │   │   ├── details.go      # Details panel
│   │   │   ├── detailview.go   # Fullscreen detail view
│   │   │   └── help.go         # Help overlay
│   │   │
│   │   └── views/
│   │       ├── main.go         # Main 3-panel view
│   │       └── detail.go       # Fullscreen detail view
│   │
│   └── models/
│       ├── workitem.go     # Work item domain model
│       ├── iteration.go    # Sprint/iteration model
│       └── filter.go       # Filter state model
│
└── pkg/
    └── browser/
        └── open.go         # Cross-platform browser open
```

---

## Implementation - Milstolpar

### Fas 1: Grundstruktur
- [ ] Initiera Go-modul
- [ ] Sätt upp projektstruktur
- [ ] Konfigurera Bubble Tea scaffold
- [ ] Implementera config-laddning (Viper)

### Fas 2: Azure DevOps API
- [ ] HTTP-klient med PAT-auth
- [ ] Hämta iterationer (sprints)
- [ ] Kör WIQL-queries
- [ ] Hämta work items med detaljer

### Fas 3: UI - Paneler
- [ ] Filter-panel (sprint, state, assigned)
- [ ] Work items-lista med tabell
- [ ] Details-panel med preview
- [ ] Panel-navigation (Tab)

### Fas 4: UI - Interaktion
- [ ] Vim-navigation (j/k/g/G)
- [ ] Filter-val påverkar work items
- [ ] Öppna i browser (Enter)
- [ ] Fullskärms-detaljvy (v)

### Fas 5: Polish
- [ ] Hjälp-overlay (?)
- [ ] Refresh (Ctrl+r)
- [ ] Felhantering och loading states
- [ ] Färgtema och styling

---

## Framtida features (post-MVP)

Dessa features är medvetet exkluderade från MVP men kan läggas till senare:

### Prioritet 1 (nästa iteration)
- [ ] Pipelines/Builds-vy
- [ ] Pull Requests-vy
- [ ] Projekt-switcher

### Prioritet 2
- [ ] Skapa work items
- [ ] Redigera work items (state, assigned)
- [ ] WIQL query editor
- [ ] Kommentarer

### Prioritet 3
- [ ] Kanban board-vy
- [ ] Notifications
- [ ] Offline cache
- [ ] Themes (dark/light/custom)

---

## Referenser

### Inspiration
- [jira-cli](https://github.com/ankitpokhrel/jira-cli) - Go/Bubble Tea, tabellbaserad
- [JiraTUI](https://jiratui.sh/) - Python/Textual, panel-baserad
- [lazygit](https://github.com/jesseduffield/lazygit) - Go/gocui, panel-layout
- [k9s](https://github.com/derailed/k9s) - Go/tview, Kubernetes TUI

### Dokumentation
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- [Azure DevOps REST API](https://learn.microsoft.com/en-us/rest/api/azure/devops/)
- [WIQL Syntax](https://learn.microsoft.com/en-us/azure/devops/boards/queries/wiql-syntax)

### Azure DevOps Go Libraries
- [microsoft/azure-devops-go-api](https://github.com/microsoft/azure-devops-go-api) - Officiellt men begränsat
- REST API direkt rekommenderas för enklare implementation
