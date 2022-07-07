package main

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/models"
	faker2 "github.com/ddosify/go-faker/faker"
	"sync"
)

func main() {
	database.Connect()
	faker := faker2.NewFaker()

	var products []models.Product

	database.DB.Find(&products)

	wg := sync.WaitGroup{}

	for _, product := range products {

		go func(p models.Product) {
			wg.Add(1)
			defer wg.Done()

			p.Title = faker.RandomProductName()
			p.Description = faker.RandomProductAdjective()
			p.Image = faker.RandomAvatarImage()

			database.DB.Save(&p)

		}(product)

	}

	wg.Wait()

}
