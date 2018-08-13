package permissions

import (
	"context"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"gopkg.in/gin-gonic/gin.v1"
)

type Roles struct {
	Roles []string `json:"roles"`
}

type IDTokenClaims struct {
	Acr               string           `json:"acr"`
	AllowedOrigins    []string         `json:"allowed-origins"`
	Aud               string           `json:"aud"`
	AuthTime          int              `json:"auth_time"`
	Azp               string           `json:"azp"`
	Email             string           `json:"email"`
	Exp               int              `json:"exp"`
	FamilyName        string           `json:"family_name"`
	GivenName         string           `json:"given_name"`
	Iat               int              `json:"iat"`
	Iss               string           `json:"iss"`
	Jti               string           `json:"jti"`
	Name              string           `json:"name"`
	Nbf               int              `json:"nbf"`
	Nonce             string           `json:"nonce"`
	PreferredUsername string           `json:"preferred_username"`
	RealmAccess       Roles            `json:"realm_access"`
	ResourceAccess    map[string]Roles `json:"resource_access"`
	SessionState      string           `json:"session_state"`
	Sub               string           `json:"sub"`
	Typ               string           `json:"typ"`
}

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenVerifier, _ := c.Get("TOKEN_VERIFIER")
		if verifier, ok := tokenVerifier.(*oidc.IDTokenVerifier); ok && verifier != nil {
			// We have a proper ID Token Verifier. Game on

			authHeader := strings.Split(strings.TrimSpace(c.Request.Header.Get("Authorization")), " ")
			if len(authHeader) == 2 || strings.ToLower(authHeader[0]) == "bearer" {
				// Authorization header provided, let's verify.

				token, err := verifier.Verify(context.TODO(), authHeader[1])
				if err != nil {
					c.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
					return
				}

				// ID Token is verified. WooHoo !
				c.Set("ID_TOKEN", token)

				// parse claims
				var claims IDTokenClaims
				if err := token.Claims(&claims); err != nil {
					c.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePublic)
					return
				}

				c.Set("ID_TOKEN_CLAIMS", claims)
			}
		}

		c.Next()
	}
}
