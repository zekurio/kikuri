package discordoauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/zekurio/daemon/internal/services/webserver/auth"
)

const (
	endpointOauth = "https://discord.com/api/oauth2/token"
	endpointMe    = "https://discord.com/api/users/@me"
	endpointAuth  = "https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify"
)

// OnErrorFunc gets called when an error occurs during the OAuth process.
type OnErrorFunc func(ctx *fiber.Ctx, status int, msg string) error

// OnSuccessFunc gets called when the OAuth process is successful.
type OnSuccessFuc func(ctx *fiber.Ctx, claims auth.Claims) error

// DiscordOAuth is the Discord OAuth handler.
type DiscordOAuth struct {
	clientID     string
	clientSecret string
	redirectURI  string

	onError   OnErrorFunc
	onSuccess OnSuccessFuc
}

type oAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
}

type userMeResponse struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}

// New creates our Discord OAuth handler.
func New(clientID string, clientSecret string, redirectURI string, onError OnErrorFunc, onSuccess OnSuccessFuc) *DiscordOAuth {
	if onError == nil {
		onError = func(ctx *fiber.Ctx, status int, msg string) error {
			return ctx.Status(status).SendString(msg)
		}
	}

	if onSuccess == nil {
		onSuccess = func(ctx *fiber.Ctx, claims auth.Claims) error {
			return ctx.Status(fiber.StatusOK).JSON(claims)
		}
	}

	return &DiscordOAuth{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		onError:      onError,
		onSuccess:    onSuccess,
	}
}

// HandlerInit returns the redirect respone to the Discord OAuth page.
func (d *DiscordOAuth) HandlerInit(ctx *fiber.Ctx) error {
	uri := fmt.Sprintf(endpointAuth, d.clientID, url.QueryEscape(d.redirectURI))

	return ctx.Redirect(uri, fiber.StatusTemporaryRedirect)
}

// HandlerCallback handles the callback from the Discord OAuth page.
// It will verify the token by getting a bearer token from Discord and
// then getting the the oauth2/user/@me endpoint.
func (d *DiscordOAuth) HandlerCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	if code == "" {
		return d.onError(ctx, fiber.StatusBadRequest, "missing code")
	}

	data := map[string][]string{
		"client_id":     {d.clientID},
		"client_secret": {d.clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {d.redirectURI},
		"scope":         {"identify"},
	}

	values := url.Values(data)
	req, err := http.NewRequest("POST", endpointOauth, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		d.onError(ctx, http.StatusInternalServerError, "failed executing request: "+err.Error())
		return nil
	}

	if res.StatusCode >= 300 {
		d.onError(ctx, http.StatusUnauthorized, "")
		return nil
	}

	resAuthBody := new(oAuthTokenResponse)
	err = parseJSONBody(res.Body, resAuthBody)
	if err != nil {
		d.onError(ctx, http.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
		return nil
	}

	if resAuthBody.Error != "" || resAuthBody.AccessToken == "" {
		d.onError(ctx, http.StatusUnauthorized, "")
		return nil
	}

	// 2. Request getting user ID

	req, err = http.NewRequest("GET", endpointMe, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resAuthBody.AccessToken))

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		d.onError(ctx, http.StatusInternalServerError, "failed executing request: "+err.Error())
		return nil
	}

	if res.StatusCode >= 300 {
		d.onError(ctx, http.StatusUnauthorized, "")
		return nil
	}

	resGetMe := new(userMeResponse)
	err = parseJSONBody(res.Body, resGetMe)
	if err != nil {
		d.onError(ctx, http.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
		return nil
	}

	if resGetMe.Error != "" || resGetMe.ID == "" {
		d.onError(ctx, http.StatusUnauthorized, "")
		return nil
	}

	d.onSuccess(ctx, auth.Claims{
		UserID: resGetMe.ID,
		Scopes: []string{string(auth.AuthOriginDiscord)},
	})

	return nil
}

func parseJSONBody(body io.Reader, v interface{}) error {
	return json.NewDecoder(body).Decode(v)
}
