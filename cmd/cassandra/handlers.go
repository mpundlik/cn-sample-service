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
	"github.com/ligato/cn-infra/db/sql/cassandra"
	"github.com/ligato/cn-infra/logging/logroot"
	"github.com/satori/go.uuid"
	"github.com/smartystreets/assertions"
	"time"
)

//insertTweets used to handle POST to insert tweets in cassandra database
func insertTweets(db *cassandra.BrokerCassa) (err error) {

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

func insertTweet(db *cassandra.BrokerCassa, id string) (err error) {

	var insertTweet1 = &tweet{ID: id, Timeline: "me1", Text: "hello world1", User: "user1"}
	insertErr1 := insert(db, insertTweet1)
	if insertErr1 != nil {
		return insertErr1
	}

	return nil
}

func getAllTweets(db *cassandra.BrokerCassa) (result []*tweet, err error) {

	result, selectAllErr := selectAll(db)
	if selectAllErr != nil {
		return nil, selectAllErr
	}

	return result, nil
}

func getTweetByID(db *cassandra.BrokerCassa, id string) (result *[]tweet, err error) {

	result, selectErr := selectByID(db, &id)
	if selectErr != nil {
		return nil, selectErr
	}

	return result, nil
}

func deleteTweetByID(db *cassandra.BrokerCassa, id string) (err error) {

	deleteErr := deleteByID(db, &id)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

//alterTable used to depict support for ALTER TABLE
func alterTable(db *cassandra.BrokerCassa) (err error) {

	err = db.Exec(`ALTER TABLE example.person ADD data text`)
	if err != nil {
		logroot.StandardLogger().Errorf("Error executing alter table %v", err)
		return err
	}

	var insertPerson = &person{ID: uuid.NewV1().String(), Name: "James Bond", Data: "new column added"}
	insertErr := insertPersonTable(db, insertPerson)
	if insertErr != nil {
		logroot.StandardLogger().Fatalf("Error executing insert %v", err)
		return insertErr
	}

	selectErr := selectPersonByID(db, &insertPerson.ID)
	if selectErr != nil {
		logroot.StandardLogger().Errorf("select error = %v", selectErr)
		return selectErr
	}

	return nil
}

//createKeySpaceIfNotExist used to depict support for IF NOT EXISTS clause while creating a keyspace
func createKeySpaceIfNotExist(db *cassandra.BrokerCassa) (err error) {

	//creating a non existing keyspace
	err1 := db.Exec(`CREATE KEYSPACE IF NOT EXISTS example2 with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`)

	if err1 != nil {
		logroot.StandardLogger().Errorf("Error creating keyspace %v", err1)
		return err1
	}

	//creating a non existing keyspace
	err2 := db.Exec(`CREATE KEYSPACE IF NOT EXISTS example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`)

	if err2 != nil {
		logroot.StandardLogger().Errorf("Error creating keyspace %v", err2)
		return err2
	}

	//does not return error for existing key space, even though key space exists since using 'IF NOT EXISTS'
	err3 := db.Exec(`CREATE KEYSPACE IF NOT EXISTS example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`)

	if err3 != nil {
		logroot.StandardLogger().Errorf("Error creating existing keyspace %v", err3)
		return err3
	}

	//will return an error since key space exists and not using 'IF NOT EXISTS'
	err4 := db.Exec(`CREATE KEYSPACE example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`)

	if err4 != nil {
		logroot.StandardLogger().Errorf("Error %v", err4)
		assertions.ShouldEqual(err4.Error(), "Cannot add existing keyspace \"example\"")
	}

	return nil
}

//insertUserDefinedType used to depict support for customized data structure, creating and using user-defined types and storing/retrieval
func insertUsers(db *cassandra.BrokerCassa) (err error) {

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
		logroot.StandardLogger().Errorf("Insert error = %v", insertErr1)
		return insertErr1
	}

	var insertUser2 = &user{ID: "user2", Addresses: addressMap}
	insertErr2 := insertUserTable(db, insertUser2)
	if insertErr2 != nil {
		logroot.StandardLogger().Errorf("Insert error = %v", insertErr2)
		return insertErr2
	}

	return nil
}

func getUserByID(db *cassandra.BrokerCassa, id string) (result []*user, err error) {
	var UserTable = &user{}
	users := []*user{}

	query1 := sql.FROM(UserTable, sql.WHERE(sql.Field(&UserTable.ID, sql.EQ(id))))
	it := db.ListValues(query1)
	for {
		user := &user{}
		stop := it.GetNext(user)
		if stop {
			break
		}
		users = append(users, user)
	}
	itErr := it.Close()

	if itErr != nil {
		logroot.StandardLogger().Errorf("Error closing iterator %v", itErr)
		return nil, itErr
	}

	logroot.StandardLogger().Infof("users = %v", users)

	return users, nil
}

func getAllUsers(db *cassandra.BrokerCassa) (result []*user, err error) {
	var UserTable = &user{}
	users := []*user{}

	query1 := sql.FROM(UserTable, sql.Exp(""))
	it := db.ListValues(query1)
	for {
		user := &user{}
		stop := it.GetNext(user)
		if stop {
			break
		}
		users = append(users, user)
	}
	itErr := it.Close()

	if itErr != nil {
		logroot.StandardLogger().Errorf("Error closing iterator %v", itErr)
		return nil, itErr
	}

	logroot.StandardLogger().Infof("users = %v", users)

	return users, nil
}

//reconnectInterval used to depict redial_interval timeout behavior
//need to manually bring down cassandra during sleep interval, after bring it back up again we can retrieve results
func reconnectInterval(db *cassandra.BrokerCassa) (err error) {

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

	//selectErr := selectByID(db, &insertTweet1.ID)
	//if selectErr != nil {
	//	return selectErr
	//}

	//sleep for 5 minutes (need to restart cassandra manually in the meantime)
	logroot.StandardLogger().Infof("Sleep for 5 min - start %v", time.Now())
	time.Sleep(5 * time.Minute)
	logroot.StandardLogger().Infof("Sleep for 5 min - end")

	//selectAllErr := selectAll(db)
	//if selectAllErr != nil {
	//	return selectAllErr
	//}

	return nil
}

//queryTimeout used to depict op_timeout timeout behavior
//need to update the config to a very low op_timeout value (600ns) to get the expected timeout error
func queryTimeout(db *cassandra.BrokerCassa) (err error) {

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

	result1, selectErr := selectByID(db, &insertTweet1.ID)
	if selectErr != nil {
		return selectErr
	}
	logroot.StandardLogger().Infof("result1 = %v", result1)

	result2, selectAllErr := selectAll(db)
	if selectAllErr != nil {
		return selectAllErr
	}
	logroot.StandardLogger().Infof("result2 = %v", result2)

	return err
}

//connectTimeout used to depict dial_timeout timeout behavior
//need to update the config to a very low dial_timeout value (600ns) to get the expected timeout error
func connectTimeout(db *cassandra.BrokerCassa) (err error) {

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

	//selectErr := selectByID(db, &insertTweet1.ID)
	//if selectErr != nil {
	//	return selectErr
	//}

	//selectAllErr := selectAll(db)
	//if selectAllErr != nil {
	//	return selectAllErr
	//}

	return err
}

//insert used to insert data in tweet table
func insert(db *cassandra.BrokerCassa, insertTweet *tweet) (err error) {
	//inserting a record (runs update behind the scene)
	start1 := time.Now()
	err1 := db.Put(sql.FieldEQ(&insertTweet.ID), insertTweet)
	if err1 != nil {
		elapsed1 := time.Since(start1)
		logroot.StandardLogger().Infof("Time taken for insert : %v", elapsed1)
		logroot.StandardLogger().Errorf("Error executing insert %v", err1)
		return err1
	}

	elapsed := time.Since(start1)
	logroot.StandardLogger().Infof("Time taken for insert : %v", elapsed)
	return nil
}

//insertPersonTable used to insert data in person table
func insertPersonTable(db *cassandra.BrokerCassa, insertPerson *person) (err error) {
	//inserting a record (runs update behind the scene)
	start1 := time.Now()
	err1 := db.Put(sql.FieldEQ(&insertPerson.ID), insertPerson)
	if err1 != nil {
		elapsed1 := time.Since(start1)
		logroot.StandardLogger().Infof("Time taken for insert : %v", elapsed1)
		logroot.StandardLogger().Errorf("Error executing insert %v", err1)
		return err1
	}

	elapsed := time.Since(start1)
	logroot.StandardLogger().Infof("Time taken for insert : %v", elapsed)
	return nil
}

//insertUserTable used to insert data in user table
func insertUserTable(db *cassandra.BrokerCassa, insertUser *user) (err error) {
	//inserting a record (runs update behind the scene)
	start1 := time.Now()
	err1 := db.Put(sql.FieldEQ(&insertUser.ID), insertUser)
	if err1 != nil {
		elapsed1 := time.Since(start1)
		logroot.StandardLogger().Infof("Time taken for insert : %v", elapsed1)
		logroot.StandardLogger().Errorf("Error executing insert %v", err1)
		return err1
	}

	elapsed := time.Since(start1)
	logroot.StandardLogger().Infof("Time taken for insert : %v", elapsed)
	return nil
}

//selectByID used to retrieve data from tweet table
func selectByID(db *cassandra.BrokerCassa, id *string) (result *[]tweet, err error) {
	start2 := time.Now()
	var TweetTable = &tweet{}
	tweets := &[]tweet{}

	query1 := sql.FROM(TweetTable, sql.WHERE(sql.Field(&TweetTable.ID, sql.EQ(id))))
	err = sql.SliceIt(tweets, db.ListValues(query1))

	if err != nil {
		elapsed2 := time.Since(start2)
		logroot.StandardLogger().Infof("Time taken for select : %v", elapsed2)
		logroot.StandardLogger().Errorf("Error executing select %v", err)
		return nil, err
	}

	elapsed2 := time.Since(start2)
	logroot.StandardLogger().Infof("Time taken for select : %v", elapsed2)
	logroot.StandardLogger().Info("Tweet:", tweets)
	return tweets, nil
}

//selectPersonByID used to retrieve data from person table
func selectPersonByID(db *cassandra.BrokerCassa, id *string) (err error) {
	start2 := time.Now()
	var PersonTable = &person{}
	people := &[]person{}

	query1 := sql.FROM(PersonTable, sql.WHERE(sql.Field(&PersonTable.ID, sql.EQ(id))))
	err = sql.SliceIt(people, db.ListValues(query1))

	if err != nil {
		elapsed2 := time.Since(start2)
		logroot.StandardLogger().Infof("Time taken for select : %v", elapsed2)
		logroot.StandardLogger().Errorf("Error executing select %v", err)
		return err
	}

	elapsed2 := time.Since(start2)
	logroot.StandardLogger().Infof("Time taken for select : %v", elapsed2)
	logroot.StandardLogger().Info("People:", people)
	return nil
}

//selectAll used to retrieve all records from tweet table
func selectAll(db *cassandra.BrokerCassa) (result []*tweet, err error) {
	start3 := time.Now()
	var TweetTable = &tweet{}
	var tweets = []*tweet{}
	query2 := sql.FROM(TweetTable, nil)
	iterator := db.ListValues(query2)
	for {
		tweetItem := &tweet{}
		stop := iterator.GetNext(tweetItem)
		if stop {
			break
		} else {
			logroot.StandardLogger().Info("Tweet Item: ", tweetItem)
			tweets = append(tweets, tweetItem)
		}
	}
	iterator.Close()
	elapsed3 := time.Since(start3)
	logroot.StandardLogger().Infof("Time taken for select all : %v", elapsed3)
	return tweets, nil
}

func deleteByID(db *cassandra.BrokerCassa, id *string) (err error) {
	start2 := time.Now()
	var TweetTable = &tweet{}

	query1 := sql.FROM(TweetTable, sql.WHERE(sql.Field(&TweetTable.ID, sql.EQ(id))))
	err = db.Delete(query1)

	if err != nil {
		elapsed2 := time.Since(start2)
		logroot.StandardLogger().Infof("Time taken for delete : %v", elapsed2)
		logroot.StandardLogger().Errorf("Error executing delete %v", err)
		return err
	}

	elapsed2 := time.Since(start2)
	logroot.StandardLogger().Infof("Time taken for delete : %v", elapsed2)
	return nil
}

// SchemaName schema name for tweet table
func (entity *tweet) SchemaName() string {
	return "example"
}

// SchemaName schema name for person table
func (entity *person) SchemaName() string {
	return "example"
}

// SchemaName schema name for user table
func (entity *user) SchemaName() string {
	return "example2"
}
