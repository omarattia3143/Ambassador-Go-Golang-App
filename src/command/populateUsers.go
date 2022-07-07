package main

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/models"
	"github.com/ddosify/go-faker/faker"
)

func main() {

	faker := faker.NewFaker()
	database.Connect()

	for i := 0; i < 30; i++ {
		ambassador := models.User{
			FirstName:    faker.RandomPersonFirstName(),
			LastName:     faker.RandomPersonLastName(),
			Email:        faker.RandomEmail(),
			IsAmbassador: true,
		}

		ambassador.SetPassword("1234")
		database.DB.Create(&ambassador)
	}
}
