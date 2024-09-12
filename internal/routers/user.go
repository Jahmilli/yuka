package routers

import (
	"net/http"

	"yuka/internal/handlers"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// createUser creates a User
// @Summary      Create User
// @Id  		 createUser
// @Tags         Users
// @Description  Creates a user
// @Accept	     json
// @Produce      json
// @Param		 create body handlers.CreateUserInput true "User Create"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.ValidationError
// @Failure      500  {object}  models.BaseError
// @Router       /api/v1/users [post]
func createUser(handler handlers.UserHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input handlers.CreateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := handler.CreateUser(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// getUser gets a User for the specified id
// @Summary      Get User for specified id
// @Id  		 getUser
// @Tags         Users
// @Description  Gets users
// @Param        id    path      string          true  "User ID"
// @Accept	     json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      500  {object}  models.BaseError
// @Router       /api/v1/users/{id} [get]
func getUser(handler handlers.UserHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := handler.FindUser(handlers.FindUserKeyAuthID, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// updateUser updates a User
// @Summary      Update User
// @Id  		 updateUser
// @Tags         Users
// @Description  Updates a user
// @Param        id    path      string          true  "User ID"
// @Accept	     json
// @Produce      json
// @Param		 create body handlers.UpdateUserInput true "User Update"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.ValidationError
// @Failure      500  {object}  models.BaseError
// @Router       /api/v1/users/{id} [put]
func updateUser(handler handlers.UserHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input handlers.UpdateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		user, err := handler.UpdateUser(c.Param("id"), input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// deleteUser deletes a User
// @Summary      Delete User
// @Id  		 deleteUser
// @Tags         Users
// @Description  Deletes a user
// @Param        id    path      string          true  "User ID"
// @Accept	     json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      500  {object}  models.BaseError
// @Router       /api/v1/users/{id} [delete]
func deleteUser(handler handlers.UserHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler.DeleteUser(c.Param("id")); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": "ok"})
	}
}
