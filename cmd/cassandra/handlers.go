package main

import (
	"github.com/unrolled/render"
	"errors"
	"net/http"
	"github.com/gorilla/mux"
)

//connectivityHandler defining route handler which performs basic connectivity test by reading/writing data to Cassandra
func (plugin *CassandraRestAPIPlugin) tweetsHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		pathParams := mux.Vars(req)

		switch req.Method {
		case "POST":
			err := insertTweets(plugin.broker)

			if err != nil {
				formatter.JSON(w, http.StatusInternalServerError, err.Error())
			} else {
				formatter.JSON(w, http.StatusOK, "Tweets inserted successfully")
			}
		case "GET":
			if pathParams != nil && len(pathParams) > 0 {
				id := pathParams["id"]
				if id != "" {
					result, err := getTweetByID(plugin.broker, id)

					if err != nil {
						formatter.JSON(w, http.StatusInternalServerError, err.Error())
					} else {
						formatter.JSON(w, http.StatusOK, result)
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
		case "PUT":
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
		case "DELETE":
			if pathParams != nil && len(pathParams) > 0 {
				id := pathParams["id"]
				if id != "" {
					err := deleteTweetByID(plugin.broker, id)

					if err != nil {
						formatter.JSON(w, http.StatusInternalServerError, err.Error())
					} else {
						formatter.JSON(w, http.StatusOK, "Tweet deleted successfully")
					}
				} else {
					formatter.JSON(w, http.StatusBadRequest, errors.New("id is nil"))
				}
			} else {
				formatter.JSON(w, http.StatusBadRequest, errors.New("Request not supported"))
			}

		default:
			formatter.JSON(w, http.StatusMethodNotAllowed, nil)

		}
	}
}

//usersHandler defining route handler which indicates use of user defined types
//used to return map of addresses as HTTP response
func (plugin *CassandraRestAPIPlugin) usersHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		pathParams := mux.Vars(req)

		switch req.Method {
		case "POST":
			err := insertUsers(plugin.broker)

			if err != nil {
				formatter.JSON(w, http.StatusInternalServerError, err.Error())
			} else {
				formatter.JSON(w, http.StatusOK, "Inserted users successfully")
			}
		case "GET":
			if pathParams != nil && len(pathParams) > 0 {
				id := pathParams["id"]
				if id != "" {
					result, err := getUserByID(plugin.broker, id)

					if err != nil {
						formatter.JSON(w, http.StatusInternalServerError, err.Error())
					} else {
						formatter.JSON(w, http.StatusOK, result)
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
		default:
			formatter.JSON(w, http.StatusMethodNotAllowed, nil)
		}
	}
}
