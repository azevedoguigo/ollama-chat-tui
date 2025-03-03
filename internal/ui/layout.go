package ui

import "github.com/rivo/tview"

func BuildMainLayout(chatList tview.Primitive, chatView tview.Primitive, inputField tview.Primitive) *tview.Flex {
	chatFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chatView, 0, 1, false).
		AddItem(inputField, 1, 1, true)

	mainFlex := tview.NewFlex().
		AddItem(chatList, 20, 1, false).
		AddItem(chatFlex, 0, 1, true)

	return mainFlex
}
