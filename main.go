package main

import (
	"fmt"
	"food-order-ui/db"
	"food-order-ui/imgProc"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var cartFoods []db.Food
var app = tview.NewApplication()
var pages = tview.NewPages()
var itemSelected bool = false

var orderListDB []db.Order

func main() {
	// app := tview.NewApplication()
	pages = listPages()

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

	// Get list of foods from DB.
	foodList := db.GetAllFoods()

	// List the foods.
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
	// var cartFoods []db.Food
	cartFoodMap := make(map[int]int)
	cartButtons := tview.NewButton("Confirm Order")
	cartDetail := tview.NewTextView().SetText("Total: 0")
	var currentSelectedFood db.Food
	var currentSelectedIndex int
	currentMenuCount := tview.NewTextView().SetText("x0").SetTextAlign(tview.AlignCenter)

	cartButtons.SetSelectedFunc(func() {
		// Insert the order data here.
		db.CreateNewOrder(cartFoods, cartFoodMap)
		changeToConfirmPage()
	})

	// Show the menu detail if user selected the item by pressing enter or space.
	menuList.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		setMenuInfoText(&foodList[index], foodInfoText)
		changeFoodImage(&foodList[index], foodImage)

		currentSelectedFood = foodList[index]
		currentSelectedIndex = index

		itemSelected = true

		reinitializeMenuCount(currentMenuCount, currentSelectedFood, cartFoodMap)
	})

	// Menu Actions related
	decreaseButton := tview.NewButton("-").SetSelectedFunc(func() {
		if itemSelected {
			refreshMenuCount("decrease", currentMenuCount, cartFoodMap, currentSelectedFood)
			refreshMenuDetail(menuList, currentSelectedIndex, cartFoodMap, currentSelectedFood)
			refreshTotalCount(cartDetail, cartFoods, cartFoodMap)
		}
	})

	increaseButton := tview.NewButton("+").SetSelectedFunc(func() {
		if itemSelected {
			refreshMenuCount("increase", currentMenuCount, cartFoodMap, currentSelectedFood)
			refreshMenuDetail(menuList, currentSelectedIndex, cartFoodMap, currentSelectedFood)
			refreshTotalCount(cartDetail, cartFoods, cartFoodMap)
		}
	})

	menuActionGrid := tview.NewGrid()
	menuActionGrid.SetRows(0)
	menuActionGrid.SetColumns(0)
	menuActionGrid.AddItem(decreaseButton, 0, 0, 1, 1, 0, 0, false)
	menuActionGrid.AddItem(currentMenuCount, 0, 1, 1, 1, 0, 0, false)
	menuActionGrid.AddItem(increaseButton, 0, 2, 1, 1, 0, 0, false)

	menuDetailGrid := tview.NewGrid()
	menuDetailGrid.SetRows(0, 0, 3)
	menuDetailGrid.SetColumns(0)
	menuDetailGrid.AddItem(foodImage, 0, 0, 1, 2, 1, 0, true)
	menuDetailGrid.AddItem(foodInfoText, 1, 0, 1, 2, 1, 0, false)
	menuDetailGrid.AddItem(menuActionGrid, 2, 0, 1, 2, 1, 0, false)
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
	flexbox.AddItem(menuDetailGrid, 0, 1, false)

	return flexbox
}

func refreshMenuDetail(menuList *tview.List, index int, cartMap map[int]int, currentFood db.Food) {
	_, itemExist := cartMap[currentFood.ID]
	if itemExist {
		menuList.SetItemText(index, currentFood.Name, strconv.Itoa(currentFood.Price)+" | x"+strconv.Itoa(cartMap[currentFood.ID]))
	} else {
		menuList.SetItemText(index, currentFood.Name, strconv.Itoa(currentFood.Price))
	}

}

func refreshMenuCount(action string, menuCountText *tview.TextView, cartMap map[int]int, currentFood db.Food) {
	_, itemExist := cartMap[currentFood.ID]
	switch action {
	case "increase":
		if itemExist {
			cartMap[currentFood.ID] = cartMap[currentFood.ID] + 1
		} else {
			cartMap[currentFood.ID] = 1
			cartFoods = append(cartFoods, currentFood)
		}
	case "decrease":
		if itemExist {
			if cartMap[currentFood.ID] > 1 {
				cartMap[currentFood.ID] = cartMap[currentFood.ID] - 1
			} else {
				delete(cartMap, currentFood.ID)
				cartFoods = removeItem(cartFoods, currentFood.ID)
			}
		}
	}

	menuCountText.Clear()
	menuCountText.SetText("x" + strconv.Itoa(cartMap[currentFood.ID]))
}

func reinitializeMenuCount(currentMenuCount *tview.TextView, currentSelectedFood db.Food, cartFoodMap map[int]int) {
	currentMenuCount.SetText("x" + strconv.Itoa(cartFoodMap[currentSelectedFood.ID]))
}

