package middleware

import (
	"7hunt-be-rest-api/auth"
	"7hunt-be-rest-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenManager auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ดึง Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
			c.Abort()
			return
		}

		// 2. ตรวจสอบและตัด "Bearer " (ถ้ามี)
		// นิยมส่งมาในรูปแบบ "Bearer <token>"
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		// 3. ใช้ Manager ที่ฉีดเข้ามา (เรียกผ่าน Interface)
		claims, err := tokenManager.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrInvalidToken)
			c.Abort()
			return
		}

		// 4. เก็บข้อมูลไว้ใน Context เพื่อให้ Handler อื่นๆ เรียกใช้ได้
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
