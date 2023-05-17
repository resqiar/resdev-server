package handlers

import (
	"resdev-server/config"
	"resdev-server/services"

	"github.com/gofiber/fiber/v2"
)

func SendAuthGoogle(c *fiber.Ctx) error {
	// create a config for google config
	conf := config.GoogleConfig()

	// create url for auth process.
	// we can pass state as someway to identify
	// and validate the login process, for now skip it.
	URL := conf.AuthCodeURL("state")

	// redirect to the google authentication URL
	return c.Redirect(URL)
}

func SendGoogleCallback(c *fiber.Ctx) error {
	// get session store for current context
	sess, err := config.SessionStore.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	conf := config.GoogleConfig()
	code := c.Query("code")

	// exchange code that retrieved from google via
	// URL query parameter into token, this token
	// can be used later to query information of current user
	// from respective provider.
	token, err := conf.Exchange(c.Context(), code)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	profile, err := services.ConvertToken(token.AccessToken)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Store the user's data in the session
	sess.Set("Name", profile.Name)
	sess.Set("Email", profile.Email)

	// Save into memory session and.
	// saving also set a session cookie containing session_id
	if err := sess.Save(); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"Name":  profile.Name,
		"Email": profile.Email,
	})
}