func refreshTotalCount(cartDetail *tview.TextView, cartFoods []db.Food, cartFoodMap map[int]int) {
	total := 0
	for _, food := range cartFoods {
		total += food.Price * cartFoodMap[food.ID]
	}
	cartDetail.Clear()
	cartDetail.SetText("Total: " + strconv.Itoa(total))
}

func listPages() *tview.Pages {
	pages := tview.NewPages()

	MenuPage := assignMenuView()
	MainMenuPage := assignMainMenuView()
	OrderListPage := assignOrderListView()

	modal := tview.NewModal()
	modal.SetText("Your order has been confirmed!")
	modal.AddButtons([]string{"OK"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.Stop()
	})

	modalOrder := tview.NewModal()
	modalOrder.SetText("The order has been processed.")
	modalOrder.AddButtons([]string{"OK"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.Stop()
	})

	pages.AddPage("modal", modal, true, false)
	pages.AddPage("orderModal", modalOrder, true, false)
	pages.AddPage("menu", MenuPage, true, false)
	pages.AddPage("mainMenu", MainMenuPage, true, true)
	pages.AddPage("orderList", OrderListPage, true, false)

	return pages
}

func detectUserInput(app *tview.Application, pages *tview.Pages) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case rune('c'): // Key a is pressed.
			pages.SwitchToPage("orderList")
		}
		return event
	})
}

func removeItem(cartFoods []db.Food, foodID int) []db.Food {
	var newArray []db.Food
	for _, food := range cartFoods {
		if food.ID != foodID {
			newArray = append(newArray, food)
		}
	}
	return newArray
}

func changeToConfirmPage() {
	pages.SwitchToPage("modal")
}

func assignMainMenuView() *tview.Grid {
	programTitle := tview.NewTextView().SetText("Food Order TUI").SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreen)

	menuListButton := tview.NewButton("Menu List").SetSelectedFunc(func() {
		pages.SwitchToPage("menu")
	})

	orderListButton := tview.NewButton("Orders").SetSelectedFunc(func() {
		pages.SwitchToPage("orderList")
	})

	mainMenuGrid := tview.NewGrid()
	mainMenuGrid.SetRows(0)
	mainMenuGrid.SetColumns(0)
	mainMenuGrid.AddItem(programTitle, 0, 0, 1, 2, 0, 0, false)
	mainMenuGrid.AddItem(menuListButton, 1, 0, 1, 1, 0, 0, false)
	mainMenuGrid.AddItem(orderListButton, 1, 1, 1, 1, 0, 0, false)
	mainMenuGrid.SetBorders(true)
	mainMenuGrid.SetBorderPadding(5, 5, 5, 5)

	return mainMenuGrid
}

func assignOrderListView() *tview.Grid {
	orderList := tview.NewList()
	orderListDB = db.GetAllOrders()
	orderDetails := tview.NewTextView()
	var currentSelectedOrder db.Order

	for _, order := range orderListDB {
		orderList.AddItem("Order ID: "+strconv.Itoa(order.ID), "", rune(0), nil)
	}

	orderList.SetSelectedFunc(func(index int, primaryText, secondaryText string, shortcut rune) {
		currentSelectedOrder = orderListDB[index]
		refreshOrderDetails(currentSelectedOrder, orderDetails)
	})

	doneButton := tview.NewButton("Finish").SetSelectedFunc(func() {
		resetOrder("done", orderList, orderDetails, currentSelectedOrder)
	})

	cancelButton := tview.NewButton("Cancel").SetSelectedFunc(func() {
		resetOrder("cancel", orderList, orderDetails, currentSelectedOrder)
	})

	orderListGrid := tview.NewGrid()
	orderListGrid.SetRows(0, 3)
	orderListGrid.SetColumns(0)
	orderListGrid.AddItem(orderList, 0, 0, 2, 1, 0, 0, true)
	orderListGrid.AddItem(orderDetails, 0, 1, 1, 2, 0, 0, false)

	// Confirm button
	orderListGrid.AddItem(doneButton, 1, 1, 1, 1, 0, 0, false)

	// Cancel button
	orderListGrid.AddItem(cancelButton, 1, 2, 1, 1, 0, 0, false)

	orderListGrid.SetBorders(true)

	return orderListGrid
}

func refreshOrderDetails(currentSelectedOrder db.Order, orderDetailText *tview.TextView) {
	detailText := fmt.Sprintf("Order ID: %v\n Date: %v\n", currentSelectedOrder.ID, currentSelectedOrder.Date)

	for _, items := range currentSelectedOrder.Items {
		detailText += fmt.Sprintf("\n - %v | x%v", items.FoodName, items.Quantity)
	}

	orderDetailText.Clear()
	orderDetailText.SetText(detailText)
	orderDetailText.SetBorder(true)
}

func resetOrder(action string, orderList *tview.List, orderDetailText *tview.TextView, currentSelectedOrder db.Order) {
	switch action {
	case "done":
		db.ChangeOrderStatusToDone(currentSelectedOrder)
	case "cancel":
		db.ChangeOrderStatusToCancelled(currentSelectedOrder)
	}

	pages.SwitchToPage("orderModal")
}
