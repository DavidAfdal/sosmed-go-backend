package response

import "github.com/labstack/echo/v4"

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, ApiResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}
