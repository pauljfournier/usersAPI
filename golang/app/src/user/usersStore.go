package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// UsersStore implements database operations
type UsersStore struct {
	collection *mongo.Collection
	ctx        context.Context
}

// NewUsersStore returns a UsersStore
func NewUsersStore(db *mongo.Database, ctx context.Context) (*UsersStore, error) {
	usersCollection := db.Collection("users")
	//set the uniq constraint for nickname and email
	_, err := usersCollection.Indexes().CreateMany(
		context.Background(),
		[]mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.D{{Key: "nickname", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return &UsersStore{
		collection: usersCollection,
		ctx:        ctx,
	}, nil
}

// Create creates a new User.
func (s *UsersStore) Create(u *User) error {
	u.ID = ""
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	err := u.Validate()
	if err != nil {
		return err
	}
	userInsertOne, err := s.collection.InsertOne(s.ctx, u)
	if err != nil {
		return err
	}
	userSingleResult := s.collection.FindOne(s.ctx, bson.M{"_id": userInsertOne.InsertedID})
	if userSingleResult.Err() != nil {
		return userSingleResult.Err()
	}

	err = userSingleResult.Decode(u)
	return err
}

// Update update an existing User.
func (s *UsersStore) Update(id string, u *User) error {
	u.UpdatedAt = time.Now()

	err := u.Validate()
	if err != nil {
		return err
	}

	primId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	//Do not allow to directly modify id, created_at and updated_at
	updateResult, err := s.collection.UpdateOne(
		s.ctx,
		bson.M{"_id": primId},
		bson.M{"$set": bson.M{
			"first_name": u.FirstName,
			"last_name":  u.LastName,
			"nickname":   u.Nickname,
			"password":   u.Password,
			"email":      u.Email,
			"country":    u.Country,
			"updated_at": u.UpdatedAt,
		}},
	)
	if err != nil {
		return err
	}
	if updateResult.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	userSingleResult := s.collection.FindOne(s.ctx, bson.M{"_id": primId})
	if userSingleResult.Err() != nil {
		return userSingleResult.Err()
	}

	err = userSingleResult.Decode(u)
	return err
}

// Delete a User from its id.
func (s *UsersStore) Delete(id string) error {

	primId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	deleteResult, err := s.collection.DeleteOne(
		s.ctx,
		bson.M{"_id": primId},
	)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Return a List of User filtered, according to the page and page_size required.
func (s *UsersStore) List(text, id, firstName, lastname, nickname, password, email, country string, startDateCreated, endDateCreated, startDateUpdated, endDateUpdated time.Time, page, pageSize int64) ([]User, int, error) {
	//rmq page start at 0
	skip := pageSize * page
	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &pageSize,
		Sort:  bson.D{{"created_at", -1}},
	}

	var filter []bson.M
	var textFilter []bson.M

	if id != "" {
		primId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, 0, err
		}
		filter = append(
			filter,
			bson.M{"_id": primId},
		)
	}

	addFilterRegex(&filter, "first_name", firstName, &textFilter, text)
	addFilterRegex(&filter, "last_name", lastname, &textFilter, text)
	addFilterRegex(&filter, "nickname", nickname, &textFilter, text)
	addFilterRegex(&filter, "password", password, &textFilter, text)
	addFilterRegex(&filter, "email", email, &textFilter, text)
	addFilterRegex(&filter, "country", country, &textFilter, text)
	addFilterDate(&filter, "created_at", "$gte", startDateCreated)
	addFilterDate(&filter, "created_at", "$lte", endDateCreated)
	addFilterDate(&filter, "updated_at", "$gte", startDateUpdated)
	addFilterDate(&filter, "updated_at", "$lte", endDateUpdated)
	filter = append(
		filter,
		bson.M{
			"$or": textFilter,
		},
	)

	cursor, err := s.collection.Find(
		s.ctx,
		bson.M{
			"$and": filter,
		},
		&opts,
	)
	if err != nil {
		return nil, 0, err
	}
	var uList []User
	err = cursor.All(s.ctx, &uList)
	return uList, len(uList), err
}

func addFilterRegex(filter *[]bson.M, field string, value string, textFilter *[]bson.M, text string) {
	if value != "" {
		*filter = append(
			*filter,
			bson.M{field: bson.D{
				{"$regex", primitive.Regex{Pattern: ".*(" + value + ").*"}},
			}},
		)
	}
	*textFilter = append(
		*textFilter,
		bson.M{field: bson.D{
			{"$regex", primitive.Regex{Pattern: ".*(" + text + ").*"}},
		}},
	)
}

func addFilterDate(filter *[]bson.M, field string, operator string, value time.Time) {
	if !value.IsZero() {
		*filter = append(
			*filter,
			bson.M{field: bson.M{
				operator: primitive.NewDateTimeFromTime(value),
			}},
		)
	}
}
