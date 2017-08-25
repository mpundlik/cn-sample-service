package main

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
)

//tweetsGetHandler defining route handler which reads a tweet or all tweets from cassandra database
func (plugin *CassandraRestAPIPlugin) tweetsGetHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		pathParams := mux.Vars(req)
		if pathParams != nil && len(pathParams) > 0 {
			id := pathParams["id"]
			if id != "" {
				result, err := getTweetByID(plugin.broker, id)

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
			result, err := getAllTweets(plugin.broker)

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
		result, err := getAllTweets(plugin.broker)

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

				tweet, err := getTweetByID(plugin.broker, id)
				if err != nil {
					formatter.JSON(w, http.StatusInternalServerError, err.Error())
				}

				if tweet != nil && tweet.ID != "" {

					err = deleteTweetByID(plugin.broker, id)

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
				result, err := getUserByID(plugin.broker, id)

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
			result, err := getAllUsers(plugin.broker)

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
		result, err := getAllUsers(plugin.broker)

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
