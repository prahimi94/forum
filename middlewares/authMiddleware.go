package middlewares

import (
	"context"
	errorManagementControllers "forum/modules/errorManagement/controllers"
	userManagementControllers "forum/modules/userManagement/controllers"
	"net/http"
)

// Context key for storing user data
type ContextKey string

const UserContextKey ContextKey = "user"

// AuthMiddleware ensures that only authenticated users can access certain routes
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r)
		if checkLoginError != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}
		if loginStatus {
			// Store user in context
			ctx := context.WithValue(r.Context(), UserContextKey, loginUser)

			// Pass request with updated context to the next handler
			next(w, r.WithContext(ctx))
			// return
		} else {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
			return
		}
	}
}
