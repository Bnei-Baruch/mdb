package permissions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/casbin/casbin"
	"github.com/stretchr/testify/suite"

	"github.com/coreos/go-oidc"

	"github.com/Bnei-Baruch/mdb/utils"
)

type PermissionsSuite struct {
	suite.Suite
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPermissions(t *testing.T) {
	suite.Run(t, new(PermissionsSuite))
}

func (suite *PermissionsSuite) TestPermissions() {
	e := casbin.NewEnforcer()
	e.EnableLog(false)

	// load model
	e.SetModel(casbin.NewModel(permissionsModel))
	e.InitWithModelAndAdapter(casbin.NewModel(permissionsModel), NewEmbeddedPolicyAdapter())

	perms := [][]string{
		{"archive_admin", "data_public", "read"},
		{"archive_admin", "data_public", "write"},
		{"archive_admin", "data_public", "write_i18n"},
		{"archive_admin", "data_public", "metadata_write"},
		{"archive_admin", "data_sensitive", "read"},
		{"archive_admin", "data_sensitive", "write"},
		{"archive_admin", "data_sensitive", "i18n_write"},
		{"archive_admin", "data_sensitive", "metadata_write"},
		{"archive_admin", "data_private", "read"},
		{"archive_admin", "data_private", "write"},
		{"archive_admin", "data_private", "i18n_write"},
		{"archive_admin", "data_private", "metadata_write"},

		{"archive_editor", "data_public", "read"},
		{"archive_editor", "data_public", "write"},
		{"archive_editor", "data_public", "i18n_write"},
		{"archive_editor", "data_public", "metadata_write"},
		{"archive_editor", "data_sensitive", "read"},
		{"archive_editor", "data_sensitive", "write"},
		{"archive_editor", "data_sensitive", "i18n_write"},
		{"archive_editor", "data_sensitive", "metadata_write"},

		{"archive_tagger", "data_public", "read"},
		{"archive_tagger", "data_public", "i18n_write"},
		{"archive_tagger", "data_public", "metadata_write"},
		{"archive_tagger", "data_sensitive", "read"},
		{"archive_tagger", "data_sensitive", "i18n_write"},
		{"archive_tagger", "data_sensitive", "metadata_write"},

		{"archive_typist", "data_public", "read"},
		{"archive_typist", "data_sensitive", "read"},
		{"archive_typist", "data_private", "read"},

		{"archive_uploader", "data_public", "read"},
		{"archive_uploader", "data_sensitive", "read"},

		{"bb_user", "data_public", "read"},
	}
	for _, perm := range perms {
		suite.True(e.Enforce(utils.ConvertArgsString(perm)...), strings.Join(perm, ", "))
	}

	perms = [][]string{
		{"archive_editor", "data_private", "read"},
		{"archive_editor", "data_private", "write"},
		{"archive_editor", "data_private", "i18n_write"},
		{"archive_editor", "data_private", "metadata_write"},

		{"archive_tagger", "data_public", "write"},
		{"archive_tagger", "data_sensitive", "write"},
		{"archive_tagger", "data_private", "read"},
		{"archive_tagger", "data_private", "write"},
		{"archive_tagger", "data_private", "i18n_write"},
		{"archive_tagger", "data_private", "metadata_write"},

		{"archive_typist", "data_public", "write"},
		{"archive_typist", "data_public", "i18n_write"},
		{"archive_typist", "data_public", "metadata_write"},
		{"archive_typist", "data_sensitive", "write"},
		{"archive_typist", "data_sensitive", "i18n_write"},
		{"archive_typist", "data_sensitive", "metadata_write"},
		{"archive_typist", "data_private", "write"},
		{"archive_typist", "data_private", "i18n_write"},
		{"archive_typist", "data_private", "metadata_write"},

		{"archive_uploader", "data_public", "write"},
		{"archive_uploader", "data_public", "i18n_write"},
		{"archive_uploader", "data_public", "metadata_write"},
		{"archive_uploader", "data_sensitive", "write"},
		{"archive_uploader", "data_sensitive", "i18n_write"},
		{"archive_uploader", "data_sensitive", "metadata_write"},
		{"archive_uploader", "data_private", "read"},
		{"archive_uploader", "data_private", "write"},
		{"archive_uploader", "data_private", "i18n_write"},
		{"archive_uploader", "data_private", "metadata_write"},

		{"bb_user", "data_public", "write"},
		{"bb_user", "data_public", "write_i18n"},
		{"bb_user", "data_public", "metadata_write"},
		{"bb_user", "data_sensitive", "read"},
		{"bb_user", "data_sensitive", "write"},
		{"bb_user", "data_sensitive", "i18n_write"},
		{"bb_user", "data_sensitive", "metadata_write"},
		{"bb_user", "data_private", "read"},
		{"bb_user", "data_private", "write"},
		{"bb_user", "data_private", "i18n_write"},
		{"bb_user", "data_private", "metadata_write"},
	}
	for _, perm := range perms {
		suite.False(e.Enforce(utils.ConvertArgsString(perm)...), strings.Join(perm, ", "))
	}

}

func (suite *PermissionsSuite) TestOIDC() {
	provider, err := oidc.NewProvider(context.TODO(), "https://accounts.kab.info/auth/realms/main")
	//provider, err := oidc.NewProvider(context.TODO(), "https://accounts.kbb1.com/auth/realms/main")
	utils.Must(err)

	verifier := provider.Verifier(&oidc.Config{
		SkipClientIDCheck: true,
		SkipExpiryCheck:   true,
	})

	//token, err := verifier.Verify(context.TODO(), "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJxVHZYMkh3elFhbjVadlNlUHJtRWxkZE0zUFYzYUU0R1liRVFvSnJ3U2hZIn0.eyJqdGkiOiI4MmI3ODc0Mi1hMWE0LTQxNjAtYWMwYy0xMzkxNjQ0NDllOTgiLCJleHAiOjE1MTczMDE4OTMsIm5iZiI6MCwiaWF0IjoxNTE3MzAwOTkzLCJpc3MiOiJodHRwczovL2FjY291bnRzLmtiYjEuY29tL2F1dGgvcmVhbG1zL21haW4iLCJhdWQiOiJtZGItYWRtaW4tdWkiLCJzdWIiOiI3OTBkMTlmNy0yYjA2LTRhODgtYTU1YS1kZWM2NTdjYjk3MmIiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJtZGItYWRtaW4tdWkiLCJub25jZSI6Ijk2YzIyMGRhN2EwOTQ0OGVhYWVlNTA4YzQzNDFhMDU1IiwiYXV0aF90aW1lIjoxNTE3MzAwOTkzLCJzZXNzaW9uX3N0YXRlIjoiNWE5ODFhNGMtOWVkYy00OGY3LWFlNzEtNDJkZGQ0NGQ5MjA2IiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJ1bWFfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sIm5hbWUiOiJFZG8gU2hvciIsInByZWZlcnJlZF91c2VybmFtZSI6ImVkb3Nob3JAZ21haWwuY29tIiwiZ2l2ZW5fbmFtZSI6IkVkbyIsImZhbWlseV9uYW1lIjoiU2hvciIsImVtYWlsIjoiZWRvc2hvckBnbWFpbC5jb20ifQ.KSqdgkqNiVyVrHfxFejMXsVZpWI_ptHAmP5Ft_iSt0YVL_VrO3YMSQ9e1G8YVLRCOqh3GbI6_iLfav4ZU8wKIBrCSrN2VW8ckZ4l3Mk2urZdnrq-2Ai7xyJ0JMEDlvZek1le1whpolMIL09xqJiuY8JU7Io4ZO8iz__GmoZQhS6yCO5qrZTzXyJgQrSfk9mbrAJ_jZE86e8D8DHRuygNlUcZsbczS4Hu6Wa0g7oc7_ZlYKUk8Q5QHBYVWeWJMbVsp_NABrrtj_-nxtnHZ6mNry8jLlS7KSa-4vVT9CxaBl1tyREz-IW01074VKiS9Qab6JvmaVprwVYF4a7mAwXv_w")
	token, err := verifier.Verify(context.TODO(), "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJxVHZYMkh3elFhbjVadlNlUHJtRWxkZE0zUFYzYUU0R1liRVFvSnJ3U2hZIn0.eyJleHAiOjE2MDcyMzk1OTMsImlhdCI6MTYwNzIyNTE5MywiYXV0aF90aW1lIjoxNjA3MjI1MTkzLCJqdGkiOiI0OTA5NTMzOS1lZDQ1LTQxNWItOTY5ZS0yOTNhZDk0NTBmMGUiLCJpc3MiOiJodHRwczovL2FjY291bnRzLmthYi5pbmZvL2F1dGgvcmVhbG1zL21haW4iLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNzkwZDE5ZjctMmIwNi00YTg4LWE1NWEtZGVjNjU3Y2I5NzJiIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoibWRiLWFkbWluLXVpIiwibm9uY2UiOiIyMjc0YzE4NThkZTk0NmY3YjAzMTcwNWU5N2Q4ZTVjOSIsInNlc3Npb25fc3RhdGUiOiI3NjdlZGVlNi1iMzUwLTQ1MTctODE5Mi03ZjZmNzk1M2NlNGYiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIndmX2NlbnNvciIsImd4eV9zbmRtYW4iLCJneHlfdXNlciIsImd4eV9zaGlkdXIiLCJiYl91c2VyIiwid2ZfYWRtaW4iLCJneHlfcm9vdCIsIndmX2luZ2VzdCIsIm9mZmxpbmVfYWNjZXNzIiwidW1hX2F1dGhvcml6YXRpb24iLCJhcmNoaXZlX2FkbWluIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSIsIm5hbWUiOiJFZG8gU2hvciIsInByZWZlcnJlZF91c2VybmFtZSI6ImVkb3Nob3JAZ21haWwuY29tIiwiZ2l2ZW5fbmFtZSI6IkVkbyIsImxvY2FsZSI6ImVuIiwiZmFtaWx5X25hbWUiOiJTaG9yIiwiZW1haWwiOiJlZG9zaG9yQGdtYWlsLmNvbSJ9.ZN8K4g909Tmq3PxsiAcaRmA9wxnspwomOKvC2BOovOz_89SItaRWJ8OKoxMy4lDsgjjCS47zVB8_eIQ-CpL90g4Tu6XTYLbPUZnkRUE0FutOUWKZytloV_RHho1wFTbGgv68A5brXleEmOmZwCFo3zrhl_pM2rHo5jIrzr_dnpeL5lg0v-sMO5ydN2b8gT6Rm7J1Vfeb6fWmGs_waVHzMyxEiGk-sOosFEeAHNEI1vZ0GWknEzKzMFIZ5IiY3t6SBEqM5IvYy_juB88lUF5LFY0Zsbd1QvS0_oDLU8CkcNwt8xXJBDBqGo6YG4JhO42vuZg5Xec75tlBJzBd7kqD8w")
	utils.Must(err)

	var claims IDTokenClaims
	err = token.Claims(&claims)
	utils.Must(err)

	fmt.Printf("%v\n", claims)

	b, _ := json.Marshal(claims)
	fmt.Printf("%s\n", string(b))

}
