package RpkNetwork

import (
	"fmt"
	"net/http"

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

func RunApi(port string, recaller Recaller) {

	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "testOk")
	})

	r.GET("/api/GameList", func(c *gin.Context) {

		in := NewData()
		in.Action = C.QUERY_GAME_LIST

		reData := recaller.ImplementRecall(in)

		c.JSON(http.StatusOK, reData)
	})

	r.Run(":" + port)

}
