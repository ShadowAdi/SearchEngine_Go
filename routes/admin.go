package routes

import (
	"fmt"
	"os"
	"search_engine/db"
	"search_engine/utils"
	"search_engine/views"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type LoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func LoginHandler(c *fiber.Ctx) error {
	return render(c, views.Login())
}

func LoginPostHandler(c *fiber.Ctx) error {
	input := LoginForm{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Something Went Wrong</h2>")
	}
	user := &db.User{}
	user, err := user.LoginAdmin(input.Email, input.Password)
	if err != nil {
		c.Status(401)
		return c.SendString("<h2>Error: Unauthorized</h2>")
	}
	signedToken, err := utils.CreateNewAuthToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.Status(401)
		return c.SendString("<h2>Error: Something went 2wrong logging in</h2>")
	}

	cookie := fiber.Cookie{
		Name:     "admin",
		Value:    signedToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	c.Append("HX-Redirect", "/")
	return c.SendStatus(200)
}

type AdminClaims struct {
	User                 string `json:"user"`
	Id                   string `json:"id"`
	jwt.RegisteredClaims `claims:"claims"`
}

func LogoutHandler(c *fiber.Ctx) error {
	c.ClearCookie("admin")
	c.Set("HX-Redirect", "/login")
	return c.SendStatus(200)
}

func AuthMiddlewar(c *fiber.Ctx) error {
	cookie := c.Cookies("admin")
	if cookie == "" {
		return c.Redirect("/Login", 302)
	}
	token, err := jwt.ParseWithClaims(cookie, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return c.Redirect("/login", 302)
	}
	_, ok := token.Claims.(*AdminClaims)
	if ok && token.Valid {
		return c.Next()
	}
	return c.Redirect("/Login", 302)

}

func DashboardHandler(c *fiber.Ctx) error {
	settings := db.SearchSettings{}
	err := settings.Get()
	if err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Cannot Get Settings</h2>")
	}

	amount := strconv.FormatUint(uint64(settings.Amount), 10)
	return render(c, views.Home(amount, settings.SearchOn, settings.AddNew))
}

type settingsform struct {
	Amount   uint   `form:"amount"`
	SearchOn string `form:"searchOn"`
	AddNew   string `form:"addNew"`
}

func DashboardPostHandler(c *fiber.Ctx) error {
	input := settingsform{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		return c.SendString("<h2>Error: Cannot Update Settings</h2>")
	}
	fmt.Println(input.Amount)
	var AddNew bool = false
	var SearchOn bool = false

	if input.AddNew == "on" {
		AddNew = true
	}
	if input.SearchOn == "on" {
		SearchOn = true
	}

	settings := db.SearchSettings{}
	settings.Amount = uint(input.Amount)
	settings.SearchOn = SearchOn
	settings.AddNew = AddNew

	if err := settings.Update(); err != nil {
		fmt.Println(err)
		return c.SendString("<h2>Error: Cannot Update Settings</h2>")
	}
	c.Append("HX-Refresh", "true")
	return c.SendStatus(200)

}
