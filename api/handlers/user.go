package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/Mubinabd/library_auth/api/docs"
	t "github.com/Mubinabd/library_auth/api/token"
	auth "github.com/Mubinabd/library_auth/genproto/auth"
)

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieve the profile of a user with the specified ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} auth.UserRes
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /user/profiles [get]
func (h *Handlers) GetProfile(c *gin.Context) {
	userID := getuserId(c)
	req := &auth.GetById{
		Id: userID,
	}

	profile, err := h.User.GetProfile(c, req)
	if err != nil {
		log.Println("Error getting profile:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// EditProfile godoc
// @Summary Edit user profile
// @Description Update the profile of a user with the specified ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body auth.EditProfileReqBpdy true "Updated profile details"
// @Success 200 {object} string "Profile updated successfully"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /user/profiles [put]
func (h *Handlers) EditProfile(c *gin.Context) {
	userID := getuserId(c)

	var body auth.EditProfileReqBpdy
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	req := &auth.UserRes{
		Id:          userID,
		FirstName:    body.FirstName,
		Email:       body.Email,
		LastName:    body.LastName,
		PhoneNumber: body.PhoneNumber,
	}

	input, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
		return
	}

	err = h.Producer.ProduceMessages("upd-user", input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Profile for user %s updated successfully", req.Id)})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Update the password of a user with the specified ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body auth.ChangePasswordReqBody true "Updated password details"
// @Success 200 {object} string "Password updated successfully"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /user/passwords [put]
func (h *Handlers) ChangePassword(c *gin.Context) {
	userID := getuserId(c)

	var body auth.ChangePasswordReqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	password, err := t.HashPassword(body.NewPassword)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	body.NewPassword = password

	req := &auth.ChangePasswordReq{
		Id:              userID,
		CurrentPassword: body.CurrentPassword,
		NewPassword:     body.NewPassword,
	}

	input, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
		return
	}

	err = h.Producer.ProduceMessages("upd-pass", input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// GetSetting godoc
// @Summary Get user settings
// @Description Retrieve the settings of a user with the specified ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} auth.Setting
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /user/setting [get]
func (h *Handlers) GetSetting(c *gin.Context) {
	userID := getuserId(c)

	req := &auth.GetById{
		Id: userID,
	}

	setting, err := h.User.GetSetting(c, req)
	if err != nil {
		log.Println("Error getting setting:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// EditSetting godoc
// @Summary Edit user settings
// @Description Update the settings of a user with the specified ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param setting body auth.Setting true "Updated setting details"
// @Success 200 {object} string "Setting updated successfully"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /user/setting [put]
func (h *Handlers) EditSetting(c *gin.Context) {
	userID := getuserId(c)

	var body auth.Setting
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	req := &auth.SettingReq{
		Id:           userID,
		PrivacyLevel: body.PrivacyLevel,
		Notification: body.Notification,
		Language:     body.Language,
		Theme:        body.Theme,
	}

	input, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "internal server error"})
		return
	}

	err = h.Producer.ProduceMessages("upd-setting", input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Setting for user %s updated successfully", req.Id)})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user with the specified ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} string "User deleted successfully"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /user [delete]
func (h *Handlers) DeleteUser(c *gin.Context) {
	userID := getuserId(c)

	req := &auth.GetById{
		Id: userID,
	}

	_, err := h.User.DeleteUser(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s deleted successfully", req.Id)})
}
