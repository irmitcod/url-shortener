package exception

import (
	"clotho/src/domain/web_response"
	"github.com/go-playground/validator/v10"
	"log"
)

func ErrorHandler(ctx echo.Context, err error) error {

	_, ok := err.(validator.ValidationErrors)
	if ok {
		return ctx.JSON(400, web_response.WebResponse{
			Code:   400,
			Status: "BAD_REQUEST",
			Data:   err.Error(),
		})
	}

	return ctx.JSON(500, web_response.WebResponse{
		Code:   500,
		Status: "INTERNAL_SERVER_ERROR",
		Data:   err.Error(),
	})
}

func ErrorEnvHandler(err error) {

	ve, ok := err.(validator.ValidationErrors)

	if ok {
		for _, fe := range ve {
			log.Panic(msgForTag(fe))
		}
	}
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " Environment variable is not set (You can set vale on .env file)"
	}
	return fe.Error() // default error
}
