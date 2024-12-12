package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Show() {
	// myApp := app.New()
	// myWindow := myApp.NewWindow("Data Management")
	// myWindow.Resize(fyne.NewSize(800, 600))
	// myWindow.SetFixedSize(true)

	// // Create content areas
	// homeLabel := widget.NewLabel("Welcome to the Home Page")
	// userTable := createTable("User Data")
	// productTable := createTable("Product Data")

	// // Default content area
	// contentArea := container.NewMax(homeLabel)
	// // Define navigation buttons
	// nav := container.NewVBox(
	// 	widget.NewButton("Home", func() {
	// 		switchContent(contentArea, homeLabel)
	// 	}),
	// 	widget.NewButton("User", func() {
	// 		switchContent(contentArea, userTable)
	// 	}),
	// 	widget.NewButton("Product", func() {
	// 		switchContent(contentArea, productTable)
	// 	}),
	// )

	// split := container.NewHSplit(nav, contentArea)
	// split.SetOffset(0.2)
	// myWindow.SetContent(split)
	// myWindow.ShowAndRun()
	a := app.NewWithID("io.fyne.demo")
	w := a.NewWindow("Shop - 后台管理系统")

	w.SetMaster()

	content := container.NewStack()
	title := widget.NewRichText(
		&widget.TextSegment{
			Text: "title",
			Style: widget.RichTextStyle{
				SizeName: theme.SizeNameHeadingText,
			},
		},
	)

	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	setTable := func(t Table) {

		title.Segments[0].(*widget.TextSegment).Text = t.Title
		title.Refresh()
		intro.SetText(t.Intro)

		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setTable, false))
	} else {
		split := container.NewHSplit(makeNav(setTable, true), tutorial)
		split.Offset = 0.2
		w.SetContent(split)
	}
	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()

}

func makeNav(setTable func(table Table), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return TablesIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := TablesIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := Tables[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := Tables[uid]; ok {
				a.Preferences().SetString("currentTable", uid)
				setTable(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback("currentTable", "welcome")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			// a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantDark})
		}),
		widget.NewButton("Light", func() {
			// a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantLight})
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}
