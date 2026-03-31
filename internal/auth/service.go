package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/pkg"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleProvider = "GOOGLE"

type AuthService struct {
	OauthConfig *oauth2.Config
	repo        *AuthRepository
}

func NewAuthService(googleClientId string, googleSecret string, authRepo *AuthRepository) *AuthService {

	return &AuthService{
		OauthConfig: &oauth2.Config{
			ClientID:     googleClientId,
			ClientSecret: googleSecret,
			RedirectURL:  "http://localhost/api/auth/google-callback",
			Scopes:       []string{"email", "profile", "openid"},
			Endpoint:     google.Endpoint,
		},

		repo: authRepo,
	}

}

func (as *AuthService) GetLoginGoogleUrl() (string, *shared.ErrorResponse) {

	state, err := pkg.GenerateSecureString(32)
	if err != nil {
		fmt.Println("couldn't generate the token for google oauth!")
		return "", shared.NewErrorResponse(500, "something went wrong with the server")
	}

	redUrl := as.OauthConfig.AuthCodeURL(state, oauth2.ApprovalForce)

	return redUrl, nil

}

func (as *AuthService) AuthenticateGoogleUser(ctx context.Context, code string) (a string, r string, svcErr *shared.ErrorResponse) {

	t, err := as.OauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", "", shared.NewErrorResponse(500, "something went wrong with the server")
	}

	client := as.OauthConfig.Client(ctx, t)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", "", shared.NewErrorResponse(500, "something went wrong with the server")
	}

	defer resp.Body.Close()

	var googleUser googleUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return "", "", shared.NewErrorResponse(500, "failed to decode user data")
	}

	existingData, isExist, err := as.repo.ExistByProviderId(ctx, googleUser.ID, googleProvider)
	if err != nil && !errors.Is(err, errAccountNotFound) {
		return "", "", shared.NewErrorResponse(500, "something wrong with the server "+err.Error())
	}

	if !isExist {
		newData, err := as.repo.CreateNewUser(ctx, createUser{
			FullName:       googleUser.Name,
			Email:          googleUser.Email,
			ProfilePicture: &googleUser.Picture,
			Provider:       googleProvider,
			ProviderUserID: googleUser.ID,
		})

		if err != nil {
			return "", "", shared.NewErrorResponse(500, "something wrong with the server"+err.Error())
		}

		accessToken, err := pkg.GenerateToken(newData.ID, newData.GlobalRole, 12)
		if err != nil {
			log.Fatal("COULDN'T GENERATE THE TOKEN! PANIC! " + err.Error())
			return "", "", shared.NewErrorResponse(500, "something wrong with the server"+err.Error())
		}

		refreshToken, err := pkg.GenerateTokenNoRole(newData.ID, 12)
		if err != nil {
			log.Fatal("COULDN'T GENERATE THE TOKEN! PANIC! " + err.Error())
			return "", "", shared.NewErrorResponse(500, "something wrong with the server"+err.Error())
		}

		return accessToken, refreshToken, nil
	}

	accessToken, err := pkg.GenerateToken(existingData.UserId, existingData.Role, 12)
	if err != nil {
		log.Fatal("COULDN'T GENERATE THE TOKEN! PANIC! " + err.Error())
		return "", "", shared.NewErrorResponse(500, "something wrong with the server"+err.Error())
	}

	refreshToken, err := pkg.GenerateTokenNoRole(existingData.UserId, 12)
	if err != nil {
		log.Fatal("COULDN'T GENERATE THE TOKEN! PANIC! " + err.Error())
		return "", "", shared.NewErrorResponse(500, "something wrong with the server"+err.Error())
	}

	return accessToken, refreshToken, nil
}

func (as *AuthService) GetRefreshSession(ctx context.Context, userId uuid.UUID) (string, *shared.ErrorResponse) {

	userData, err := as.repo.GetUserById(ctx, userId)
	if err != nil {

		if errors.Is(err, errAccountNotFound) {
			return "", shared.NewErrorResponse(404, "no account with this id found!")
		}

		return "", shared.NewErrorResponse(500, "something wrong with the server"+err.Error())
	}

	accessToken, err := pkg.GenerateToken(userData.ID, userData.GlobalRole, 12)

	return accessToken, nil
}
