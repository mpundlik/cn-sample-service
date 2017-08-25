// Copyright (c) 2017 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/ligato/cn-infra/db/sql"
	"github.com/satori/go.uuid"
)

//insertTweets used to handle POST to insert tweets in cassandra database
func insertTweets(db sql.Broker) (err error) {

	var insertTweet1 = &tweet{ID: uuid.NewV4().String(), Timeline: "me1", Text: "hello world1", User: "user1"}
	insertErr1 := insert(db, insertTweet1)
	if insertErr1 != nil {
		return insertErr1
	}

	var insertTweet2 = &tweet{ID: uuid.NewV4().String(), Timeline: "me2", Text: "hello world2", User: "user2"}
	insertErr2 := insert(db, insertTweet2)
	if insertErr2 != nil {
		return insertErr2
	}

	return nil
}

//insertTweet used to handle PUT to insert a tweet in cassandra database
func insertTweet(db sql.Broker, id string) (err error) {

	var insertTweet1 = &tweet{ID: id, Timeline: "me1", Text: "hello world1", User: "user1"}
	insertErr1 := insert(db, insertTweet1)
	if insertErr1 != nil {
		return insertErr1
	}

	return nil
}

//getAllTweets used to handle GET to get all tweets from cassandra database
func getAllTweets(db sql.Broker) (result *[]tweet, err error) {

	result, selectAllErr := selectAll(db)
	if selectAllErr != nil {
		return nil, selectAllErr
	}

	return result, nil
}

//getTweetByID used to handle GET to get a tweet by ID from cassandra database
func getTweetByID(db sql.Broker, id string) (result *[]tweet, err error) {

	result, selectErr := selectByID(db, &id)
	if selectErr != nil {
		return nil, selectErr
	}

	return result, nil
}

//deleteTweetByID used to handle DELETE to delete a tweet from cassandra database
func deleteTweetByID(db sql.Broker, id string) (err error) {

	deleteErr := deleteByID(db, &id)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

//insertUsers used to handle POST for storing a user-defined type in cassandra database
func insertUsers(db sql.Broker) (err error) {

	homePhone := phone{CountryCode: 1, Number: "408-123-1234"}
	cellPhone := phone{CountryCode: 1, Number: "408-123-1235"}
	workPhone := phone{CountryCode: 1, Number: "408-123-1236"}

	phoneMap1 := map[string]phone{"home": homePhone, "cell": cellPhone}
	homeAddr := address{City: "San Jose", Street: "123 Tasman Drive", Zip: "95135", Phones: phoneMap1}

	phoneMap2 := map[string]phone{"work": workPhone}
	workAddr := address{City: "San Jose", Street: "255 E Tasman Drive", Zip: "95134", Phones: phoneMap2}

	addressMap := map[string]address{"home": homeAddr, "work": workAddr}

	var insertUser1 = &user{ID: "user1", Addresses: addressMap}
	insertErr1 := insertUserTable(db, insertUser1)
	if insertErr1 != nil {
		return insertErr1
	}

	var insertUser2 = &user{ID: "user2", Addresses: addressMap}
	insertErr2 := insertUserTable(db, insertUser2)
	if insertErr2 != nil {
		return insertErr2
	}

	return nil
}

//getUserByID used to handle GET for retrieving a user-defined type from cassandra database
func getUserByID(db sql.Broker, id string) (result *[]user, err error) {
	var UserTable = &user{}
	users := &[]user{}

	query1 := sql.FROM(UserTable, sql.WHERE(sql.Field(&UserTable.ID, sql.EQ(id))))
	err = sql.SliceIt(users, db.ListValues(query1))

	return users, nil
}

//getAllUsers used to handle GET for retrieving all entries for a user-defined type from cassandra database
func getAllUsers(db sql.Broker) (result *[]user, err error) {
	var UserTable = &user{}
	users := &[]user{}

	query1 := sql.FROM(UserTable, sql.Exp(""))
	err = sql.SliceIt(users, db.ListValues(query1))

	if err != nil {
		return nil, err
	}

	return users, nil
}

//insert used to insert data in tweet table
func insert(db sql.Broker, insertTweet *tweet) (err error) {
	//inserting a record (runs update behind the scene)
	err1 := db.Put(sql.FieldEQ(&insertTweet.ID), insertTweet)
	if err1 != nil {
		return err1
	}

	return nil
}

//insertUserTable used to insert data in user table
func insertUserTable(db sql.Broker, insertUser *user) (err error) {
	//inserting a record (runs update behind the scene)
	err1 := db.Put(sql.FieldEQ(&insertUser.ID), insertUser)
	if err1 != nil {
		return err1
	}
	return nil
}

//selectByID used to retrieve data from tweet table
func selectByID(db sql.Broker, id *string) (result *[]tweet, err error) {
	var TweetTable = &tweet{}
	tweets := &[]tweet{}

	query1 := sql.FROM(TweetTable, sql.WHERE(sql.Field(&TweetTable.ID, sql.EQ(id))))
	err = sql.SliceIt(tweets, db.ListValues(query1))

	if err != nil {
		return nil, err
	}

	return tweets, nil
}

//selectAll used to retrieve all records from tweet table
func selectAll(db sql.Broker) (result *[]tweet, err error) {
	var TweetTable = &tweet{}
	var tweets = &[]tweet{}

	query2 := sql.FROM(TweetTable, nil)
	err = sql.SliceIt(tweets, db.ListValues(query2))

	if err != nil {
		return nil, err
	}

	return tweets, nil
}

//deleteByID used to delete a tweet from tweet table
func deleteByID(db sql.Broker, id *string) (err error) {
	var TweetTable = &tweet{}

	query1 := sql.FROM(TweetTable, sql.WHERE(sql.Field(&TweetTable.ID, sql.EQ(id))))
	err = db.Delete(query1)

	if err != nil {
		return err
	}

	return nil
}

// SchemaName schema name for tweet table
func (entity *tweet) SchemaName() string {
	return "example"
}

// SchemaName schema name for user table
func (entity *user) SchemaName() string {
	return "example2"
}
