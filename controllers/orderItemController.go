package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/essaubaid/my_first_go_project/database"
	"github.com/essaubaid/my_first_go_project/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := orderItemCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while listing the order items",
			})
			return
		}
		var allOrderItems []bson.M
		if err = result.All(ctx, &allOrderItems); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")

		allOrderItems, err := ItemsByOrder(orderId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while listing order items by order ID",
			})
			return
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderItemId := c.Param("order_item_id")
		var orderItem models.OrderItem

		err := orderItemCollection.FindOne(ctx, bson.M{"orderItem_id": orderItemId}).Decode(&orderItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing ordered item"})
			return
		}
		c.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItemPack OrderItemPack
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		orderItemsToBeInserted := []interface{}{}

		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)

		for _, orderItem := range orderItemPack.Order_items {
			orderItem.Order_id = order_id

			if validationError := validate.Struct(orderItem); validationError != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": validationError.Error(),
				})
				return
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()

			num := toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		insertedOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)

		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, insertedOrderItems)

	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItem models.OrderItem

		orderItemId := c.Param("order_item_id")

		filter := bson.M{
			"order_item_id": orderItemId,
		}

		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		updateObj := bson.M{}

		if orderItem.Unit_price != nil {
			updateObj["unit_price"] = *orderItem.Unit_price
		}

		if orderItem.Quantity != nil {
			updateObj["quantity"] = *orderItem.Quantity
		}

		if orderItem.Food_id != nil {
			updateObj["food_id"] = *orderItem.Food_id
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj["updated_at"] = orderItem.Updated_at

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			updateObj,
			&opt,
		)

		if err != nil {
			msg := "Order item update failed"
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})
		}

		c.JSON(http.StatusOK, result)
	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchStage := bson.D{
		{
			Key: "$match", Value: bson.D{
				{Key: "order_id", Value: id},
			},
		},
	}
	lookupStage := bson.D{
		{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "food"},
				{Key: "localField", Value: "food_id"},
				{Key: "foreignField", Value: "food_id"},
				{Key: "as", Value: "food"},
			},
		},
	}
	unwindStage := bson.D{
		{
			Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$food"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			},
		},
	}

	lookupOrderStage := bson.D{
		{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "order"},
				{Key: "localField", Value: "order_id"},
				{Key: "foreignField", Value: "order_id"},
				{Key: "as", Value: "order"},
			},
		},
	}
	unwindOrderStage := bson.D{
		{
			Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$order"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			},
		},
	}

	lookupTableStage := bson.D{
		{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "table"},
				{Key: "localField", Value: "order.table_id"},
				{Key: "foreignField", Value: "table_id"},
				{Key: "as", Value: "table"},
			},
		},
	}
	unwindTableStage := bson.D{
		{
			Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$table"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			},
		},
	}
	projectStage := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "id", Value: 0},
				{Key: "amount", Value: "$food.price"},
				{Key: "total_count", Value: 1},
				{Key: "food_name", Value: "$food.name"},
				{Key: "food_image", Value: "$food.food_image"},
				{Key: "table_number", Value: "$table.table_number"},
				{Key: "table_id", Value: "$table.table_id"},
				{Key: "order_id", Value: "$order.order_id"},
				{Key: "price", Value: "$food.price"},
				{Key: "quantity", Value: 1},
			},
		},
	}
	groupStage := bson.D{
		{
			Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "order_id", Value: "$order_id"},
					{Key: "table_id", Value: "$table_id"},
					{Key: "table_number", Value: "$table_number"},
				}},
				{Key: "payment_due", Value: bson.D{
					{Key: "$sum", Value: "$amount"},
				}},
				{Key: "total_count", Value: bson.D{
					{Key: "$sum", Value: 1},
				}},
				{Key: "order_items", Value: bson.D{
					{Key: "$push", Value: "$$ROOT"},
				}},
			},
		},
	}

	projectStage2 := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "id", Value: 0},
			{Key: "payment_due", Value: 1},
			{Key: "total_count", Value: 1},
			{Key: "table_number", Value: "$_id.table_number"},
			{Key: "order_items", Value: 1},
		}},
	}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})

	if err != nil {
		panic(err)
	}

	if err := result.All(ctx, &OrderItems); err != nil {
		panic(err)
	}

	return OrderItems, err

}
