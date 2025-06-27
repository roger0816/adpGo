package Adp

import (
	"fmt"
	"net/http"
	"strconv"

	C "adpGo/common"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func StartApiServer(port string) {
	http.HandleFunc("/", apiHandler)
	fmt.Println("API server listening on port", port)
	http.ListenAndServe(":"+port, nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from the RESTful API!")
}

func getIntParam(c *gin.Context, key string, defaultVal int) (int, error) {
	valStr := c.DefaultQuery(key, fmt.Sprintf("%d", defaultVal))
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s", key)
	}
	return val, nil
}

func RunApi(port string) {

	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {

		c.JSON(http.StatusOK, "testOk")
	})

	r.GET("/newCusId", func(c *gin.Context) {

		c.JSON(http.StatusOK, GetNewCusId())
	})

	r.GET("/cusId", func(c *gin.Context) {

		id, _ := strconv.Atoi(c.Query("id"))

		c.JSON(http.StatusOK, GenerateCusID(id))
	})

	r.GET("/parsingCusId", func(c *gin.Context) {
		id := c.Query("id")

		if id == "" {
			c.JSON(400, gin.H{"error": "缺少參數 id"})

		}

		n, err := ParseCusID(id)
		if err != nil {

			c.JSON(400, gin.H{"error": err.Error()})

		}

		c.JSON(200, gin.H{"id": id, "num": n})

	})

	r.GET("/api/GameList", func(c *gin.Context) {

		in := C.NewData()
		in.Action = C.QUERY_GAME_LIST
		query := AdpRecaller{}
		reData := query.ImplementRecall(in)

		c.JSON(http.StatusOK, reData)
	})

	r.Run(":" + port)

}
