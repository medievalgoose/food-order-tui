package main

import (
	"fmt"
	"food-order-ui/db"
	"food-order-ui/imgProc"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := listPages()

	detectUserInput(app, pages)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func setMenuInfoText(food *db.Food, foodInfoText *tview.TextView) {
	foodInfoText.Clear()
	text := fmt.Sprintf("Nama makanan: %v \nHarga: %v \nDeskripsi: %v", food.Name, food.Price, food.Description)
	foodInfoText.SetText(text)

	foodInfoText.SetBorder(true)
	foodInfoText.SetTitle("   DETAIL MENU   ")
}

func assignMenuView() *tview.Flex {
	menuList := tview.NewList()
	foodList := db.GetAllFoods()

	for _, food := range foodList {
		idString := strconv.Itoa(food.ID)

		var idRune []rune
		for _, char := range idString {
			idRune = append(idRune, char)
		}

		menuList.AddItem(food.Name, strconv.Itoa(food.Price), rune(idRune[0]), nil)
	}

	menuList.SetBorder(true)
	menuList.SetTitle("   MENU MAKANAN   ")

	foodInfoText := tview.NewTextView()
	foodImage := tview.NewImage()
	foodImage.SetImage(imgProc.GetImage("sate-padang.png"))

	// menuList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
	// 	setMenuInfoText(&foodList[index], foodInfoText)
	// })

	menuList.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		setMenuInfoText(&foodList[index], foodInfoText)
	})

	// menuDetailFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	// menuDetailFlex.AddItem(foodImage, 0, 4, false).SetBorder(true).SetBorderPadding(1, 0, 0, 0)
	// menuDetailFlex.AddItem(foodInfoText, 0, 6, false)

	menuDetailGrid := tview.NewGrid()
	menuDetailGrid.SetRows(0)
	menuDetailGrid.SetColumns(0)

	menuDetailGrid.AddItem(foodImage, 0, 0, 1, 2, 1, 0, true)
	menuDetailGrid.AddItem(foodInfoText, 1, 0, 1, 2, 1, 0, false)

	menuDetailGrid.SetBorders(true)

	flexbox := tview.NewFlex()
	flexbox.AddItem(menuList, 0, 1, true)
	// flexbox.AddItem(menuDetailFlex, 0, 1, false)
	flexbox.AddItem(menuDetailGrid, 0, 1, false)

	return flexbox
}

func listPages() *tview.Pages {
	pages := tview.NewPages()

	MenuPage := assignMenuView()

	pages.AddPage("menu", MenuPage, true, true)
	pages.AddPage("test", tview.NewBox().SetBorder(true), false, false)

	return pages
}

func detectUserInput(app *tview.Application, pages *tview.Pages) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 97: // Key a is pressed.
			pages.SwitchToPage("test")
		}
		return event
	})
}
