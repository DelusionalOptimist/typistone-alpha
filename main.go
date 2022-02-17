package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

var (
	term    = termenv.EnvColorProfile()
	rawText = "The quick brown fox jumps over the lazy black dogs."
	red     = makeFgBgStyle("000", "211")
	green   = makeFgStyle("002")
	timeout = 60
)

type Model struct {
	// InputText stores the text entered by user
	InputText string

	// RawText is for playing the game
	RawText string

	// TextInput is the a model for input operations provided by bubbletea
	TextInput textinput.Model

	// Timeout
	Timer *time.Timer

	// Records the time user starts typing
	StartTime time.Time

	// Stores the duration of the game
	Duration float64

	// Stores the speed
	Speed float64
}

func main() {
	m := initialModel()

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}

}

// initialModel constructs a new Model with some defaults
func initialModel() Model {
	ti := textinput.New()
	ti.Focus()
	return Model{
		RawText:    rawText,
		TextInput:  ti,
	}
}

// prints the initial text
func (m Model) Init() tea.Cmd {
	fmt.Println(rawText)
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// exit if written completely
	if len(m.InputText) == len(m.RawText) {
		// print speed and shits
		endTime := time.Now()
		m.Duration = endTime.Sub(m.StartTime).Minutes()
		m.Speed = float64(len(m.RawText)/5)/m.Duration

		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:

		// start the timeout only when the user starts typing
		if len(m.InputText) == 0 && m.Timer == nil {
			m.StartTime = time.Now()
			m.Timer = time.NewTimer(time.Duration(timeout) * time.Second)

			// timeout
			go func() (tea.Model, tea.Cmd){
				select {
				case <-m.Timer.C:
					endTime := time.Now()
					fmt.Println("\nTime up")

					m.Duration = endTime.Sub(m.StartTime).Minutes()

					// Gross WPM
					m.Speed = float64(len(m.RawText)/5)/m.Duration

					return m, tea.Quit
				}
			}()
		}

		// save proper characters entered by user
		if msg.Type == tea.KeyRunes {
			m.InputText = m.InputText + msg.String()
		}

		switch msg.String() {
		case tea.KeyCtrlC.String(), tea.KeyEsc.String():
			return m, tea.Quit
		case tea.KeyBackspace.String():
			if inputLength := len(m.InputText); inputLength > 0 {
			m.InputText = m.InputText[0 : inputLength-1]
		}

		// increase counter for each wrong char
		default:
		}
	}

	return m, nil
}

// renders the UI
func (m Model) View() string {
	var wrongChars = 0
	var displayText string

	// this colours the text char by char
	for i, char := range m.InputText {
		if char == rune(m.RawText[i]) {
			displayText += green(string(char))
		} else {
			wrongChars++
			displayText += red(string(char))
		}
	}

	accuracy := float64(float64((len(m.RawText) - wrongChars))/float64(len(m.RawText)))*100
	status := fmt.Sprintf("Speed: %f WPM\n\nAccuracy: %f\n\n", m.Speed, accuracy)
	displayText += fmt.Sprintf("\n\n%s", status)

	return displayText
}

// color foreground
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

// color foreground and background with the given value.
func makeFgBgStyle(fg, bg string) func(string) string {
	return termenv.Style{}.
		Foreground(term.Color(fg)).
		Background(term.Color(bg)).
		Styled
}
