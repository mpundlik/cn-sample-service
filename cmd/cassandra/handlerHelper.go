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
)

//TweetTable used to represent the tweet table
var TweetTable = &tweet{}

//UserTable used to represent the user table
var UserTable = &user{}

//insertTweet used to handle PUT to insert a tweet in cassandra database
func insertTweet(db sql.Broker, id string) (err error) {

	var insertTweet1 = &tweet{ID: id, Timeline: "me1", Text: "hello world1", User: "user1"}
	err = db.Put(sql.FieldEQ(&insertTweet1.ID), insertTweet1)
	if err != nil {
		return err
	}

	return nil
}

//getAllTweets used to handle GET to get all tweets from cassandra database
func getAllTweets(db sql.Broker) (result *[]tweetResource, err error) {

	var tweets = &[]tweetResource{}

	query2 := sql.FROM(TweetTable, nil)
	err = sql.SliceIt(tweets, db.ListValues(query2))

	if err != nil {
		return nil, err
	}

	return tweets, nil
}

//getTweetByID used to handle GET to get a tweet by ID from cassandra database
func getTweetByID(db sql.Broker, id string) (result *tweetResource, err error) {

	tweet := &tweetResource{}

	query1 := sql.FROM(TweetTable, sql.WHERE(sql.Field(&TweetTable.ID, sql.EQ(id))))
	_, err = db.GetValue(query1, tweet)

	if err != nil {
		return nil, err
	}

	return tweet, nil
}

//deleteTweetByID used to handle DELETE to delete a tweet from cassandra database
func deleteTweetByID(db sql.Broker, id string) (err error) {

	err = db.Delete(sql.FROM(TweetTable, sql.WHERE(sql.Field(&TweetTable.ID, sql.EQ(id)))))

	if err != nil {
		return err
	}

	return nil
}

//insertUser used to handle PUT to insert a user in cassandra database
func insertUser(db sql.Broker, id string) (err error) {

	homePhone := phone{CountryCode: 1, Number: "408-123-1234"}
	cellPhone := phone{CountryCode: 1, Number: "408-123-1235"}
	workPhone := phone{CountryCode: 1, Number: "408-123-1236"}

	phoneMap1 := map[string]phone{"home": homePhone, "cell": cellPhone}
	homeAddr := address{City: "San Jose", Street: "123 Tasman Drive", Zip: "95135", Phones: phoneMap1}

	phoneMap2 := map[string]phone{"work": workPhone}
	workAddr := address{City: "San Jose", Street: "255 E Tasman Drive", Zip: "95134", Phones: phoneMap2}

	addressMap := map[string]address{"home": homeAddr, "work": workAddr}

	var insertUser1 = &user{ID: id, Addresses: addressMap}
	err = db.Put(sql.FieldEQ(&insertUser1.ID), insertUser1)
	if err != nil {
		return err
	}

	return nil
}

//getUserByID used to handle GET for retrieving a user-defined type from cassandra database
func getUserByID(db sql.Broker, id string) (result *userResource, err error) {
	user := &userResource{}

	query1 := sql.FROM(UserTable, sql.WHERE(sql.Field(&UserTable.ID, sql.EQ(id))))
	_, err = db.GetValue(query1, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

//getAllUsers used to handle GET for retrieving all entries for a user-defined type from cassandra database
func getAllUsers(db sql.Broker) (result *[]userResource, err error) {
	users := &[]userResource{}

	query1 := sql.FROM(UserTable, sql.Exp(""))
	err = sql.SliceIt(users, db.ListValues(query1))

	if err != nil {
		return nil, err
	}

	return users, nil
}

// SchemaName schema name for tweet table
func (entity *tweet) SchemaName() string {
	return "example"
}

// SchemaName schema name for user table
func (entity *user) SchemaName() string {
	return "example2"
}
