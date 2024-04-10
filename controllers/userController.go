package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/essaubaid/my_first_go_project/database"
	"github.com/essaubaid/my_first_go_project/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")
		var user models.User

		if err := userCollection.FindOne(
			ctx, bson.M{"user_id": userId},
		).Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while fetching the user item",
			})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
