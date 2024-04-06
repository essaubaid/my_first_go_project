package controllers

import (
	"github.com/essaubaid/my_first_go_project/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")
