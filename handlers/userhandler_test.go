package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"../user"
)

func TestBody(t *testing.T) {
	valid := &user.User{
		ID:   bson.NewObjectId(),
		Name: "Joshn",
		Role: "Tester",
	}
	valid2 := &user.User{
		ID:   "",
		Name: "Joshn",
		Role: "Developer",
	}
	js, err := json.Marshal(valid)
	if err != nil {
		t.Errorf("Error marshallaung %s", err)
	}
	ts := []struct {
		txt string
		r   *http.Request
		u   *user.User
		err bool
		exp *user.User
	}{
		{
			txt: "nil request",
			err: true,
		}, {
			txt: "empty request body",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString("{}")),
			},
			err: true,
		}, {
			txt: "malformed request body",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"id": 12}`)),
			},
			u:   &user.User{},
			err: true,
		}, {
			txt: "valid request body",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBuffer(js)),
			},
			u:   &user.User{},
			exp: valid,
		}, {
			txt: "valid partial request body",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"role": "Developer", "name": "Joshn"}`)),
			},
			u:   &user.User{},
			exp: valid2,
		},
	}
	for _, tc := range ts {
		t.Log(tc.txt)
		err := createUserBody(tc.r, tc.u)
		if tc.err {
			if err == nil {
				t.Error("Expected error, got none")
			}
			continue
		}
		if err != nil {
			t.Errorf("Erro marshalling a valid user : %s", err)
			t.FailNow()
			continue
		}
		if !reflect.DeepEqual(tc.u, tc.exp) {
			t.Error("Unmarshalled data is different")
			t.Error(tc.u)
			t.Error(tc.exp)
		}
	}
}
