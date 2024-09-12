package pages

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mhtml/h"
	"time"
)

func Test(ctx *fiber.Ctx) *h.Page {
	time.Sleep(time.Second * 1)
	text := fmt.Sprintf("News ID: %s", ctx.Params("id"))
	return h.NewPage(
		h.Div(h.Text(text)),
	)
}
