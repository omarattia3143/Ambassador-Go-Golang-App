package main

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/models"
	"context"
	"github.com/go-redis/redis/v9"
)

func main() {
	database.Connect()
	database.SetupRedis()

	ctx := context.Background()

	var users []models.User

	database.DB.Find(&users, models.User{
		IsAmbassador: true,
	})

	for _, user := range users {
		ambassador := models.Ambassador(user)

		database.Cache.ZAdd(ctx, "rankings", redis.Z{
			Score:  ambassador.CalculateTotalRevenue(database.DB),
			Member: user.FullName(),
		})
	}
}
