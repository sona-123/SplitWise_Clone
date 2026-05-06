package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sona-123/splitwise_clone/utils"
)

func (h *Handler) GoogleLogin(c *gin.Context) {
	url := utils.GoogleOAuthConfig().AuthCodeURL("randomstate")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GoogleCallback(c *gin.Context) {

	//1.) Get authorization code from Google
	code := c.Query("code")

	//2.) Exchange code for token
	token, err := utils.GoogleOAuthConfig().Exchange(
		context.Background(),
		code,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to exchange token",
		})
		return
	}

	//3.) Get user info from google
	client := utils.GoogleOAuthConfig().Client(
		context.Background(),
		token,
	)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get user info",
		})
	}

	defer resp.Body.Close()

	var userData struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}

	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse User Info",
		})
		return
	}

	user, err := h.Service.HandleGoogleLogin(
		userData.Name,
		userData.Email,
		userData.Picture,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to process user",
		})
		return
	}

	jwtToken, err := utils.GenerateToken(user.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate JWT",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "google login Successful",
		"token":   jwtToken,
		"user":    user,
	})

}
