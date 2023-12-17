package db

import (
	"time"
)

type Order struct {
	ID     int
	Date   string
	Status string
	Items  []OrderDetail
}

type OrderDetail struct {
	OrderID  int
	FoodID   int
	FoodName string
	Quantity int
}

func CreateNewOrder(cartFoods []Food, cartFoodMap map[int]int) {
	db := OpenDBConnection()
	defer db.Close()

	currentDate := time.Now()
	orderId := 0

	createNewOrderQuery := "INSERT INTO orders (date, status) VALUES ($1, 'PENDING') RETURNING id;"
	err := db.QueryRow(createNewOrderQuery, currentDate).Scan(&orderId)
	if err != nil {
		panic(err)
	}

	addItemToOrderQuery := "INSERT INTO order_details (order_id, food_id, quantity) VALUES($1, $2, $3);"

	for _, food := range cartFoods {
		_, err := db.Exec(addItemToOrderQuery, orderId, food.ID, cartFoodMap[food.ID])
		if err != nil {
			panic(err)
		}
	}
}

func GetAllOrders() []Order {
	db := OpenDBConnection()
	defer db.Close()

	getAllOrderQuery := "SELECT id, date, status FROM orders WHERE status = 'PENDING';"
	getOrderDetailQuery := "SELECT od.order_id, od.food_id, f.name, od.quantity FROM order_details od JOIN foods f ON od.food_id = f.id WHERE order_id = $1;"
	rows, err := db.Query(getAllOrderQuery)
	if err != nil {
		panic(err)
	}

	var orderList []Order

	for rows.Next() {
		var relevantOrder Order

		err := rows.Scan(&relevantOrder.ID, &relevantOrder.Date, &relevantOrder.Status)
		if err != nil {
			panic(err)
		}

		orderDetails, err := db.Query(getOrderDetailQuery, relevantOrder.ID)
		if err != nil {
			panic(err)
		}

		for orderDetails.Next() {
			var relevantOrderDetail OrderDetail
			err := orderDetails.Scan(&relevantOrderDetail.OrderID, &relevantOrderDetail.FoodID, &relevantOrderDetail.FoodName, &relevantOrderDetail.Quantity)
			if err != nil {
				panic(err)
			}
			relevantOrder.Items = append(relevantOrder.Items, relevantOrderDetail)
		}

		orderList = append(orderList, relevantOrder)
	}

	return orderList
}

func ChangeOrderStatusToDone(currentSelectedOrder Order) {
	db := OpenDBConnection()
	defer db.Close()

	updateStatusQuery := "UPDATE orders SET status = 'DONE' WHERE id = $1;"
	_, err := db.Exec(updateStatusQuery, currentSelectedOrder.ID)
	if err != nil {
		panic(err)
	}
}

func ChangeOrderStatusToCancelled(currentSelectedOrder Order) {
	db := OpenDBConnection()
	defer db.Close()

	updateStatusQuery := "UPDATE orders SET status = 'CANCELLED' WHERE id = $1;"
	_, err := db.Exec(updateStatusQuery, currentSelectedOrder.ID)
	if err != nil {
		panic(err)
	}
}
