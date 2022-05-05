package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"groceryMate/api/auth"
	"groceryMate/api/models"
	"groceryMate/api/responses"
	"groceryMate/api/utils/formaterror"
)

func (server *Server) CreateFollow(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the follow id is valid
	followeeID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	follow := models.Follow{}
	err = json.Unmarshal(body, &follow)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	follow.Prepare()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != follow.FollowerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	followCreated, err := follow.CreateFollow(server.DB, int32(followeeID))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, followCreated.ID))
	responses.JSON(w, http.StatusCreated, followCreated)
}

func (server *Server) GetFollows(w http.ResponseWriter, r *http.Request) {

	follow := models.Follow{}

	follows, err := follow.FindAllFollows(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, follows)
}

func (server *Server) GetFollow(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	follow := models.Follow{}

	followReceived, err := follow.FindFollowByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, followReceived)
}

func (server *Server) DeleteFollow(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid follow id given to us?
	fid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the follow exist
	follow := models.Follow{}
	err = server.DB.Debug().Model(models.Follow{}).Where("id = ?", fid).Take(&follow).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this follow?
	if uid != follow.FollowerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = follow.DeleteAFollow(server.DB, fid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", fid))
	responses.JSON(w, http.StatusNoContent, "")
}
