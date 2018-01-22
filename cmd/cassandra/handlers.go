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
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/unrolled/render"
)

//tweetsGetHandler defining route handler which reads a tweet or all tweets from cassandra database
func (plugin *CassandraRestAPIPlugin) tweetsGetHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		pathParams := mux.Vars(req)
		if pathParams != nil && len(pathParams) > 0 {
			id := pathParams["id"]
			if id != "" {
				result, err := getTweetByID(plugin.broker, plugin.tweetTable, id)

				if err != nil {
					formatter.JSON(w, http.StatusInternalServerError, err.Error())
				} else {

					if result != nil && result.ID != "" {
						formatter.JSON(w, http.StatusOK, result)
					} else {
						formatter.JSON(w, http.StatusNotFound, "Tweet not found")
					}
				}
			} else {
				formatter.JSON(w, http.StatusBadRequest, errors.New("id is nil"))
			}
		} else {
			result, err := getAllTweets(plugin.broker, plugin.tweetTable)

			if err != nil {
				formatter.JSON(w, http.StatusInternalServerError, err.Error())
			} else {
				formatter.JSON(w, http.StatusOK, result)
			}
		}
	}
}

//tweetsPutHandler defining route handler which writes a tweet to cassandra database
func (plugin *CassandraRestAPIPlugin) tweetsPutHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		pathParams := mux.Vars(req)
		if pathParams != nil && len(pathParams) > 0 {
			id := pathParams["id"]
			if id != "" {
				err := insertTweet(plugin.broker, id)

				if err != nil {
					formatter.JSON(w, http.StatusInternalServerError, err.Error())
				} else {
					formatter.JSON(w, http.StatusCreated, "Tweet inserted successfully")
				}
			} else {
				formatter.JSON(w, http.StatusBadRequest, errors.New("id is nil"))
			}
		} else {
			formatter.JSON(w, http.StatusBadRequest, errors.New("Request not supported"))
		}
	}
}

//tweetsPostHandler defining route handler which performs gets all tweets from cassandra database
func (plugin *CassandraRestAPIPlugin) tweetsPostHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		result, err := getAllTweets(plugin.broker, plugin.tweetTable)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, result)
		}
	}
}

//tweetsDeleteHandler defining route handler which deletes a tweet from cassandra database
func (plugin *CassandraRestAPIPlugin) tweetsDeleteHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		pathParams := mux.Vars(req)
		if pathParams != nil && len(pathParams) > 0 {
			id := pathParams["id"]
			if id != "" {

				tweet, err := getTweetByID(plugin.broker, plugin.tweetTable, id)
				if err != nil {
					formatter.JSON(w, http.StatusInternalServerError, err.Error())
				}

				if tweet != nil && tweet.ID != "" {

					err = deleteTweetByID(plugin.broker, plugin.tweetTable, id)

					if err != nil {
						formatter.JSON(w, http.StatusInternalServerError, err.Error())
					} else {
						formatter.JSON(w, http.StatusOK, "Tweet deleted successfully")
					}
				} else {
					formatter.JSON(w, http.StatusNotFound, "Tweet not found")
				}
			} else {
				formatter.JSON(w, http.StatusBadRequest, errors.New("id is nil"))
			}
		} else {
			formatter.JSON(w, http.StatusBadRequest, errors.New("Request not supported"))
		}
	}
}

//usersGetHandler defining route handler which indicates retrieval of user defined type from cassandra
func (plugin *CassandraRestAPIPlugin) usersGetHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		pathParams := mux.Vars(req)
		if pathParams != nil && len(pathParams) > 0 {
			id := pathParams["id"]
			if id != "" {
				result, err := getUserByID(plugin.broker, plugin.userTable, id)

				if err != nil {
					formatter.JSON(w, http.StatusInternalServerError, err.Error())
				} else {
					if result != nil && result.ID != "" {
						formatter.JSON(w, http.StatusOK, result)
					} else {
						formatter.JSON(w, http.StatusNotFound, "User not found")
					}
				}
			} else {
				formatter.JSON(w, http.StatusBadRequest, errors.New("id is nil"))
			}
		} else {
			result, err := getAllUsers(plugin.broker, plugin.userTable)

			if err != nil {
				formatter.JSON(w, http.StatusInternalServerError, err.Error())
			} else {
				formatter.JSON(w, http.StatusOK, result)
			}
		}
	}
}

//usersPostHandler defining route handler which indicates retrieval of user defined type from cassandra
func (plugin *CassandraRestAPIPlugin) usersPostHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		result, err := getAllUsers(plugin.broker, plugin.userTable)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			formatter.JSON(w, http.StatusOK, result)
		}
	}
}

//usersPutHandler defining route handler which indicates storing a user defined type to cassandra
//used to return map of addresses as HTTP response
func (plugin *CassandraRestAPIPlugin) usersPutHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		pathParams := mux.Vars(req)
		if pathParams != nil && len(pathParams) > 0 {
			id := pathParams["id"]
			if id != "" {
				err := insertUser(plugin.broker, id)

				if err != nil {
					formatter.JSON(w, http.StatusInternalServerError, err.Error())
				} else {
					formatter.JSON(w, http.StatusCreated, "User inserted successfully")
				}
			} else {
				formatter.JSON(w, http.StatusBadRequest, errors.New("id is nil"))
			}
		} else {
			formatter.JSON(w, http.StatusBadRequest, errors.New("Request not supported"))
		}
	}
}
