package db

import (
	"context"
	"time"

	"github.com/programmingbunny/epub-backend/configs"
	"github.com/programmingbunny/epub-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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

func UpdateChapterWithHeaderImage(imageLoc string, book string, chNum int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(book)

	_, err := ChapterCollection.UpdateOne(
		ctx,
		bson.M{"bookID": objId, "chapterNum": chNum},
		bson.M{
			"$set": bson.M{
				"imageLocation": imageLoc}},
	)

	if err != nil {
		return err
	}

	return nil
}
