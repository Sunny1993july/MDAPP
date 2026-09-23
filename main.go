package main 

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

  "github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"  // tea is an alias
	"github.com/charmbracelet/lipgloss"
	"github.com/go-git/go-git/v5"
)

type activePane int  // custom type for the active pane 

const (
	paneNotes activePane = iota  //  it give the numerical 0 and increment by 1 for the subsequent line.
  paneEditor
	paneGraph
)

// Below is the a custom container that holds the custom variables togather
type model struct {
	vaultPath string 
	notes   []string  // basically a list of the strings that contains the files named 
	selected int 
	active activePane

	//  bubbles's text area
	textarea textarea.Model 

	// text scroll view 
  viewport viewport.Model 

	// map: Key(Note), value(list of the strings) use of the Graph 
	links map[string][]string


	// Layout dimensions
	width int
	height int

	// Bottom staus bar text 
	status string

}


func initialModel(vault string) model {
	// return the initial state 
	
	ta := textarea.New()
	ta.Placeholder = "Write markdown here.... Use [[Note]] to link and - [ ] for the Checklist"
	ta.Focus()

	vp := viewport.New(30, 20)

	m := model{
		vaultPath: vault,
		active: paneNotes,
		textarea: ta,
		viewport: vp,
		links: make(map[string][]string),
		status: "Ready | Ctrl+s: Save/Sync | Tab: Switch Pane | Space: Checklist",

	}
	m.scanVault()
	if len(m.notes) > 0 {
		 m.loadNoteContent(m.notes[0])
	}

	return m 
}

func (m *model) scanVault(){
	// it is a custom function , a method, to check all the files in the md folder to find the [[Linked notes]] as stores
	
	m.notes = []string{}
  m.links = make(map[string][]string)
	wikiLinkRegex := regexp.MustCompile(`\[\[(.*?)\]\]`)

	filepath.Walk(m.vaultPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !string.HasSuffix(info.Name(), ".md"){
				return nil
			}

			rel, _ := filepath.Rel(m.vaultPath, path)
			noteName := strings.TrimSuffix(rel, ".md")
			m.notes = append(m.notes, noteName)

			content, err := os.ReadFile(path)
			if err == nil {
				matches := wikiLinkRegex.FindAllStringSubmatch(string(content), -1)
				for _, match := range matches {
                   if len(match) > 1 {
										 m.links[noteName] = append(m.links[noteName], match[1])
									 }
				}
			}

  return nil
	})
}


