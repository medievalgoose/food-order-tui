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
	text := fmt.Sprintf("\nNama makanan: %v \n\nHarga: %v \n\nDeskripsi: %v", food.Name, food.Price, food.Description)
	foodInfoText.SetText(text)

	foodInfoText.SetBorder(true)
	foodInfoText.SetTitle("   DETAIL MENU   ")
}

func changeFoodImage(food *db.Food, foodImage *tview.Image) {
	foodImage.SetImage(imgProc.GetImage(food.Image))
}

func assignMenuView() *tview.Flex {
	menuList := tview.NewList()
	menuMap := make(map[string]int)
	// menuList.ShowSecondaryText(false)
	foodList := db.GetAllFoods()

	for _, food := range foodList {
		// No shortcut
		menuList.AddItem(food.Name, strconv.Itoa(food.Price), rune(0), nil)
		menuMap[food.Name] = 0
	}

	menuList.SetBorder(true)
	menuList.SetTitle("   MENU MAKANAN   ")

	foodInfoText := tview.NewTextView()
	foodImage := tview.NewImage()

	// Ordered items data
	var cartFoods []db.Food
	cartFoodMap := make(map[int]int)
	cartButtons := tview.NewButton("Confirm Order")
	cartDetail := tview.NewTextView().SetText("Total: 0")

	menuList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		setMenuInfoText(&foodList[index], foodInfoText)
		changeFoodImage(&foodList[index], foodImage)
	})

	// Show the menu detail if user selected the item by pressing enter or space.
	menuList.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		setMenuInfoText(&foodList[index], foodInfoText)
		changeFoodImage(&foodList[index], foodImage)

		// Increment the order count
		menuMap[foodList[index].Name] = menuMap[foodList[index].Name] + 1
		menuList.SetItemText(index, foodList[index].Name, strconv.Itoa(foodList[index].Price)+" | x"+strconv.Itoa(menuMap[foodList[index].Name]))

		// Check if the item already exist on the list
		_, itemExist := cartFoodMap[foodList[index].ID]
		if itemExist {
			cartFoodMap[foodList[index].ID] = cartFoodMap[foodList[index].ID] + 1
		} else {
			cartFoodMap[foodList[index].ID] = 1
			cartFoods = append(cartFoods, foodList[index])
		}

		refreshCartDetail(cartDetail, cartFoods, cartFoodMap)
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

	menuListGrid := tview.NewGrid()
	menuListGrid.SetRows(0, 3, 3)
	menuListGrid.SetColumns(0)
	menuListGrid.AddItem(menuList, 0, 0, 1, 2, 0, 0, true)
	menuListGrid.AddItem(cartDetail, 1, 0, 1, 2, 0, 0, false)
	menuListGrid.AddItem(cartButtons, 2, 0, 1, 2, 0, 0, false)
	menuListGrid.SetBorders(true)

	flexbox := tview.NewFlex()
	flexbox.AddItem(menuListGrid, 0, 1, true)
	// flexbox.AddItem(menuDetailFlex, 0, 1, false)
	flexbox.AddItem(menuDetailGrid, 0, 1, false)

	return flexbox
}

func refreshCartDetail(cartDetail *tview.TextView, cartFoods []db.Food, cartMap map[int]int) {
	cartDetail.Clear()
	total := 0

	for _, food := range cartFoods {
		total += cartMap[food.ID] * food.Price
	}

	cartDetail.SetText("Total: " + strconv.Itoa(total))
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
