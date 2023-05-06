package db

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/programmingbunny/epub-backend/configs"
	"github.com/programmingbunny/epub-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/mongo"
)


var UserCollection *mongo.Collection = configs.GetCollection(configs.DB, "OnWord", "UserDetails")

var BookCollection *mongo.Collection = configs.GetCollection(configs.DB, "OnWord", "BookDetails")

var ChapterCollection *mongo.Collection = configs.GetCollection(configs.DB, "OnWord", "Chapters")

var ImageCollection *mongo.Collection = configs.GetCollection(configs.DB, "OnWord", "Images")

var VersionCollection *mongo.Collection = configs.GetCollection(configs.DB, "OnWord", "Versions")

var NoteCollection *mongo.Collection = configs.GetCollection(configs.DB, "OnWord", "Notes")

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
    passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(passwordHash), nil
}

// InsertUser inserts the given user into the UserCollection
func InsertUser(ctx context.Context, newUser models.User) error {
   

    _, err := UserCollection.InsertOne(ctx, newUser)
    if err != nil {
        return err
    }
    return nil
}

// DeleteUserByID deletes a user by their ID from the UserCollection
func DeleteUserByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := UserCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserById(ctx context.Context, id primitive.ObjectID, updatedUser models.User) error {
	validate := validator.New()
	err := validate.Struct(updatedUser)
	if err != nil {
		return err
	}
	// Update the user document
	update := bson.M{
		"$set": updatedUser,
	}
	_, err = UserCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID retrieves a user by their ID from the UserCollection
func GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := UserCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func InsertBook(ctx context.Context, newBook models.Book) (result *mongo.InsertOneResult, err error) {
	result, err = BookCollection.InsertOne(ctx, newBook)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DeleteBookByID(ctx context.Context, id primitive.ObjectID) error {
	// Delete the book from the BookDetails collection
	_, err := BookCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	// Delete all chapters associated with the book ID
	_, err = ChapterCollection.DeleteMany(ctx, bson.M{"bookID": id})
	if err != nil {
		return err
	}

	// Delete all images associated with the book ID
	_, err = ImageCollection.DeleteMany(ctx, bson.M{"bookID": id})
	if err != nil {
		return err
	}

	// Delete all versions associated with the book ID
	_, err = VersionCollection.DeleteMany(ctx, bson.M{"bookID": id})
	if err != nil {
		return err
	}

	// Delete all notes associated with the book ID
	_, err = NoteCollection.DeleteMany(ctx, bson.M{"bookID": id})
	if err != nil {
		return err
	}

	return nil
}

func DeleteNoteByID(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	// Delete the note from the database
	result, err := NoteCollection.DeleteOne(ctx, bson.M{"_id": id})
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

func DeleteChapterByID(ctx context.Context, chapterID primitive.ObjectID) (*mongo.DeleteResult, error) {
    // Delete the chapter from the database
    result, err := ChapterCollection.DeleteOne(ctx, bson.M{"_id": chapterID})
    if err != nil {
        return nil, err
    }

    return result, nil
}

func UpdateChapterTitleByID(ctx context.Context, bookID primitive.ObjectID, chNum int, title string) (*mongo.UpdateResult, error) {
	filter := bson.M{"bookId": bookID, "chNum": chNum}
	update := bson.M{"$set": bson.M{
		"title": title,
	}}

	result, err := ChapterCollection.UpdateOne(ctx, filter, update)
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

func InsertVersion(ctx context.Context, newVersion models.Version) (result *mongo.InsertOneResult, err error) {
	result, err = VersionCollection.InsertOne(ctx, newVersion)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func InsertNotes(ctx context.Context, newNote models.Notes) (result *mongo.InsertOneResult, err error) {
	result, err = NoteCollection.InsertOne(ctx, newNote)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func UpdateNoteByID(ctx context.Context, id primitive.ObjectID, updatedNote models.Notes) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"title":     updatedNote.Title,
		"text":      updatedNote.Text,
		"type":      updatedNote.Type,
		"bookID":    updatedNote.BookID,
		"versionID": updatedNote.VersionID,
	}}

	result, err := NoteCollection.UpdateOne(ctx, filter, update)
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
