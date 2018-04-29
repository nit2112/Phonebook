package controllers

import (
	"encoding/json"
	"net/http"
	"Phonebook/common"
	"Phonebook/data"
	"Phonebook/models"
) 

func Register(w http.ResponseWriter, r *http.Request){
	var dataResource UserResource

	err := json.NewDecoder(r.Body).Decode(&dataResource)
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"Invalid User Data",
			500,
			)
		return
	}

	user := &dataResource.Data
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("users")
	repo := &data.UserRepository{c}
	repo.CreateUser(user)
	user.HashPassword = nil

	if j, err := json.Marshal(UserResource{Data: *user}); err != nil {
		common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
		return
	}  else {
		w.Header().Set("Context-Type","application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(j)
	}
}


func Login(w http.ResponseWriter, r *http.Request){
	var dataResource LoginResource
	var token string

	err := json.NewDecoder(r.Body).Decode(&dataResource)
	if err != nil {
		common.DisplayAppError(
			w,
			err,
			"Invalid Login Data",
			500,
		)
		return
	}
	loginModel := dataResource.Data
	loginUser := models.User{
		Email: loginModel.Email,
		Password: loginModel.Password,
	}
	context := NewContext()
	defer context.Close()
	c := context.DbCollection("users")
	repo := &data.UserRepository{c}
	if user, err := repo.Login(loginUser); err != nil {
		common.DisplayAppError(
			w,
			err,
			"Invalid login credentials",
			401,
		)
		return
	} else {
		token, err = common.GenerateJWT(user.Email,"member")
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"Error while generating the access token",
				500,
			)
			return
		}

		w.Header().Set("Content-Type","application/json")
		user.HashPassword = nil
		authUser := AuthUserModel{
			User: user,
			Token: token,
		}
		j, err := json.Marshal(AuthUserResource{Data: authUser})
		if err != nil {
			common.DisplayAppError(
				w,
				err,
				"An unexpected error has occurred",
				500,
			)
		  return	
		}
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}

}