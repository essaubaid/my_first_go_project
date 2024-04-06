package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/essaubaid/my_first_go_project/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			recordPerPage = 1
		}

		startIndex := (page - 1) * recordPerPage
		// startIndex, _ = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{
			{Key: "$match", Value: bson.D{}},
		}
		groupStage := bson.D{
			{
				Key: "$group", Value: bson.D{
					{
						Key: "_id", Value: bson.D{
							{Key: "_id", Value: "null"},
						},
					},
					{
						Key: "total_count", Value: bson.D{
							{Key: "$sum", Value: 1},
						},
					},
					{
						Key: "data", Value: bson.D{
							{Key: "$push", Value: "$$ROOT"},
						},
					},
				},
			},
		}
		projectStage := bson.D{
			{
				Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "total_count", Value: 1},
					{
						Key: "food_items", Value: bson.D{
							{
								Key: "$slice", Value: []interface{}{
									"$data",
									startIndex,
									recordPerPage,
								},
							},
						},
					},
				},
			},
		}

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while listing food items",
			})
		}

		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allFoods[0])

	}
}
