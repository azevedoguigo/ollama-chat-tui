package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputField struct {
	input *tview.InputField
}

func NewInputField() *InputField {
	in := tview.NewInputField().SetLabel("You: ").SetFieldWidth(0)

	return &InputField{
		input: in,
	}
}

func (ifd *InputField) GetPrimitive() tview.Primitive {
	return ifd.input
}

func (ifd *InputField) GetText() string {
	return ifd.input.GetText()
}

func (ifd *InputField) SetText(text string) {
	ifd.input.SetText(text)
}

func (ifd *InputField) SetDoneFunc(handler func(key tcell.Key)) {
	ifd.input.SetDoneFunc(handler)
}
