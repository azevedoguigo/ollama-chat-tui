// internal/ui/settings.go
package ui

import (
	"github.com/azevedoguigo/ollama-chat-tui/internal/service"
	"github.com/rivo/tview"
)

type SettingsPage struct {
	form         *tview.Form
	currentModel string
	models       []string
	onSave       func(newModel string)
}

func NewSettingsPage(
	currentModel string,
	models *service.FindModelsResponse,
	onSave func(newModel string),
) *SettingsPage {
	var modelsNames []string

	if len(models.Models) > 0 {
		for i := range models.Models {
			modelsNames = append(modelsNames, models.Models[i].Name)
		}
	} else {
		modelsNames = append(modelsNames, "No models avaliable")
	}

	sp := &SettingsPage{
		form:         tview.NewForm(),
		currentModel: currentModel,
		models:       modelsNames,
		onSave:       onSave,
	}

	dropdown := tview.NewDropDown().
		SetLabel("Model: ").
		SetOptions(modelsNames, nil)

	for index, m := range modelsNames {
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
