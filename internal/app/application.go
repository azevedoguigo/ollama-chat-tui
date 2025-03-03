package app

import (
	"fmt"
	"log"

	"github.com/azevedoguigo/ollama-chat-tui/internal/handler"
	"github.com/azevedoguigo/ollama-chat-tui/internal/ollama"
	"github.com/azevedoguigo/ollama-chat-tui/internal/storage"
	"github.com/azevedoguigo/ollama-chat-tui/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Run() error {
	configDir := ".ollama-chat-tui"
	chatsDir := "chats"

	chatManager, err := handler.NewChatManager(configDir, chatsDir)
	if err != nil {
		return err
	}

	currentModel := "deepseek-r1"
	availableModels := []string{"deepseek-r1", "gemma2", "mistral"}

	chatList := ui.NewChatList(chatManager)
	chatView := ui.NewChatView()
	inputField := ui.NewInputField()

	chatFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chatView.GetPrimitive(), 0, 1, false).
		AddItem(inputField.GetPrimitive(), 3, 1, true)

	mainChatLayout := tview.NewFlex().
		AddItem(chatList.GetPrimitive(), 20, 1, false).
		AddItem(chatFlex, 0, 1, true)

	pages := tview.NewPages().
		AddPage("chat", mainChatLayout, true, true)

	app := tview.NewApplication()

	var currentChat *storage.ChatSession

	updateChatView := func() {
		if currentChat != nil {
			chatView.Update(currentChat)
		}
	}

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			userInput := inputField.GetText()
			inputField.SetText("")
			if userInput == "" {
				return
			}

			if currentChat == nil {
				title := fmt.Sprintf("Chat %d", len(chatManager.GetAllChats())+1)
				currentChat = chatManager.AddChat(title)
				chatList.Refresh()
			}

			chatID := currentChat.ID.String()

			if err := chatManager.AppendMessage(chatID, "user", userInput); err != nil {
				chatView.SetText(fmt.Sprintf("[red]Error adding message: %v[-]", err))
				return
			}

			if err := chatManager.AppendMessage(chatID, "assistant", ""); err != nil {
				chatView.SetText(fmt.Sprintf("[red]Error creating response message: %v[-]", err))
				return
			}

			updateChatView()

			history := make([]storage.Message, len(currentChat.Messages))
			copy(history, currentChat.Messages)

			go func() {
				err := ollama.QueryOllamaStream(currentModel, history[:len(history)-1], func(chunk string) {
					if err := chatManager.UpdateLastMessage(chatID, chunk); err != nil {
						log.Printf("Error updating message: %v", err)
					}

					updatedChat, _ := chatManager.GetChatByID(chatID)
					currentChat = updatedChat
					app.QueueUpdateDraw(func() {
						updateChatView()
					})
				})
				if err != nil {
					errMsg := fmt.Sprintf("\n\n[red]Error: %v[-]", err)
					if updateErr := chatManager.UpdateLastMessage(chatID, errMsg); updateErr != nil {
						log.Printf("Error updating message with error: %v", updateErr)
					}
					app.QueueUpdateDraw(func() {
						updateChatView()
					})
				}
			}()
		}
	})

	chatList.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		if index == 0 {
			currentChat = nil
			chatView.Clear()
			inputField.SetText("")
		} else {
			chat, exists := chatManager.GetChatByID(secondaryText)
			if exists {
				currentChat = chat
				updateChatView()
			}
		}
	})

	chatList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDelete || event.Rune() == 'd' {
			currentItem := chatList.GetCurrentItemIndex()
			if currentItem > 0 {
				chatID := chatList.GetSecondaryText(currentItem)
				modal := tview.NewModal().
					SetText("Delete this chat permanently?").
					AddButtons([]string{"Cancel", "Delete"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "Delete" {
							if err := chatManager.DeleteChat(chatID); err != nil {
								chatView.SetText(fmt.Sprintf("[red]Error deleting chat: %v[-]", err))
							} else {
								if currentChat != nil && currentChat.ID.String() == chatID {
									currentChat = nil
									chatView.Clear()
								}
								chatList.Refresh()
							}
						}
						app.SetRoot(pages, true)
					})
				app.SetRoot(modal, false)
			}
			return nil
		}

		return event
	})

	mainChatLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == rune(tcell.KeyCtrlD) {
			openSettings(app, pages, currentModel, availableModels, func(newModel string) {
				currentModel = newModel
				pages.SwitchToPage("chat")
				app.SetFocus(inputField.GetPrimitive())
			})
			return nil
		}
		return event
	})

	return app.SetRoot(pages, true).EnableMouse(true).Run()
}

func openSettings(
	app *tview.Application,
	pages *tview.Pages,
	currentModel string,
	availableModels []string,
	onSave func(newModel string),
) {
	settingsPage := ui.NewSettingsPage(currentModel, availableModels, onSave)
	pages.AddAndSwitchToPage("settings", settingsPage.GetPrimitive(), true)
}
