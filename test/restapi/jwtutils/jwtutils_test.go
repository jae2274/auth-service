package jwtresolver

import (
	"testing"
	"time"

	"github.com/jae2274/auth-service/auth_service/restapi/jwtresolver"
	"github.com/stretchr/testify/require"
)

func TestJwtresolver(t *testing.T) {
	secretKey := "testKey"
	userId := "Jyo Liar"
	authorities := []string{"admin", "user"}

	now := time.Now()
	// accessTokenDuration := 30 * time.Minute

	jwtResolver := jwtresolver.NewJwtResolver(secretKey)
	// jwtResolver.SetAccessTokenDuration(accessTokenDuration)

	t.Run("return default expiration duration", func(t *testing.T) {
		//Given
		//When
		accessTokenDuration := jwtResolver.GetAccessTokenDuration()
		refreshTokenDuration := jwtResolver.GetRefreshTokenDuration()

		//Then
		require.Equal(t, 10*time.Minute, accessTokenDuration)
		require.Equal(t, 24*time.Hour, refreshTokenDuration)
	})

	t.Run("return custom expiration duration", func(t *testing.T) {
		//Given
		accessTokenDuration := 10 * time.Minute
		refreshTokenDuration := 48 * time.Hour

		//When
		jwtResolver.SetAccessTokenDuration(accessTokenDuration)
		jwtResolver.SetRefreshTokenDuration(refreshTokenDuration)

		//Then
		require.Equal(t, accessTokenDuration, jwtResolver.GetAccessTokenDuration())
		require.Equal(t, refreshTokenDuration, jwtResolver.GetRefreshTokenDuration())
	})

	t.Run("return valid claims", func(t *testing.T) {
		//Given
		tokenInfo, err := jwtResolver.CreateToken(userId, authorities, time.Now())
		require.NoError(t, err)

		//When
		claims, ok, err := jwtResolver.ParseToken(tokenInfo.AccessToken)

		//Then
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, userId, claims.UserId)
		require.Equal(t, authorities, claims.Authorities)
		require.Equal(t, now.Add(jwtResolver.GetAccessTokenDuration()).Unix(), claims.ExpiresAt.Time.Unix())
	})

	t.Run("return valid claims even empty authorities", func(t *testing.T) {
		//Given
		tokenInfo, err := jwtResolver.CreateToken(userId, []string{}, time.Now())
		require.NoError(t, err)

		//When
		claims, ok, err := jwtResolver.ParseToken(tokenInfo.AccessToken)

		//Then
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, userId, claims.UserId)
		require.Empty(t, claims.Authorities)
		require.Equal(t, now.Add(jwtResolver.GetAccessTokenDuration()).Unix(), claims.ExpiresAt.Time.Unix())
	})

	t.Run("return invalid claims if after expiresAt", func(t *testing.T) {
		//Given
		tokenInfo, err := jwtResolver.CreateToken(userId, authorities, now.Add(-jwtResolver.GetAccessTokenDuration()-time.Minute))
		require.NoError(t, err)

		//When
		_, ok, err := jwtResolver.ParseToken(tokenInfo.AccessToken)

		//Then
		require.NoError(t, err)
		require.False(t, ok)
	})

	//아래 두 케이스는 각 서비스간 CustomClaims 및 토큰 규약이 제대로 지켜지지 않았을 때 발생하므로 에러 발생으로 처리
	t.Run("return error when secret key is different", func(t *testing.T) {
		diffJwtResolver := jwtresolver.NewJwtResolver("differentSecretKey")
		//Given
		tokenInfo, err := diffJwtResolver.CreateToken(userId, authorities, time.Now())
		require.NoError(t, err)

		//When
		_, ok, err := jwtResolver.ParseToken(tokenInfo.AccessToken)

		//Then
		require.Error(t, err)
		require.False(t, ok)
	})

	t.Run("return error if empty userId", func(t *testing.T) {
		//Given
		tokenInfo, err := jwtResolver.CreateToken("", authorities, time.Now())
		require.NoError(t, err)

		//When
		_, ok, err := jwtResolver.ParseToken(tokenInfo.AccessToken)

		//Then
		require.Error(t, err)
		require.False(t, ok)
	})

	t.Run("return valid if refresh token is valid", func(t *testing.T) {
		//Given
		tokenInfo, err := jwtResolver.CreateToken(userId, authorities, time.Now())
		require.NoError(t, err)

		//When
		claims, ok, err := jwtResolver.ParseToken(tokenInfo.RefreshToken)

		//Then
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, userId, claims.UserId)
		require.Equal(t, now.Add(jwtResolver.GetRefreshTokenDuration()).Unix(), claims.ExpiresAt.Time.Unix())
	})

	t.Run("return invalid if refresh token is expired", func(t *testing.T) {
		//Given
		tokenInfo, err := jwtResolver.CreateToken(userId, authorities, now.Add(-jwtResolver.GetRefreshTokenDuration()-time.Minute))
		require.NoError(t, err)

		//When
		_, ok, err := jwtResolver.ParseToken(tokenInfo.RefreshToken)

		//Then
		require.NoError(t, err)
		require.False(t, ok)
	})
}
