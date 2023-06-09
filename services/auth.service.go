package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"resqiar.com-server/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/imagekit-developer/imagekit-go"
)

type AuthService interface {
	ConvertToken(accessToken string) (*entities.GooglePayload, error)
	SignIK(c *fiber.Ctx) imagekit.SignedToken
}

type AuthServiceImpl struct{}

func (service *AuthServiceImpl) ConvertToken(accessToken string) (*entities.GooglePayload, error) {
	resp, httpErr := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v3/userinfo?access_token=%s", accessToken))
	if httpErr != nil {
		return nil, httpErr
	}

	// clean up when this function returns (destroyed)
	defer resp.Body.Close()

	respBody, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		return nil, bodyErr
	}

	// Unmarshal raw response body to a map
	var body map[string]interface{}
	if err := json.Unmarshal(respBody, &body); err != nil {
		return nil, err
	}

	// if json body containing error,
	// then the token is indeed invalid. return invalid token err
	if body["error"] != nil {
		return nil, errors.New("Invalid token")
	}

	// Bind JSON into struct
	var data entities.GooglePayload
	err := json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (service *AuthServiceImpl) SignIK(c *fiber.Ctx) imagekit.SignedToken {
	IMAGE_KIT_KEY := os.Getenv("IMAGE_KIT_KEY")
	IMAGE_KIT_KEY_PUBLIC := os.Getenv("IMAGE_KIT_KEY_PUBLIC")
	IMAGE_KIT_URL := os.Getenv("IMAGE_KIT_URL")

	// Initialize image kit with provided params
	ik := imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  IMAGE_KIT_KEY,
		PublicKey:   IMAGE_KIT_KEY_PUBLIC,
		UrlEndpoint: IMAGE_KIT_URL,
	})

	// return an Object containing Token, Signature and Expire
	signed := ik.SignToken(imagekit.SignTokenParam{})
	return signed
}
