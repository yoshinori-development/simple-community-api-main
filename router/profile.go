package router

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yoshinori-development/simple-community-api-main/models"
	"github.com/yoshinori-development/simple-community-api-main/services"
	"gorm.io/gorm"
)

type ProfileHandler struct {
	ProfileService services.ProfileService
}

type NewProfileHandlerInput struct {
	ProfileService services.ProfileService
}

func NewProfileHandler(input NewProfileHandlerInput) *ProfileHandler {
	return &ProfileHandler{
		ProfileService: input.ProfileService,
	}
}

type ProfileResponse struct {
	Nickname  string    `json:"nickname"`
	Age       uint      `json:"age"`
	UpdatedAt time.Time `json:"updateAt"`
}

func (controller *ProfileHandler) Get(c *gin.Context) {
	sub, _ := c.Get("sub")
	profile, err := controller.ProfileService.Get(services.ProfileServiceGetInput{
		Sub: sub.(string),
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, RenderMessageError(err, "プロフィールが登録されていません"))
		} else {
			c.JSON(http.StatusInternalServerError, RenderMessageError(err, "エラーが発生しました"))
		}
		return
	}

	response := ProfileResponse{
		Nickname:  profile.Nickname,
		Age:       profile.Age,
		UpdatedAt: profile.UpdatedAt,
	}
	c.JSON(http.StatusOK, response)
}

type ProfileHandlerCreateOrUpdateInput struct {
	Nickname string `form:"nickname" binding:"required,min=5"`
	Age      uint   `form:"age" binding:"numeric,max=150"`
}

func (controller *ProfileHandler) CreateOrUpdate(c *gin.Context) {
	var input ProfileHandlerCreateOrUpdateInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, RenderValidationError(err))
		return
	}

	sub, _ := c.Get("sub")
	err := controller.ProfileService.CreateOrUpdate(services.ProfileServiceCreateOrUpdateInput{
		Profile: models.Profile{
			Sub:      sub.(string),
			Nickname: input.Nickname,
			Age:      input.Age,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, RenderMessageError(err, "エラーが発生しました"))
	}

	c.Status(http.StatusOK)
}
