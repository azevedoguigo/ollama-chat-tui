package ui

import (
	"fmt"

	"github.com/azevedoguigo/ollama-chat-tui/internal/storage"
	"github.com/rivo/tview"
)

const INITIAL_MESSAGE = "Welcome to Ollama Chat TUI!\nPress CTRL + D to open settings and change AI model or CTRL + C to exit."

type ChatView struct {
	textView *tview.TextView
}

func NewChatView() *ChatView {
	cv := &ChatView{
		textView: tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetWordWrap(true),
	}
	cv.textView.SetBorder(true).SetTitle("Chat")
	cv.textView.SetText(INITIAL_MESSAGE)

	return cv
}

func (cv *ChatView) Update(chat *storage.ChatSession) {
	cv.textView.Clear()
	for _, msg := range chat.Messages {
		var speaker, color string
		if msg.Role == "user" {
			speaker = "You"
			color = "[green]"
		} else {
			speaker = fmt.Sprintf("Assistant (%s)", chat.Model)
			color = "[blue]"
		}

		fmt.Fprintf(cv.textView, "%s%s:[-] %s\n", color, speaker, msg.Content)
	}

	cv.textView.ScrollToEnd()
}

func (cv *ChatView) Clear() {
	cv.textView.Clear()
	cv.textView.SetText(INITIAL_MESSAGE)
}

func (cv *ChatView) GetPrimitive() tview.Primitive {
	return cv.textView
}

func (cv *ChatView) SetText(text string) {
	cv.textView.SetText(text)
}
