package db

type Food struct {
	ID          int
	Name        string
	Description string
	Price       int
}

func GetAllFoods() []Food {
	var foodList []Food

	db := OpenDBConnection()
	defer db.Close()

	allFoodQuery := "SELECT * FROM foods;"
	rows, err := db.Query(allFoodQuery)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var relevantFood Food
		err := rows.Scan(&relevantFood.ID, &relevantFood.Name, &relevantFood.Description, &relevantFood.Price)
		if err != nil {
			panic(err)
		}

		foodList = append(foodList, relevantFood)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return foodList
}
