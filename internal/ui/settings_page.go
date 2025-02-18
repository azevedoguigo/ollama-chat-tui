// internal/ui/settings.go
package ui

import (
	"github.com/rivo/tview"
)

type SettingsPage struct {
	form         *tview.Form
	currentModel string
	models       []string
	onSave       func(newModel string)
}

func NewSettingsPage(currentModel string, models []string, onSave func(newModel string)) *SettingsPage {
	sp := &SettingsPage{
		form:         tview.NewForm(),
		currentModel: currentModel,
		models:       models,
		onSave:       onSave,
	}

	dropdown := tview.NewDropDown().
		SetLabel("Model: ").
		SetOptions(models, nil)

	for index, m := range models {
		if m == currentModel {
			dropdown.SetCurrentOption(index)
			break
		}
	}

	dropdown.SetSelectedFunc(func(text string, index int) {
		sp.currentModel = text
	})

	sp.form.AddFormItem(dropdown)
	sp.form.AddButton("Save", func() {
		if sp.onSave != nil {
			sp.onSave(sp.currentModel)
		}
	})
	sp.form.AddButton("Cancel", func() {
		if sp.onSave != nil {
			sp.onSave(currentModel)
		}
	})
	sp.form.SetBorder(true).SetTitle("Settings")
	return sp
}

func (sp *SettingsPage) GetPrimitive() tview.Primitive {
	return sp.form
}
