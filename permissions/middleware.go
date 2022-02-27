package permissions

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"
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
	Aud               interface{}      `json:"aud"`
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
		tokenVerifiers, _ := c.Get("TOKEN_VERIFIERS")
		if verifiers, ok := tokenVerifiers.([]*oidc.IDTokenVerifier); ok && verifiers != nil {
			// We have some ID Token Verifiers. Game on

			authHeader := strings.Split(strings.TrimSpace(c.Request.Header.Get("Authorization")), " ")
			if len(authHeader) == 2 || strings.ToLower(authHeader[0]) == "bearer" {
				// Authorization header provided, let's verify.

				token, err := verifyWithFallback(verifiers, authHeader[1])
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

				mdb := c.MustGet("MDB").(*sql.DB)
				user, err := getOrCreateUser(mdb, &claims)
				if err != nil {
					c.AbortWithError(http.StatusBadRequest, err).SetType(gin.ErrorTypePrivate)
					return
				}

				c.Set("USER", user)
			}
		}

		c.Next()
	}
}

func getOrCreateUser(mdb *sql.DB, claims *IDTokenClaims) (*models.User, error) {
	user, err := fetchUserFromDB(mdb, claims.Sub)
	if err != nil || user != nil {
		return user, err
	}

	user = &models.User{
		AccountID: null.StringFrom(claims.Sub),
		Email:     claims.Email,
		Name:      null.StringFrom(fmt.Sprintf("% %", claims.GivenName, claims.FamilyName)),
	}

	tx, err := mdb.Begin()
	utils.Must(err)

	if err := user.Insert(tx); err != nil {
		utils.Must(tx.Rollback())
	} else {
		utils.Must(tx.Commit())
	}
	return fetchUserFromDB(mdb, claims.Sub)
}

func fetchUserFromDB(mdb *sql.DB, accountID string) (*models.User, error) {
	user, err := models.Users(mdb, qm.Where("account_id = ?", accountID)).One()
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return user, nil
}

func verifyWithFallback(verifiers []*oidc.IDTokenVerifier, tokenStr string) (*oidc.IDToken, error) {
	var token *oidc.IDToken
	var err error
	for _, verifier := range verifiers {
		token, err = verifier.Verify(context.TODO(), tokenStr)
		if err == nil {
			return token, nil
		}
	}
	return nil, err
}
