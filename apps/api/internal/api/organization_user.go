package api

import (
	"api/internal/models"
	"api/pkg/firebase"
	"api/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (a *api) addUserToOrganization(w http.ResponseWriter, r *http.Request) {
	var payload models.AddUserOrganizationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := utils.GetTokenClaims(r)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "Token claims missing")
		return
	}

	userID, ok := utils.GetUserIDFromClaims(claims)
	if !ok {
		utils.JSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// check authorization for create user
	organizationID := r.Header.Get("organization_id")
	role, err := a.organizationUsersRepo.GetOrganizationUser(userID, organizationID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		utils.JSONError(w, http.StatusInternalServerError, "Unauthoruized")
		return
	}

	if role == "admin" || role == "owner" {
		user, err := a.userRepo.GetUserWithEmail(payload.Email)
		if err.Error() == "user not found" {
			err = firebase.InitFirebase()
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, err.Error())
				return
			}

			password, err := utils.GeneratePassword(12)
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, err.Error())
				return
			}

			err = firebase.CreateUserWithEmail(payload.Email, password)
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, err.Error())
				return
			}

			firebaseUser, err := firebase.GetUserByEmail(payload.Email)
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, err.Error())
				return
			}

			newUser := &models.User{
				Id:        firebaseUser.UID,
				FirstName: "Unknown",
				LastName:  "Unknown",
				Email:     payload.Email,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			_, err = a.userRepo.CreateUser(newUser)
			if err != nil {
				utils.JSONError(w, http.StatusInternalServerError, err.Error())
				return
			}

			user, err = a.userRepo.GetUserWithEmail(payload.Email)
			if err != nil {
				log.Printf("Error get user: %v", err)
				utils.JSONError(w, http.StatusInternalServerError, "No permission")
				return
			}

			utils.SendEmail(payload.Email, password, &a.cfg)
		} else if err != nil && err.Error() != "user not found" {
			log.Printf("Error get user: %v", err)
			utils.JSONError(w, http.StatusInternalServerError, "No permission")
			return
		}

		organization_user := &models.OrganizationUserCreated{
			OrganizationId: organizationID,
			UserId:         user.Id,
			Role:           payload.Role,
		}

		create_organization_user, err := a.organizationUsersRepo.CreateOrganizationUser(organization_user)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			utils.JSONError(w, http.StatusInternalServerError, "Failed to create organization user")
			return
		}

		utils.JSONResponse(w, http.StatusCreated, create_organization_user)
		return
	}

	utils.JSONError(w, http.StatusInternalServerError, "Failed to create organization user")
}
