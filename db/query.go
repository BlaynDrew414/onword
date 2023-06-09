package db

import (
	"net/url"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ParseQuery(query url.Values) (bson.M, error) {
	pipe := make([]bson.M, 0)
	for key, values := range query {
		switch key {
		case "bookID":
			in := make([]interface{}, 0)
			for _, value := range values {
				in = append(in, value)
				if id, err := primitive.ObjectIDFromHex(value); err == nil {
					in = append(in, id)
				}
			}
			if len(in) > 0 {
				pipe = append(pipe, bson.M{"bookID": bson.M{"$in": in}})
			}
		case "versionID":
			in := make([]interface{}, 0)
			for _, value := range values {
				in = append(in, value)
				if id, err := primitive.ObjectIDFromHex(value); err == nil {
					in = append(in, id)
				}
			}
			if len(in) > 0 {
				pipe = append(pipe, bson.M{"versionID": bson.M{"$in": in}})
			}
		case "noteID":
            in := make([]interface{}, 0)
            for _, value := range values {
                if primitive.IsValidObjectID(value) {
                    if id, err := primitive.ObjectIDFromHex(value); err == nil {
                        in = append(in, id)
                    }
                } else {
                    in = append(in, value)
                }
            }
            if len(in) > 0 {
                pipe = append(pipe, bson.M{"_id": bson.M{"$in": in}})
            }
		}
	}
	if len(pipe) > 0 {
		return bson.M{"$and": pipe}, nil
	}
	return bson.M{}, nil
}
