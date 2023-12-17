package db

import "time"

func CreateNewOrder(cartFoods []Food, cartFoodMap map[int]int) {
	db := OpenDBConnection()
	defer db.Close()

	currentDate := time.Now()
	orderId := 0

	createNewOrderQuery := "INSERT INTO orders (date) VALUES ($1) RETURNING id;"
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
