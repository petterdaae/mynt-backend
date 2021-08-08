package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	utils "mynt/internal/utils"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
)

type Jwks struct {
	Keys []jsonWebKeys `json:"keys"`
}

type jsonWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func AuthGuard(database *utils.Database) gin.HandlerFunc {
	middleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := os.Getenv("AUTH0_CLIENT_ID")
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			// Verify 'iss' claim
			iss := os.Getenv("AUTH0_API_DOMAIN")
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	return func(c *gin.Context) {
		if err := middleware.CheckJWT(c.Writer, c.Request); err != nil {
			c.AbortWithStatus(401)
		}
	}
}

type CustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// TODO : Clean up this mess
func GetUserId(c *gin.Context) (int64, error) {
	// Parse id token
	authHeaderParts := strings.Split(c.Request.Header.Get("Authorization"), " ")
	tokenString := authHeaderParts[1]
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})
	claims, ok := token.Claims.(*CustomClaims)

	// TODO : do we have to check if the token is valid here?
	if !ok {
		c.AbortWithStatus(500)
		return 0, errors.New("Invalid token or failed to parse token.")
	}

	// Check if email is in database
	database, _ := c.MustGet("database").(*utils.Database)
	email := claims.Email
	connection, err := database.Connect()
	if err != nil {
		c.AbortWithStatus(500)
		return 0, err
	}
	defer connection.Close()

	rows, err := connection.Query(`SELECT id FROM users WHERE email = $1`, email)
	if err != nil {
		c.AbortWithStatus(500)
		return 0, err
	}
	defer rows.Close()

	var id int64
	count := 0
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			c.AbortWithStatus(500)
			return 0, err
		}
		count++
	}

	if count == 0 {
		err := connection.QueryRow(`INSERT INTO users (email) VALUES ($1) RETURNING id`, email).Scan(&id)
		if err != nil {
			c.AbortWithStatus(500)
			return 0, err
		}
	}

	return id, nil
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(os.Getenv("AUTH0_API_DOMAIN") + ".well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}
