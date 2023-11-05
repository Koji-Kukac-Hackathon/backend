package response

import "github.com/gin-gonic/gin"

func Success(message string, datas ...interface{}) gin.H {
	var data interface{}
	if len(datas) > 0 {
		data = datas[0]
	}

	return gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	}
}

func Error(message string) gin.H {
	return gin.H{
		"status":  "error",
		"message": message,
	}
}
