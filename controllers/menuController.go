package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/essaubaid/my_first_go_project/database"
	"github.com/essaubaid/my_first_go_project/models"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var validate = validator.New()

/*
Check logic latter for this specific function.
Why is start time in the future.
And check parameter is not used at all
I've removed the check parameter
*/
func inTimeSpan(start, end time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := menuCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while listing the menu items",
			})
			return
		}
		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)

	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		menuId := c.Param("menu_id")
		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error while fetching the menu",
			})
			return
		}

		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if validationErr := validate.Struct(menu); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(
			ctx,
			menu,
		)

		if insertErr != nil {
			msg := "Menu item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}

		updateObj := bson.M{}

		if menu.Start_date == nil || menu.End_date == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "time is missing from the request",
			})
			return
		}

		if !inTimeSpan(*menu.Start_date, *menu.End_date) {
			var msg string = "kindly retype the time"
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})
			return
		}

		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		if menu.Name != "" {
			updateObj["name"] = menu.Name
		}
		if menu.Category != "" {
			updateObj["category"] = menu.Category
		}
		updateObj["start_date"] = menu.Start_date
		updateObj["end_date"] = menu.End_date
		updateObj["updated_at"] = menu.Updated_at

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, upsertError := menuCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if upsertError != nil {
			msg := "Menu update failed"
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": msg,
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			result,
		)
	}
}
