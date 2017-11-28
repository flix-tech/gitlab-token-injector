package main

import (
    "time"
    rsa "crypto/rsa"
    jwt "github.com/dgrijalva/jwt-go"
)

func generateToken(iat time.Time, ttl time.Duration, user string, groups []string, key *rsa.PrivateKey ) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
        "iat": iat.Unix(),
        "exp": iat.Add(ttl).Unix(),
        "user": user,
        "groups": groups,
    })
    log.Debugf("-- Generated token with payload: %#v", token.Claims)

    tokenString, err := token.SignedString(key)
    if err != nil {
        return "", err
    }
    log.Debug("-- Token signed with key.")
    return tokenString, nil
}
