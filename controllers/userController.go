package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phamtrunghieu/tinder-clone-backend/config"
	"github.com/phamtrunghieu/tinder-clone-backend/helpers/common"
	"github.com/phamtrunghieu/tinder-clone-backend/helpers/respond"
	"github.com/phamtrunghieu/tinder-clone-backend/helpers/util"
	"github.com/phamtrunghieu/tinder-clone-backend/models"
	request "github.com/phamtrunghieu/tinder-clone-backend/request/user"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type ListUser struct {
	Data []UserId `json:"data"`
}
type UserDetail struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Picture     string `json:"picture"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"dateOfBirth"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}
type UserId struct {
	Id       string `json:"id"`
	LastName string `json:"lastName"`
}

type UserController struct {
}

func (uCtrl UserController) DumpData(c *gin.Context) {
	cfg := config.GetConfig()
	apiKey := cfg.GetString("dummy.appId")
	baseUrl := cfg.GetString("dummy.Url")
	client := http.Client{
		Timeout: 20 * time.Second,
	}

	url := baseUrl + "/user?limit=10"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[ERROR] Get store recommendation http.NewRequest %s\n", err)
		fmt.Println("url: ", url)
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorCommonNotFound("Create request failed"))
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("app-id", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] Get store recommendation client.Do %s\n", err)
		fmt.Println("url: ", url)
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorCommonNotFound("Send request failed"))
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[ERROR] Get store recommendation ioutil.ReadAll %s\n", err)
		fmt.Println("url: ", url)
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorCommonNotFound("Read response failed"))
		return
	}
	var responseData *ListUser
	if err := json.Unmarshal(b, &responseData); err != nil {
		fmt.Printf("[ERROR] json.Unmarshal %s %s\n", err, string(b))
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorCommonNotFound("Unmarshall json failed"))
		return
	}

	util.LogInfo("responseData.Data: ", responseData.Data)

	if len(responseData.Data) > 0 {
		var wg sync.WaitGroup
		for _, item := range responseData.Data {
			wg.Add(1)
			go getDummyAndSaveData(item.Id, apiKey, baseUrl, client, &wg)
		}
		wg.Wait()
	}
	c.JSON(http.StatusOK, respond.Success(nil, "Successfully"))
}

func getDummyAndSaveData(id string, apiKey string, baseUrl string, client http.Client, wg *sync.WaitGroup) {
	url := baseUrl + "/user/" + id
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[ERROR] Get store recommendation http.NewRequest %s\n", err)
		fmt.Println("url: ", url)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("app-id", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] Get store recommendation client.Do %s\n", err)
		fmt.Println("url: ", url)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[ERROR] Get store recommendation ioutil.ReadAll %s\n", err)
		fmt.Println("url: ", url)
		return
	}
	var responseData *UserDetail
	if err := json.Unmarshal(b, &responseData); err != nil {
		fmt.Printf("[ERROR] json.Unmarshal %s %s\n", err, string(b))
		return
	}
	matches := []string{}
	user := models.User{
		Uuid:        common.GenerateUUID(),
		LastName:    responseData.LastName,
		FirstName:   responseData.FirstName,
		Gender:      responseData.Gender,
		Picture:     responseData.Picture,
		DateOfBirth: responseData.DateOfBirth,
		Email:       responseData.Email,
		Phone:       responseData.Phone,
		Matches:     matches,
		CreatedAt:   util.GetNowUTC(),
		UpdatedAt:   util.GetNowUTC(),
	}
	user.Insert()
	defer wg.Done()
}

func (uCtrl UserController) GetList(c *gin.Context) {
	userModel := new(models.User)

	cond := bson.M{}

	users, err := userModel.Find(c, cond)

	if err != nil {
		util.LogError(err)
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorResponse("Get user fail"))
		return
	}

	response := []*request.GetListResponse{}
	for _, item := range users {
		response = append(response, &request.GetListResponse{
			Uuid:      item.Uuid,
			LastName:  item.LastName,
			FirstName: item.FirstName,
			Gender:    util.ConvertGender(item.Gender),
			Age:       util.GetAge(item.DateOfBirth),
			Picture:   item.Picture,
		})
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}

func (uCtrl UserController) GetDetail(c *gin.Context) {
	userModel := new(models.User)
	var req request.DetailRequest
	errUri := c.ShouldBindUri(&req)
	if errUri != nil {
		_ = c.Error(errUri)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	cond := bson.M{
		"uuid": req.Uuid,
	}
	user, err := userModel.FindOne(c, cond)

	if err != nil {
		util.LogError(err)
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorResponse("Can not user by uuid"))
		return
	}

	response := request.GetDetailResponse{
		Uuid:      user.Uuid,
		LastName:  user.LastName,
		FirstName: user.FirstName,
		Gender:    util.ConvertGender(user.Gender),
		Age:       util.GetAge(user.DateOfBirth),
		Picture:   user.Picture,
		Email:     user.Email,
		Phone:     user.Phone,
		Matches:   user.Matches,
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}

func (uCtrl UserController) GetUserRandom(c *gin.Context) {
	userModel := new(models.User)
	var req request.DetailRequest
	errUri := c.ShouldBindUri(&req)
	if errUri != nil {
		_ = c.Error(errUri)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	user, err := userModel.FindRandomUser(c)

	if err != nil {
		util.LogError(err)
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorResponse("Can not get random user"))
		return
	}
	response := request.GetDetailResponse{
		Uuid:      user.Uuid,
		LastName:  user.LastName,
		FirstName: user.FirstName,
		Gender:    util.ConvertGender(user.Gender),
		Age:       util.GetAge(user.DateOfBirth),
		Picture:   user.Picture,
		Email:     user.Email,
		Phone:     user.Phone,
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}

func (uCtrl UserController) GetUserAvailable(c *gin.Context) {
	userModel := new(models.User)
	userActionModel := new(models.UserAction)
	var req request.DetailRequest
	errUri := c.ShouldBindUri(&req)
	if errUri != nil {
		_ = c.Error(errUri)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	cond := bson.M{
		"uuid": req.Uuid,
	}
	user, err := userModel.FindOne(c, cond)
	if err != nil {
		util.LogError(err)
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorResponse("Can not found user by uuid"))
		return
	}
	exceptUuids := []string{}
	exceptUuids = user.Matches
	condAction := bson.M{
		"user_uuid": req.Uuid,
	}

	userActions, err := userActionModel.Find(c, condAction)
	if err != nil {
		util.LogError(err)
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorResponse("Can not found user action"))
		return
	}
	for _, userAction := range userActions {
		exceptUuids = append(exceptUuids, userAction.GuestUuid)
	}
	exceptUuids = append(exceptUuids, req.Uuid)
	fmt.Println("exceptUuids: ", exceptUuids)
	condAvl := bson.M{
		"uuid": bson.M{"$nin": exceptUuids},
	}
	fmt.Println("condAvl: ")
	util.LogError(condAvl)
	users, err := userModel.Find(c, condAvl)
	response := []*request.GetDetailResponse{}
	if err != nil {
		util.LogError(err)
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorResponse("Can not get list user available"))
		return
	}

	for _, item := range users {
		response = append(response, &request.GetDetailResponse{
			Uuid:      item.Uuid,
			LastName:  item.LastName,
			FirstName: item.FirstName,
			Gender:    util.ConvertGender(item.Gender),
			Age:       util.GetAge(item.DateOfBirth),
			Picture:   item.Picture,
			Matches:   item.Matches,
		})
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}

func (uCtrl UserController) GetMatchesUser(c *gin.Context) {
	userModel := new(models.User)
	var req request.DetailRequest
	errUri := c.ShouldBindUri(&req)
	if errUri != nil {
		_ = c.Error(errUri)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	cond := bson.M{
		"uuid": req.Uuid,
	}
	user, err := userModel.FindOne(c, cond)

	if err != nil {
		util.LogError(err)
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorResponse("Can not user by uuid"))
		return
	}
	response := []*request.GetDetailResponse{}
	if len(user.Matches) == 0 {
		c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
		return
	}

	condUsers := bson.M{
		"uuid": bson.M{"$in": user.Matches},
	}

	users, errUser := userModel.Find(c, condUsers)
	if errUser != nil {
		util.LogError(errUser)
		c.JSON(http.StatusUnprocessableEntity, respond.ErrorResponse("Get user liked fail"))
		return
	}
	for _, item := range users {
		response = append(response, &request.GetDetailResponse{
			Uuid:      item.Uuid,
			LastName:  item.LastName,
			FirstName: item.FirstName,
			Gender:    util.ConvertGender(item.Gender),
			Age:       util.GetAge(item.DateOfBirth),
			Picture:   item.Picture,
		})
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}
