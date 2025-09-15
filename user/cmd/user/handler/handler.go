package handler

import (
	"net/http"
	"user/cmd/user/usecase"
	"user/infrastructure/log"
	"user/models"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		UserUsecase: userUsecase,
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	var param models.LoginParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Logger.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid input parameter",
		})
		return
	}

	if len(param.Password) < 8 {
		log.Logger.Info("Invalid Input")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Password must longer than 8 characters",
		})
		return
	}

	token, err := h.UserUsecase.Login(c.Request.Context(), &param)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email or password mismatch",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *UserHandler) Register(c *gin.Context) {
	var param models.RegisterParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Logger.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid parameter input",
		})
		return
	}

	if len(param.Password) < 8 || len(param.ConfirmPassword) < 8 {
		log.Logger.Info("Invalid Input")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Password must longer than 8 characters",
		})
		return
	}

	if param.Password != param.ConfirmPassword {
		log.Logger.Info("Invalid Input")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Password mismatched",
		})
		return
	}

	user, err := h.UserUsecase.GetUserByEmail(c.Request.Context(), param.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err.Error(),
		})
		return
	}

	if user != nil && user.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Email already exists!",
		})
		return
	}

	err = h.UserUsecase.RegisterUser(c.Request.Context(), &models.User{
		Name:     param.Name,
		Email:    param.Email,
		Password: param.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User successfully registered!",
	})
}

func (h *UserHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

// func (h *UserHandler) GetUserInfo(c *gin.Context) {
// 	userIDStr, isExist := c.Get("user_id")
// 	if !isExist {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"error_message": "Unauthorized",
// 		})
// 		return
// 	}

// 	userID, ok := userIDStr.(float64)
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"error_message": "Invalid user id",
// 		})
// 		return
// 	}

// 	user, err := h.UserUsecase.GetUserByUserID(c.Request.Context(), int64(userID))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error_message": err.Error(),
// 		})
// 		return
// 	}

// 	if user.ID == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error_message": "User not found!",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"name":  user.Name,
// 		"email": user.Email,
// 	})
// }

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userIDVal, isExist := c.Get("user_id")
	if !isExist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error_message": "Unauthorized",
		})
		return
	}

	switch id := userIDVal.(type) {
	case string:
		// If user_id is email
		user, err := h.UserUsecase.GetUserByEmail(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
			return
		}
		if user == nil || user.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error_message": "User not found!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"name": user.Name, "email": user.Email})

		// case int64:
		// 	// If user_id is numeric ID
		// 	user, err := h.UserUsecase.GetUserByUserID(c.Request.Context(), id)
		// 	if err != nil {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})
		// 		return
		// 	}
		// 	if user == nil || user.ID == 0 {
		// 		c.JSON(http.StatusNotFound, gin.H{"error_message": "User not found!"})
		// 		return
		// 	}
		// 	c.JSON(http.StatusOK, gin.H{"name": user.Name, "email": user.Email})

		// default:
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error_message": "Invalid user_id type"})
	}
}
