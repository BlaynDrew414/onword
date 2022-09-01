package db

import (
	"context"

	"github.com/programmingbunny/epub-backend/configs"
	"github.com/programmingbunny/epub-backend/models"

	"go.mongodb.org/mongo-driver/mongo"
)

var BookCollection *mongo.Collection = configs.GetCollection(configs.DB, "BookDetails")

var ChapterCollection *mongo.Collection = configs.GetCollection(configs.DB, "Chapters")

var ImageCollection *mongo.Collection = configs.GetCollection(configs.DB, "Images")

func InsertBook(ctx context.Context, newBook models.Book) (result *mongo.InsertOneResult, err error) {
	result, err = BookCollection.InsertOne(ctx, newBook)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func InsertChapter(ctx context.Context, newChapter models.Chapter) (result *mongo.InsertOneResult, err error) {
	result, err = ChapterCollection.InsertOne(ctx, newChapter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func InsertImage(ctx context.Context, newImage models.ChapterImages) (result *mongo.InsertOneResult, err error) {
	result, err = ImageCollection.InsertOne(ctx, newImage)
	if err != nil {
		return nil, err
	}
	return result, nil
}
