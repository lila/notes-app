package tui

import (
	"fmt"
	"note-app/concept"
	"note-app/note"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewState int

const (
	viewingDate viewState = iota
	viewingConcepts
	viewingNotesInConcept
	editingNote
)

// A message to indicate that the concepts have been loaded.
type conceptsLoadedMsg struct{ concepts []concept.Concept }

type model struct {
	notes           []note.Note
	concepts        []concept.Concept
	groupedNotes    map[string][]note.Note // For date view
	sortedDates     []string               // For date view
	cursor          int
	view            viewState
	selectedConcept *concept.Concept
	loading         bool
	textarea        textarea.Model
	editingNoteID   int
	width           int
}

func initialModel() model {
	notes := note.All()
	groupedNotes, sortedDates := groupNotesByDate(notes)

	ta := textarea.New()
	ta.Placeholder = "Enter your note..."
	ta.Focus()

	return model{
		notes:        notes,
		groupedNotes: groupedNotes,
		sortedDates:  sortedDates,
		view:         viewingDate,
		loading:      true, // Start in loading state
		textarea:     ta,
	}
}

// A command to generate the concepts in the background.
func generateConceptsCmd(notes []note.Note, force bool) tea.Cmd {
	return func() tea.Msg {
		concepts := concept.Cluster(notes, force)
		return conceptsLoadedMsg{concepts: concepts}
	}
}

func groupNotesByDate(notes []note.Note) (map[string][]note.Note, []string) {
	grouped := make(map[string][]note.Note)
	for _, n := range notes {
		dateStr := n.Created.Format("January 2, 2006")
		grouped[dateStr] = append(grouped[dateStr], n)
	}

	dates := make([]string, 0, len(grouped))
	for d := range grouped {
		dates = append(dates, d)
	}

	sort.Slice(dates, func(i, j int) bool {
		t1, _ := time.Parse("January 2, 2006", dates[i])
		t2, _ := time.Parse("January 2, 2006", dates[j])
		return t1.After(t2)
	})

	return grouped, dates
}

func (m model) Init() tea.Cmd {
	// Start the concept generation as soon as the app starts.
	return tea.Batch(textarea.Blink, generateConceptsCmd(m.notes, false))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.textarea.SetWidth(msg.Width - docStyle.GetHorizontalMargins() - docStyle.GetHorizontalPadding() - 2)
		return m, nil

	// Handle the message from our background task.
	case conceptsLoadedMsg:
		m.concepts = msg.concepts
		m.loading = false
		return m, nil

	case tea.KeyMsg:
		if m.view == editingNote {
			switch msg.String() {
			case "ctrl+s":
				note.Update(m.editingNoteID, m.textarea.Value())
				m.notes = note.All()
				m.groupedNotes, m.sortedDates = groupNotesByDate(m.notes)
				m.view = viewingDate
			case "esc":
				m.view = viewingDate
			default:
				m.textarea, cmd = m.textarea.Update(msg)
				cmds = append(cmds, cmd)
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "r":
				if m.view == viewingConcepts {
					m.loading = true
					m.concepts = nil // Clear old concepts
					return m, generateConceptsCmd(m.notes, true)
				}
			case "tab":
				if m.view == viewingDate {
					m.view = viewingConcepts
				} else {
					m.view = viewingDate
				}
				m.cursor = 0
				m.selectedConcept = nil
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				max := 0
				switch m.view {
				case viewingDate:
					max = len(m.notes) - 1
				case viewingConcepts:
					max = len(m.concepts) - 1
				case viewingNotesInConcept:
					if m.selectedConcept != nil {
						max = len(m.selectedConcept.Notes) - 1
					}
				}
				if m.cursor < max {
					m.cursor++
				}
			case "enter":
				if m.view == viewingDate && len(m.notes) > 0 {
					m.view = editingNote
					m.editingNoteID = m.notes[m.cursor].ID
					m.textarea.SetValue(m.notes[m.cursor].Content)
				} else if m.view == viewingConcepts && !m.loading && len(m.concepts) > 0 {
					m.selectedConcept = &m.concepts[m.cursor]
					m.view = viewingNotesInConcept
					m.cursor = 0
				}
			case "backspace":
				if m.view == viewingNotesInConcept {
					m.view = viewingConcepts
					m.cursor = 0
					m.selectedConcept = nil
				}
			case "d":
				if m.view == viewingDate && len(m.notes) > 0 {
					note.Delete(m.notes[m.cursor].ID)
					m.notes = note.All()
					m.groupedNotes, m.sortedDates = groupNotesByDate(m.notes)
					if m.cursor >= len(m.notes) && len(m.notes) > 0 {
						m.cursor = len(m.notes) - 1
					}
				}
			}
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s strings.Builder

	// Title
	title := titleStyle.Render(" Note App ")
	s.WriteString(title + "\n\n")

	// Content
	if len(m.notes) == 0 {
		s.WriteString("No notes yet!\n\n")
		s.WriteString("You can add a new note from your terminal:\n")
		s.WriteString(selectedItemStyle.Render("go run . add --content \"My first note\"") + "\n")
	} else {
		switch m.view {
		case viewingDate:
			s.WriteString(m.dateView())
		case viewingConcepts:
			s.WriteString(m.conceptView())
		case viewingNotesInConcept:
			s.WriteString(m.notesInConceptView())
		case editingNote:
			s.WriteString(m.editView())
		}
	}

	// Help
	s.WriteString("\n" + m.helpView())

	return docStyle.Width(m.width - docStyle.GetHorizontalMargins()).Render(s.String())
}

func (m model) dateView() string {
	var s strings.Builder
	s.WriteString("Journal (Tab to switch to Concepts)\n")

	noteIndex := 0
	for _, date := range m.sortedDates {
		s.WriteString(dateStyle.Render(date) + "\n")
		for _, n := range m.groupedNotes[date] {
			style := itemStyle
			if m.cursor == noteIndex {
				style = selectedItemStyle
			}
			s.WriteString(style.Render("  - "+n.Content) + "\n")
			noteIndex++
		}
	}
	return s.String()
}

func (m model) conceptView() string {
	var s strings.Builder
	s.WriteString("Concepts (Tab to switch to Journal)\n\n")

	if m.loading {
		s.WriteString("Regenerating concept groupings...\n")
	} else {
		for i, c := range m.concepts {
			style := itemStyle
			if m.cursor == i {
				style = selectedItemStyle
			}
			s.WriteString(style.Render(fmt.Sprintf("> %s (%d notes)", c.Name, len(c.Notes))) + "\n")
		}
	}
	return s.String()
}

func (m model) notesInConceptView() string {
	var s strings.Builder
	if m.selectedConcept != nil {
		s.WriteString(conceptStyle.Render(m.selectedConcept.Name) + " (Backspace to return)\n\n")
		for i, n := range m.selectedConcept.Notes {
			style := itemStyle
			if m.cursor == i {
				style = selectedItemStyle
			}
			s.WriteString(style.Render("- "+n.Content) + "\n")
		}
	}
	return s.String()
}

func (m model) editView() string {
	return m.textarea.View()
}

func (m model) helpView() string {
	switch m.view {
	case editingNote:
		help := "ctrl+s: save • esc: cancel"
		return helpStyle.Render(help)
	case viewingConcepts:
		left := "↑/k up • ↓/j down • tab switch"
		right := "enter select • r refresh • q quit"

		contentWidth := m.width - docStyle.GetHorizontalMargins() - docStyle.GetHorizontalPadding() - 2

		spacerWidth := contentWidth - lipgloss.Width(left) - lipgloss.Width(right)
		if spacerWidth < 1 {
			spacerWidth = 1
		}
		spacer := strings.Repeat(" ", spacerWidth)

		return helpStyle.Render(left + spacer + right)
	default:
		left := "↑/k up • ↓/j down • tab switch"
		right := "enter select • d delete • q quit"

		contentWidth := m.width - docStyle.GetHorizontalMargins() - docStyle.GetHorizontalPadding() - 2

		spacerWidth := contentWidth - lipgloss.Width(left) - lipgloss.Width(right)
		if spacerWidth < 1 {
			spacerWidth = 1
		}
		spacer := strings.Repeat(" ", spacerWidth)

		return helpStyle.Render(left + spacer + right)
	}
}

func NewProgram() *tea.Program {
	return tea.NewProgram(initialModel())
}
