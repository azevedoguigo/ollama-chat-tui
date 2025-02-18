package ui

import (
	"fmt"

	"github.com/azevedoguigo/ollama-chat-tui/internal/storage"
	"github.com/rivo/tview"
)

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
			speaker = "Assistant"
			color = "[blue]"
		}
		fmt.Fprintf(cv.textView, "%s%s:[-] %s\n", color, speaker, msg.Content)
	}
	cv.textView.ScrollToEnd()
}

func (cv *ChatView) Clear() {
	cv.textView.Clear()
}

func (cv *ChatView) GetPrimitive() tview.Primitive {
	return cv.textView
}

func (cv *ChatView) SetText(text string) {
	cv.textView.SetText(text)
}
