package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/asdine/storm"

	"../user"
	"gopkg.in/mgo.v2/bson"
)

func usersGetAll(w http.ResponseWriter, r *http.Request) {
	users, err := user.All()
	if err != nil {
		postError(w, http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodHead {
		postBodyResponse(w, http.StatusOK, jsonResponse{})
		return
	}
	postBodyResponse(w, http.StatusOK, jsonResponse{"users": users})
}

func usersGetOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	users, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			postError(w, http.StatusNotFound)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodHead {
		postBodyResponse(w, http.StatusOK, jsonResponse{})
		return
	}
	postBodyResponse(w, http.StatusOK, jsonResponse{"users": users})
}

func createUserBody(r *http.Request, u *user.User) error {
	if r == nil  || r.Body == nil {
		return errors.New("Request body is empty")
	}
	if u == nil {
		return errors.New("A user is req")
	}

	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bd, u)
}

func usersPostOne(w http.ResponseWriter, r *http.Request) {
	u := new(user.User)
	err := createUserBody(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	u.ID = bson.NewObjectId()
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
		} else {
			postError(w, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Location", "/users/"+u.ID.Hex())
	w.WriteHeader(http.StatusCreated)

}

func usersPutOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	u := new(user.User)
	err := createUserBody(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
		} else {
			postError(w, http.StatusInternalServerError)
		}
		return
	}
	postBodyResponse(w, http.StatusOK, jsonResponse{"users": u})

}

func usersPatchOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			postError(w, http.StatusNotFound)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	err = createUserBody(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
		} else {
			postError(w, http.StatusInternalServerError)
		}
		return
	}
	postBodyResponse(w, http.StatusOK, jsonResponse{"users": u})

}

func usersDeleteOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	err := user.Delete(id)
	if err != nil {
		if err == storm.ErrNotFound {
			postError(w, http.StatusNotFound)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
