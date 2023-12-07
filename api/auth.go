package api

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	audience = os.Getenv("AUTH0_AUDIENCE")
	domain   = os.Getenv("AUTH0_DOMAIN")
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractTokenFromHeader(c.GetHeader("Authorization"))
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token, err := validateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Validate audience and issuer
		if !claims.VerifyAudience(audience, false) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid audience"})
			return
		}
		iss := fmt.Sprintf("https://%s/", domain)
		if !claims.VerifyIssuer(iss, false) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid issuer"})
			return
		}

		c.Next()
	}
}

func extractTokenFromHeader(authorizationHeader string) string {
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

func validateToken(tokenString string) (*jwt.Token, error) {
	key, err := getPublicKey()
	if err != nil {
		return nil, err
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
}

func getPublicKey() (*rsa.PublicKey, error) {
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", domain)
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			Alg string   `json:"alg"`
			Kid string   `json:"kid"`
			Kty string   `json:"kty"`
			N   string   `json:"n"`
			E   string   `json:"e"`
			Use string   `json:"use"`
			X5c []string `json:"x5c"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	var key *rsa.PublicKey
	for _, k := range jwks.Keys {
		if k.Kty == "RSA" && k.Alg == "RS256" {
			n, err := base64.RawURLEncoding.DecodeString(k.N)
			if err != nil {
				return nil, err
			}
			e, err := base64.RawURLEncoding.DecodeString(k.E)
			if err != nil {
				return nil, err
			}

			key = &rsa.PublicKey{
				N: (&big.Int{}).SetBytes(n),
				E: int((&big.Int{}).SetBytes(e).Int64()),
			}
			break
		}
	}

	if key == nil {
		return nil, errors.New("RSA public key not found in JWKS")
	}

	return key, nil
}
