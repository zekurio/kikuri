package discordoauth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
)

const (
	endpointOauth = "https://discord.com/api/oauth2/token"
	endpointMe    = "https://discord.com/api/users/@me"
)

type SuccessResponse struct {
	UserID string
	State  map[string]interface{}
}

type tokenResponse struct {
	Error       string `json:"error"`
	AccessToken string `json:"access_token"`
}

type userMeResponse struct {
	Error string `json:"error"`
	ID    string `json:"id"`
}

// OnErrorFunc is the function to be used to handle errors during
// authentication.
type OnErrorFunc func(ctx *fiber.Ctx, status int, msg string) error

// OnSuccessFuc is the func to be used to handle the successful
// authentication.
type OnSuccessFuc func(ctx *fiber.Ctx, res SuccessResponse) error

type DiscordOAuth struct {
	clientID        string
	clientSecret    string
	redirectURI     string
	stateSigningKey []byte

	onError   OnErrorFunc
	onSuccess OnSuccessFuc
}

func New(clientID, clientSecret, redirectURI string, onError OnErrorFunc, onSuccess OnSuccessFuc) (doauth *DiscordOAuth, err error) {
	if onError == nil {
		onError = func(ctx *fiber.Ctx, status int, msg string) error { return nil }
	}
	if onSuccess == nil {
		onSuccess = func(ctx *fiber.Ctx, res SuccessResponse) error { return nil }
	}

	signingKey := make([]byte, 128)
	_, err = rand.Read(signingKey)
	if err != nil {
		return nil, err
	}

	doauth = &DiscordOAuth{
		clientID:        clientID,
		clientSecret:    clientSecret,
		redirectURI:     redirectURI,
		stateSigningKey: signingKey,

		onError:   onError,
		onSuccess: onSuccess,
	}

	return
}

// HandlerInit returns a redirect response to discord oauth2 authorization endpoint.
func (d *DiscordOAuth) HandlerInit(ctx *fiber.Ctx) error {
	return d.HandlerInitWithState(ctx, nil)
}

// HandlerInitWithState returns a redirect response to discord oauth2 authorization endpoint
// with a state token that contains the state map encoded and signed.
func (d *DiscordOAuth) HandlerInitWithState(ctx *fiber.Ctx, state map[string]interface{}) error {
	stateToken, err := d.encodeAndSignWithPayload(state)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify&state=%s",
		d.clientID, url.QueryEscape(d.redirectURI), stateToken)
	return ctx.Redirect(uri, fiber.StatusTemporaryRedirect)
}

// HandlerCallback handles the callback from discord oauth2 authorization endpoint.
func (d *DiscordOAuth) HandlerCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if state == "" {
		return d.onError(ctx, fiber.StatusBadRequest, "missing state")
	}

	claims, err := d.decodeAndValidateState(state)
	if err != nil {
		switch err {
		case jwt.ErrSignatureInvalid:
			return d.onError(ctx, fiber.StatusBadRequest, "invalid state")
		case jwt.ErrTokenExpired, jwt.ErrTokenNotValidYet, jwt.ErrTokenInvalidClaims:
			return d.onError(ctx, fiber.StatusBadRequest, "invalid state")
		default:
			return err
		}
	}

	// request bearer token by exchanging code
	data := map[string][]string{
		"client_id":     {d.clientID},
		"client_secret": {d.clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {d.redirectURI},
		"scope":         {"identify"},
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	values := url.Values(data)

	req.Header.SetMethod("POST")
	req.SetRequestURI(endpointOauth)
	req.SetBody([]byte(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err := fasthttp.Do(req, res); err != nil {
		return d.onError(ctx, fasthttp.StatusInternalServerError, "failed executing request: "+err.Error())
	}

	if res.StatusCode() >= 300 {
		return d.onError(ctx, fasthttp.StatusUnauthorized, "invalid auth code")
	}

	resAuthBody := new(tokenResponse)
	err = parseJSONBody(res.Body(), resAuthBody)
	if err != nil {
		return d.onError(ctx, fasthttp.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
	}

	if resAuthBody.Error != "" || resAuthBody.AccessToken == "" {
		return d.onError(ctx, fasthttp.StatusUnauthorized, "empty auth response")
	}

	// request user id

	req.Header.Reset()
	req.ResetBody()
	req.Header.SetMethod("GET")
	req.SetRequestURI(endpointMe)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resAuthBody.AccessToken))

	if err = fasthttp.Do(req, res); err != nil {
		return d.onError(ctx, fasthttp.StatusInternalServerError, "failed executing request: "+err.Error())
	}

	if res.StatusCode() >= 300 {
		return d.onError(ctx, fasthttp.StatusUnauthorized, "user request failed")
	}

	userMeRes := new(userMeResponse)
	err = parseJSONBody(res.Body(), userMeRes)
	if err != nil {
		return d.onError(ctx, fasthttp.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
	}

	if userMeRes.Error != "" || userMeRes.ID == "" {
		return d.onError(ctx, fasthttp.StatusUnauthorized, "empty user response")
	}

	return d.onSuccess(ctx, SuccessResponse{
		UserID: userMeRes.ID,
		State:  claims["payload"].(map[string]interface{}),
	})
}

func (d *DiscordOAuth) encodeAndSignWithPayload(payload map[string]interface{}) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	}).SignedString(d.stateSigningKey)
}

func (d *DiscordOAuth) decodeAndValidateState(token string) (claims jwt.MapClaims, err error) {
	_, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return d.stateSigningKey, nil
	})

	return claims, err
}

func parseJSONBody(body []byte, v interface{}) error {
	return json.Unmarshal(body, v)
}
