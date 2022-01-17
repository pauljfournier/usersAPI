package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"test/utils"
	"testing"
	"time"
)

var testUsersStore *UsersStore

func TestMain(m *testing.M) {
	//connect to the database
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.TODO()
	err = client.Connect(ctx)

	defer client.Disconnect(ctx)
	testDb := client.Database("testDb")

	dbConnection := utils.DbConnection{
		Client:   client,
		Database: testDb,
		Ctx:      ctx,
	}

	testUsersStore, err = NewUsersStore(dbConnection.Database, dbConnection.Ctx)
	if err != nil {
		log.Fatal(err)
	}

	//run the tests
	exitVal := m.Run()

	//DROP the db to clean
	err = testDb.Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitVal)
}

type storeCreateTest struct {
	user        User
	expectedErr bool
}

func TestStoreCreate(t *testing.T) {
	presetId := "61e41ed578752c5997718aff"

	presetIdUser := storeCreateTest{
		user: User{
			ID:        presetId,
			FirstName: "FirstName",
			LastName:  "LastName",
			Nickname:  "presetNickname",
			Password:  "Password",
			Email:     "presetEmail@email.com",
			Country:   "Country",
			CreatedAt: time.Now().Add(time.Duration(-5) * time.Minute),
		},
		expectedErr: false,
	}

	existingUser := User{
		FirstName: "XFirstName",
		LastName:  "XLastName",
		Nickname:  "XNickname",
		Password:  "XPassword",
		Email:     "XEmail@email.com",
		Country:   "XCountry",
	}
	resultExistingErr := testUsersStore.Create(&existingUser)
	if resultExistingErr != nil {
		t.Errorf("Create user failled for creating a pre-existing user with err %v", resultExistingErr)
	}

	storeCreateTests := []storeCreateTest{
		//test normal behavior
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			expectedErr: false,
		},
		{
			user: User{
				FirstName: "FirstName2",
				LastName:  "LastName3",
				Nickname:  "Nick_name4",
				Password:  "Password5$",
				Email:     "Email2@email.com",
				Country:   "Count$ry",
			},
			expectedErr: false,
		},
		//test with missing field
		{
			user: User{
				LastName: "LastName",
				Nickname: "Nickname3",
				Password: "Password",
				Email:    "Email3@email.com",
				Country:  "Country",
			},
			expectedErr: true,
		},
		{
			user: User{
				FirstName: "FirstName",
				Nickname:  "Nickname4",
				Password:  "Password",
				Email:     "Email4@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Password:  "Password",
				Email:     "Email5@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname6",
				Email:     "Email6@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname7",
				Password:  "Password",
				Country:   "Country",
			},
			expectedErr: true,
		},
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname8",
				Password:  "Password",
				Email:     "Email8@email.com",
			},
			expectedErr: true,
		},
		{
			user: User{
				FirstName: "FirstName",
				Email:     "Email9@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		//test with invalid email
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname10",
				Password:  "Password",
				Email:     "Emailemailcom",
				Country:   "Country",
			},
			expectedErr: true,
		},
		//test with already existing email
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname11",
				Password:  "Password",
				Email:     "XEmail@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		//test with already existing nickname
		{
			user: User{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "XNickname",
				Password:  "Password",
				Email:     "Email12@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		//test with pre-set id and created_at
		presetIdUser,
	}

	for _, item := range storeCreateTests {
		resultErr := testUsersStore.Create(&item.user)
		if !item.expectedErr && resultErr != nil {
			t.Errorf("usersStore.Create for %v output err %v not expected", item.user, resultErr.Error())
		}
		if item.expectedErr && resultErr == nil {
			t.Errorf("usersStore.Create for %v output err expected but not found", item.user)
		}
		if item.user.CreatedAt.After(time.Now()) || item.user.CreatedAt.Before(time.Now().Add(time.Duration(-1)*time.Minute)) {
			t.Errorf("usersStore.Create for %v output wrong CreateAt.", item.user)
		}
		if item.user.UpdatedAt.After(time.Now()) || item.user.UpdatedAt.Before(time.Now().Add(time.Duration(-1)*time.Minute)) {
			t.Errorf("usersStore.Create for %v output wrong UpdatedAt.", item.user)
		}
		//verify the presetIdUser is not set with the presetId
		if item.user.ID == presetId {
			primPresetId, err := primitive.ObjectIDFromHex(presetId)
			if err != nil {
				t.Errorf(err.Error())
			}
			var presetIdFoundUser User
			err = testUsersStore.collection.FindOne(
				testUsersStore.ctx,
				bson.M{"_id": primPresetId},
			).Decode(&presetIdFoundUser)
			if err != mongo.ErrNoDocuments {
				t.Errorf("usersStore.Create for %v (expecting that the given id was not used) output err %v but %v was expected.", primPresetId, err, mongo.ErrNoDocuments)
			}
		}
		//Delete the entry from the db to clean
		if item.expectedErr == false {
			primId, err := primitive.ObjectIDFromHex(item.user.ID)
			if err != nil {
				t.Errorf(err.Error())
			}
			_, err = testUsersStore.collection.DeleteOne(
				testUsersStore.ctx,
				bson.M{"_id": primId},
			)
			if err != nil {
				t.Errorf(err.Error())
			}
		}
	}

	//delete pre-existing user
	primId, err := primitive.ObjectIDFromHex(existingUser.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = testUsersStore.collection.DeleteOne(
		testUsersStore.ctx,
		bson.M{"_id": primId},
	)
	if err != nil {
		t.Errorf("Failled to delete pre-existing user with err %v", err)
	}
}

type storeUpdateTest struct {
	userCreated  []*User //the id of the first user is the one used if present
	userUpdate   User
	userExpected User
	expectedErr  bool
}

func TestStoreUpdate(t *testing.T) {
	presetId := "61e41ed578752c5997718aff"

	presetIdUser := storeUpdateTest{
		userCreated: []*User{
			{
				FirstName: "FirstName",
				LastName:  "LastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
		},
		userUpdate: User{
			ID:        presetId,
			FirstName: "FirstName",
			LastName:  "UpdatedLastName",
			Nickname:  "Nickname",
			Password:  "Password",
			Email:     "Email@email.com",
			Country:   "Country",
			CreatedAt: time.Now().Add(time.Duration(-5) * time.Minute),
		},
		userExpected: User{
			FirstName: "FirstName",
			LastName:  "UpdatedLastName",
			Nickname:  "Nickname",
			Password:  "Password",
			Email:     "Email@email.com",
			Country:   "Country",
		},
		expectedErr: false,
	}

	storeUpdateTests := []storeUpdateTest{
		//test normal behavior
		{
			userCreated: []*User{
				{
					FirstName: "FirstName",
					LastName:  "LastName",
					Nickname:  "Nickname",
					Password:  "Password",
					Email:     "Email@email.com",
					Country:   "Country",
				},
			},
			userUpdate: User{
				FirstName: "FirstName",
				LastName:  "UpdatedLastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			userExpected: User{
				FirstName: "FirstName",
				LastName:  "UpdatedLastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			expectedErr: false,
		},
		{
			userCreated: []*User{
				{
					FirstName: "FirstName",
					LastName:  "LastName",
					Nickname:  "Nickname",
					Password:  "Password",
					Email:     "Email@email.com",
					Country:   "Country",
				},
			},
			userUpdate: User{
				FirstName: "FirstName",
				LastName:  "UpdatedLastName",
				Nickname:  "Nickname",
				Password:  "Updated@Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			userExpected: User{
				FirstName: "FirstName",
				LastName:  "UpdatedLastName",
				Nickname:  "Nickname",
				Password:  "Updated@Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			expectedErr: false,
		},
		//test with missing field
		{
			userCreated: []*User{
				{
					FirstName: "FirstName",
					LastName:  "LastName",
					Nickname:  "Nickname",
					Password:  "Password",
					Email:     "Email@email.com",
					Country:   "Country",
				},
			},
			userUpdate: User{
				LastName: "UpdatedLastName",
				Nickname: "Nickname",
				Password: "Password",
				Email:    "Email@email.com",
				Country:  "Country",
			},
			expectedErr: true,
		},
		{
			userCreated: []*User{
				{
					FirstName: "FirstName",
					LastName:  "LastName",
					Nickname:  "Nickname",
					Password:  "Password",
					Email:     "Email@email.com",
					Country:   "Country",
				},
			},
			userUpdate: User{
				FirstName: "UpdatedFirstName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		{
			userCreated: []*User{
				{
					FirstName: "FirstName",
					LastName:  "LastName",
					Nickname:  "Nickname",
					Password:  "Password",
					Email:     "Email@email.com",
					Country:   "Country",
				},
			},
			userUpdate: User{
				FirstName: "FirstName",
				LastName:  "UpdatedLastName",
				Password:  "Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		{
			userCreated: []*User{
				{
					FirstName: "FirstName",
					LastName:  "LastName",
					Nickname:  "Nickname",
					Password:  "Password",
					Email:     "Email@email.com",
					Country:   "Country",
				},
			},
			userUpdate: User{
				FirstName: "FirstName",
				LastName:  "UpdatedLastName",
				Nickname:  "Nickname",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		{
			userCreated: []*User{
				{
					FirstName: "FirstName",
					LastName:  "LastName",
					Nickname:  "Nickname",
					Password:  "Password",
					Email:     "Email@email.com",
					Country:   "Country",
				},
			},
			userUpdate: User{
				FirstName: "FirstName",
				LastName:  "UpdatedLastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Country:   "Country",
			},
			expectedErr: true,
		},
		{
			userCreated: []*User{
				{
					FirstName: "FirstName",
					LastName:  "LastName",
					Nickname:  "Nickname",
					Password:  "Password",
					Email:     "Email@email.com",
					Country:   "Country",
				},
			},
			userUpdate: User{
				FirstName: "FirstName",
				LastName:  "UpdatedLastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "Email@email.com",
			},
			expectedErr: true,
		},
		{

			userCreated: []*User{
				{
					FirstName: "FirstName",
					LastName:  "LastName",
					Nickname:  "Nickname",
					Password:  "Password",
					Email:     "Email@email.com",
					Country:   "Country",
				},
			},
			userUpdate: User{
				FirstName: "FirstName",
				LastName:  "UpdatedLastName",
				Country:   "Country",
			},
			expectedErr: true,
		},
		//test with invalid email
		{

			userCreated: []*User{
				{
					FirstName: "FirstName",
					LastName:  "LastName",
					Nickname:  "Nickname",
					Password:  "Password",
					Email:     "Email@email.com",
					Country:   "Country",
				},
			},
			userUpdate: User{
				FirstName: "FirstName",
				LastName:  "UpdatedLastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "Emailemailcom",
				Country:   "Country",
			},
			expectedErr: true,
		},
		//test with already existing email
		{
			userCreated: []*User{
				{
					FirstName: "FirstName",
					LastName:  "LastName",
					Nickname:  "Nickname",
					Password:  "Password",
					Email:     "Email@email.com",
					Country:   "Country",
				},
				{
					FirstName: "FirstName2",
					LastName:  "LastName2",
					Nickname:  "Nickname2",
					Password:  "Password2",
					Email:     "Email2@email.com",
					Country:   "Country2",
				},
			},
			userUpdate: User{
				FirstName: "FirstName",
				LastName:  "UpdatedLastName",
				Nickname:  "Nickname",
				Password:  "Password",
				Email:     "Email2@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		//test with already existing nickname
		{
			userCreated: []*User{
				{
					FirstName: "FirstName",
					LastName:  "LastName",
					Nickname:  "Nickname",
					Password:  "Password",
					Email:     "Email@email.com",
					Country:   "Country",
				},
				{
					FirstName: "FirstName2",
					LastName:  "LastName2",
					Nickname:  "Nickname2",
					Password:  "Password2",
					Email:     "Email2@email.com",
					Country:   "Country2",
				},
			},
			userUpdate: User{
				FirstName: "FirstName",
				LastName:  "UpdatedLastName",
				Nickname:  "Nickname2",
				Password:  "Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		//test with non existing User
		{
			userCreated: []*User{},
			userUpdate: User{
				FirstName: "FirstName",
				LastName:  "UpdatedLastName",
				Nickname:  "Nickname",
				Password:  "Updated@Password",
				Email:     "Email@email.com",
				Country:   "Country",
			},
			expectedErr: true,
		},
		//test with pre-set id and created_at
		presetIdUser,
	}

	for itemIndex, item := range storeUpdateTests {
		var id = "61e41ed578752c5997718aee" //possible but non-existing id
		for index, itemToCreate := range item.userCreated {
			resultErr := testUsersStore.Create(itemToCreate)
			if resultErr != nil {
				t.Errorf("Create user failled for update test of index %v item %v with err %v", itemIndex, item, resultErr)
			}
			if index == 0 {
				id = itemToCreate.ID
			}
		}
		resultErr := testUsersStore.Update(id, &item.userUpdate)
		if !item.expectedErr {
			if resultErr != nil {
				t.Errorf("usersStore.Update for %v output err %v not expected", item.userExpected, resultErr.Error())
			} else {
				if !item.userUpdate.IsSoftEqual(&item.userExpected) {
					t.Errorf("usersStore.Update for %v output %v not expected", item.userExpected, item.userUpdate)
				}
			}
		}
		if item.expectedErr && resultErr == nil {
			t.Errorf("usersStore.Update for %v output err expected but not found", item.userUpdate)
		}
		//verify the presetIdUser is not set with the presetId
		if item.userUpdate.ID == presetId {
			primPresetId, err := primitive.ObjectIDFromHex(presetId)
			if err != nil {
				t.Errorf(err.Error())
			}
			var presetIdFoundUser User
			err = testUsersStore.collection.FindOne(
				testUsersStore.ctx,
				bson.M{"_id": primPresetId},
			).Decode(&presetIdFoundUser)
			if err != mongo.ErrNoDocuments {
				t.Errorf("usersStore.Update for %v (expecting that the given id was not used) output err %v but %v was expected.", primPresetId, err, mongo.ErrNoDocuments)
			}
		}
		//Delete the entry from the db to clean
		for _, itemToCreate := range item.userCreated {
			primId, err := primitive.ObjectIDFromHex(itemToCreate.ID)
			if err != nil {
				t.Errorf(err.Error())
			}
			_, err = testUsersStore.collection.DeleteOne(
				testUsersStore.ctx,
				bson.M{"_id": primId},
			)
			if err != nil {
				t.Errorf(err.Error())
			}
		}
	}
}

func TestStoreDelete(t *testing.T) {
	user := User{
		FirstName: "FirstName",
		LastName:  "LastName",
		Nickname:  "Nickname",
		Password:  "Password",
		Email:     "Email@email.com",
		Country:   "Country",
	}
	//insert the user to be deleted
	resultErr := testUsersStore.Create(&user)
	if resultErr != nil {
		t.Errorf("Create user failled for delete test of item %v with err %v", user, resultErr)
	}

	fakeId := "61e41ed578752c5997718aff"

	//normal behavior
	err := testUsersStore.Delete(user.ID)
	if err != nil {
		t.Errorf("usersStore.Delete failled with err %v", err)
	}

	//test non-existing id
	err = testUsersStore.Delete(fakeId)
	if err == nil {
		t.Errorf("usersStore.Delete did not fail with fake id.")
	}
}

type listParameters struct {
	text, id, firstName, lastname, nickname, password, email, country  string
	startDateCreated, endDateCreated, startDateUpdated, endDateUpdated time.Time
	page, pageSize                                                     int64
}

type storeListTest struct {
	parameters    listParameters
	usersExpected []*User
	expectedErr   bool
}

func TestStoreList(t *testing.T) {
	usersInit := []*User{
		{
			FirstName: "6FirstName",
			LastName:  "6LastName",
			Nickname:  "6Nickname",
			Password:  "6Password",
			Email:     "6Email@email.com",
			Country:   "6Country",
		},
		{
			FirstName: "5FirstName",
			LastName:  "5LastName",
			Nickname:  "5Nickname",
			Password:  "5Password",
			Email:     "5Email@email.com",
			Country:   "5Country",
		},
		{
			FirstName: "4FirstName",
			LastName:  "4LastName",
			Nickname:  "4Nickname",
			Password:  "4Password",
			Email:     "4Email@email.com",
			Country:   "4Country",
		},
		{
			FirstName: "3FirstName",
			LastName:  "3LastName",
			Nickname:  "3Nickname",
			Password:  "3Password",
			Email:     "3Email@email.com",
			Country:   "3Country",
		},
		{
			FirstName: "2FirstName",
			LastName:  "2LastName",
			Nickname:  "2Nickname",
			Password:  "2Password",
			Email:     "2Email@email.com",
			Country:   "2Country",
		},
		{
			FirstName: "1FirstName",
			LastName:  "12LastName",
			Nickname:  "1Nickname",
			Password:  "1Password",
			Email:     "1Email@email.com",
			Country:   "1Country",
		},
		{
			FirstName: "0FirstName",
			LastName:  "0LastName",
			Nickname:  "0Nickname",
			Password:  "0Password",
			Email:     "0Email@email.com",
			Country:   "0Country",
		},
	}

	//insert for test
	for index, itemToCreate := range usersInit {
		resultErr := testUsersStore.Create(itemToCreate)
		if resultErr != nil {
			t.Errorf("List user failled to create for test of index %v item %v with err %v", index, itemToCreate, resultErr)
		}
		time.Sleep(10 * time.Millisecond)
	}

	storeListTests := []storeListTest{
		//normal return all
		{
			parameters:    listParameters{},
			usersExpected: usersInit,
			expectedErr:   false,
		},
		//normal with page and page_size
		{
			parameters: listParameters{page: 1, pageSize: 2},
			usersExpected: []*User{
				usersInit[len(usersInit)-3],
				usersInit[len(usersInit)-4],
			},
			expectedErr: false,
		},
		//partial page
		{
			parameters: listParameters{page: 2, pageSize: 3},
			usersExpected: []*User{
				usersInit[len(usersInit)-7],
			},
			expectedErr: false,
		},
		//no result page too far
		{
			parameters:    listParameters{page: 3, pageSize: 3},
			usersExpected: []*User{},
			expectedErr:   false,
		},
		//text in multiple field
		{
			parameters: listParameters{text: "2", page: 0, pageSize: 3},
			usersExpected: []*User{
				usersInit[len(usersInit)-3],
				usersInit[len(usersInit)-2],
			},
			expectedErr: false,
		},
		{
			parameters: listParameters{text: "0P", page: 0, pageSize: 3},
			usersExpected: []*User{
				usersInit[len(usersInit)-1],
			},
			expectedErr: false,
		},
		//various parameters
		{
			parameters: listParameters{country: "0Country", page: 0, pageSize: 3},
			usersExpected: []*User{
				usersInit[len(usersInit)-1],
			},
			expectedErr: false,
		},
		{
			parameters: listParameters{email: "1Email@email.com", page: 0, pageSize: 3},
			usersExpected: []*User{
				usersInit[len(usersInit)-2],
			},
			expectedErr: false,
		},
		{
			parameters: listParameters{password: "2Password", page: 0, pageSize: 3},
			usersExpected: []*User{
				usersInit[len(usersInit)-3],
			},
			expectedErr: false,
		},
		{
			parameters: listParameters{nickname: "3Nickname", page: 0, pageSize: 3},
			usersExpected: []*User{
				usersInit[len(usersInit)-4],
			},
			expectedErr: false,
		},
		{
			parameters: listParameters{lastname: "4LastName", page: 0, pageSize: 3},
			usersExpected: []*User{
				usersInit[len(usersInit)-5],
			},
			expectedErr: false,
		},
		{
			parameters: listParameters{firstName: "5FirstName", page: 0, pageSize: 3},
			usersExpected: []*User{
				usersInit[len(usersInit)-6],
			},
			expectedErr: false,
		},
		//dates
		{
			parameters: listParameters{startDateCreated: usersInit[len(usersInit)-2].CreatedAt, page: 0, pageSize: 3},
			usersExpected: []*User{
				usersInit[len(usersInit)-2],
				usersInit[len(usersInit)-1],
			},
			expectedErr: false,
		},
		{
			parameters: listParameters{endDateCreated: usersInit[1].CreatedAt, page: 0, pageSize: 3},
			usersExpected: []*User{
				usersInit[1],
				usersInit[0],
			},
			expectedErr: false,
		},
		{
			parameters: listParameters{startDateUpdated: usersInit[len(usersInit)-2].UpdatedAt, page: 0, pageSize: 3},
			usersExpected: []*User{
				usersInit[len(usersInit)-2],
				usersInit[len(usersInit)-1],
			},
			expectedErr: false,
		},
		{
			parameters: listParameters{endDateUpdated: usersInit[1].UpdatedAt, page: 0, pageSize: 3},
			usersExpected: []*User{
				usersInit[1],
				usersInit[0],
			},
			expectedErr: false,
		},
		// complex query
		{
			parameters: listParameters{startDateUpdated: usersInit[1].UpdatedAt, text: "12", country: "oun", page: 0, pageSize: 3},
			usersExpected: []*User{
				usersInit[len(usersInit)-2],
			},
			expectedErr: false,
		},
		{
			parameters:    listParameters{startDateUpdated: usersInit[1].UpdatedAt, text: "12", country: "oun", firstName: "boris", page: 0, pageSize: 3},
			usersExpected: []*User{},
			expectedErr:   false,
		},
	}

	for _, item := range storeListTests {
		resultUsers, _, resultErr := testUsersStore.List(
			item.parameters.text,
			item.parameters.id,
			item.parameters.firstName,
			item.parameters.lastname,
			item.parameters.nickname,
			item.parameters.password,
			item.parameters.email,
			item.parameters.country,
			item.parameters.startDateCreated,
			item.parameters.endDateCreated,
			item.parameters.startDateUpdated,
			item.parameters.endDateUpdated,
			item.parameters.page,
			item.parameters.pageSize,
		)
		if !item.expectedErr {
			if resultErr != nil {
				t.Errorf("usersStore.List for %v output err %v not expected", item.parameters, resultErr.Error())
			} else {
				if !areSameUsers(item.usersExpected, resultUsers) {
					var usersExpected []User
					for _, expected := range item.usersExpected {
						usersExpected = append(usersExpected, *expected)
					}
					t.Errorf("usersStore.List for %v output %v but expected %v", item.parameters, resultUsers, usersExpected)

				}
			}
		}
		if item.expectedErr && resultErr == nil {
			t.Errorf("usersStore.List for %v output err expected but not found", item.parameters)
		}
	}

	//delete to clean
	for index, itemToCreate := range usersInit {
		resultErr := testUsersStore.Delete(itemToCreate.ID)
		if resultErr != nil {
			t.Errorf("List user failled to delete for test of index %v item %v with err %v", index, itemToCreate, resultErr)
		}
	}
}

func areSameUsers(usersA []*User, usersB []User) bool {
	if len(usersA) != len(usersB) {
		return false
	}
	for i := range usersA {
		found := false
		for j := range usersB {
			if usersA[i].IsSoftEqual(&usersB[j]) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
