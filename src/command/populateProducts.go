package main

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/models"
	"github.com/ddosify/go-faker/faker"
	"math/rand"
)

func main() {
	faker := faker.NewFaker()
	database.Connect()

	for i := 0; i < 30; i++ {

		products := models.Product{
			Title:       faker.RandomProductName(),
			Description: faker.RandomProductAdjective(),
			Image:       faker.RandomImageURL(),
			Price:       float64(rand.Intn(90) + 10),
		}

		database.DB.Create(&products)
	}
}
