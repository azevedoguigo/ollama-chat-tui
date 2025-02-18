package ui

import (
	"github.com/azevedoguigo/ollama-chat-tui/internal/handler"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ChatList struct {
	list        *tview.List
	chatManager *handler.ChatManager
}

func NewChatList(chatManager *handler.ChatManager) *ChatList {
	cl := &ChatList{
		list:        tview.NewList().ShowSecondaryText(false),
		chatManager: chatManager,
	}

	cl.list.AddItem("New Chat", "", 'n', nil)
	cl.Refresh()
	cl.list.SetBorder(true).SetTitle("Chats")
	return cl
}

func (cl *ChatList) Refresh() {
	cl.list.Clear()
	cl.list.AddItem("New Chat", "", 'n', nil)
	chats := cl.chatManager.GetAllChats()
	for _, chat := range chats {
		cl.list.AddItem(chat.Title, chat.ID.String(), 0, nil)
	}
}

func (cl *ChatList) GetPrimitive() tview.Primitive {
	return cl.list
}

func (cl *ChatList) SetSelectedFunc(handler func(index int, mainText, secondaryText string, shortcut rune)) {
	cl.list.SetSelectedFunc(handler)
}

func (cl *ChatList) SetInputCapture(handler func(event *tcell.EventKey) *tcell.EventKey) {
	cl.list.SetInputCapture(handler)
}

func (cl *ChatList) GetCurrentItemIndex() int {
	return cl.list.GetCurrentItem()
}

func (cl *ChatList) GetSecondaryText(index int) string {
	_, secondary := cl.list.GetItemText(index)
	return secondary
}
