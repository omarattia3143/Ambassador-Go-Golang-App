package main

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/models"
	"github.com/ddosify/go-faker/faker"
	"math/rand"
)

func main() {
	database.Connect()
	faker := faker.NewFaker()
	for i := 0; i < 30; i++ {
		var orderItems []models.OrderItem

		for j := 0; j < rand.Intn(5); j++ {
			price := float64(rand.Intn(90) + 10)
			orderItems = append(orderItems, models.OrderItem{
				ProductTitle:      faker.RandomProductName(),
				Price:             price,
				Quantity:          uint(rand.Intn(5) + 1),
				AdminRevenue:      0.9 * price,
				AmbassadorRevenue: 0.1 * price,
			})
		}

		database.DB.Create(&models.Order{
			UserId:          uint(rand.Intn(30) + 1),
			Code:            faker.RandomBankAccountBic(),
			AmbassadorEmail: faker.RandomEmail(),
			FirstName:       faker.RandomPersonFirstName(),
			LastName:        faker.RandomPersonLastName(),
			Email:           faker.RandomEmail(),
			Complete:        true,
			OrderItems:      orderItems,
		})
	}

}
