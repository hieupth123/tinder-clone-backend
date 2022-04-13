package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phamtrunghieu/tinder-clone-backend/constant"
	"github.com/phamtrunghieu/tinder-clone-backend/helpers/common"
	"github.com/phamtrunghieu/tinder-clone-backend/helpers/respond"
	"github.com/phamtrunghieu/tinder-clone-backend/helpers/util"
	"github.com/phamtrunghieu/tinder-clone-backend/models"
	request "github.com/phamtrunghieu/tinder-clone-backend/request/user_action"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"sync"
)


type UserActionController struct {
}

func (uCtrl UserActionController) LikeUser(c *gin.Context) {
	var req request.LikeUserRequest
	errUri := c.ShouldBindUri(&req)
	if errUri != nil {
		_ = c.Error(errUri)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	userActionModel := new(models.UserAction)
	cond := bson.M{
		"user_uuid": req.GuestUuid,
		"guest_uuid": req.UserUuid,
	}
	isGuestLike, err := userActionModel.FindOne(c, cond)
	if err != nil {
		fmt.Println(err)
	}
	obj := models.UserAction{
		Uuid: common.GenerateUUID(),
		UserUuid: req.UserUuid,
		GuestUuid: req.GuestUuid,
		Type: constant.LIKED,
		CreatedAt: util.GetNowUTC(),
		UpdatedAt: util.GetNowUTC(),
	}
	obj.Insert()
	if isGuestLike != nil {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			PushMatchesUser(req.UserUuid, req.GuestUuid)
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			PushMatchesUser(req.GuestUuid, req.UserUuid)
		}()

		wg.Wait()
	}
	c.JSON(http.StatusOK, respond.Success(nil, "Liked successfully"))
}

func (uCtrl UserActionController) PassUser(c *gin.Context) {
	var req request.LikeUserRequest
	errUri := c.ShouldBindUri(&req)
	if errUri != nil {
		_ = c.Error(errUri)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	userActionModel := new(models.UserAction)
	cond := bson.M{
		"user_uuid": req.GuestUuid,
		"guest_uuid": req.UserUuid,
	}
	_, err := userActionModel.FindOne(c, cond)
	if err != nil {
		fmt.Println(err)
		obj := models.UserAction{
			Uuid: common.GenerateUUID(),
			UserUuid: req.UserUuid,
			GuestUuid: req.GuestUuid,
			Type: constant.PASSED,
			CreatedAt: util.GetNowUTC(),
			UpdatedAt: util.GetNowUTC(),
		}
		obj.Insert()
	}
	c.JSON(http.StatusOK, respond.Success(nil, "Liked successfully"))
}
func PushMatchesUser (userUuid string, guestUuid string) {
	userModel := new(models.User)
	condUpdateUser := bson.M{
		"uuid":        userUuid,
	}
	dataU := make(map[string]interface{})
	dataU["matches"] = guestUuid

	_, errUpdateUser := userModel.UpdatePushDataToArray(condUpdateUser, dataU)
	if errUpdateUser != nil {
		util.LogError(errUpdateUser)
	}
}

func (uCtrl UserActionController) GetUserLiked(c *gin.Context) {
	var req request.LikeUserRequest
	errUri := c.ShouldBindUri(&req)
	if errUri != nil {
		_ = c.Error(errUri)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	userActionModel := new(models.UserAction)
	userModel := new (models.User)
	cond := bson.M{
		"user_uuid": req.GuestUuid,
		"type": constant.LIKED,
	}
	userActions, err := userActionModel.Find(c, cond)
	if err != nil {
		util.LogError(err)
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorResponse("Get User Action fail"))
		return
	}
	var userUuids []string
	for _, userAction := range userActions{
		userUuids = append(userUuids, userAction.GuestUuid)
	}
	response := []*request.UserLikedResponse{}
	if len(userUuids) == 0 {
		c.JSON(http.StatusOK, respond.Success(response, "Get Liked User Successfully"))
		return
	}
	condUser := bson.M{
		"uuid": bson.M{"$in": userUuids},
	}
	users, errUser := userModel.Find(c, condUser)
	if errUser != nil {
		util.LogError(errUser)
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorResponse("Get user liked fail"))
		return
	}
	for _, item := range users {
		response = append(response, &request.UserLikedResponse{
			Uuid:      item.Uuid,
			LastName:  item.LastName,
			FirstName: item.FirstName,
			Gender:    util.ConvertGender(item.Gender),
			Age:       util.GetAge(item.DateOfBirth),
			Picture:   item.Picture,
		})
	}
	c.JSON(http.StatusOK, respond.Success(nil, "Liked successfully"))
}