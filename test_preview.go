package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/mybatis-generator-gui-go/internal/api"
)

func main() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/preview", api.PreviewSnippet)

	reqBody := `{"tableName":"user","mapperName":"UserMapper","modelType":"com.example.model.User","snippetConfigs":[{"operation":"select","whereFields":[{"columnName":"id","fieldName":"id","javaType":"Long"}],"isBatch":false}, {"operation":"delete","whereFields":[{"columnName":"id","fieldName":"id","javaType":"Long"}],"isBatch":true}]}`
	req, _ := http.NewRequest("POST", "/preview", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())
}
