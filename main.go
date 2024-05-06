package main

import (
	"errors"
	"log"

	"github.com/AdluAghnia/nyoba-fiber/connection"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name   string
	Passwd string
}

type Auth struct {
	Message           string
	IsAunthetificated bool
}

var auth Auth

func createUser(name, password string) User {
	return User{
		Name:   name,
		Passwd: password,
	}
}

func hashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func saveUserToDB(u User) error {
	conn, err := connection.InitiliazedDB()
	if err != nil {
		return err
	}

	_, err = conn.Exec("INSERT INTO User (Name, Passwd) VALUE (?,?)", u.Name, u.Passwd)
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}

func isUsernameExist(name string) (bool, error) {
	conn, err := connection.InitiliazedDB()
	if err != nil {
		return false, err
	}
	result, err := conn.Prepare("SELECT COUNT(*) FROM User WHERE Name = ?")
	if err != nil {
		log.Println("PREPARE")
		return false, err
	}
	defer result.Close()

	var count int
	err = result.QueryRow(name).Scan(&count)
	if err != nil {
		log.Println("QUERY")
		return false, err
	}

	if count != 0 {
		log.Printf("count : %v", count)
		return false, nil
	}

	return true, nil
}

func validateRegister(name, password string) (bool, error) {
	lenValid := len(name) >= 3 && len(password) >= 6
	nameExist, err := isUsernameExist(name)
	if !lenValid {
		return false, errors.New("password and username should atleast contain 6 characters")
	}

	if err != nil {
		return false, err
	}
	return nameExist && lenValid, nil
}

func validateLogin(name, password string) (bool, error) {
	var storedPassword []byte

	conn, err := connection.InitiliazedDB()
	if err != nil {
		return false, err
	}

	err = conn.QueryRow("SELECT Passwd FROM User WHERE Name = ?", name).Scan(&storedPassword)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(storedPassword, []byte(password))
	if err != nil {
		log.Println("Compare Failed")
		return false, err
	}
	return true, nil
}

func loginHandler(c *fiber.Ctx) error {
	// get User information from Form
	name := c.FormValue("name")
	password := c.FormValue("passwd")
	// Validate Login
	val, err := validateLogin(name, password)
	if err != nil {
		log.Println("Val Login")
		return err
	}

	if !val {
		auth.Message = "Login Gagal"
		return c.Redirect("/auth/login")
	}
	log.Printf("Login Valid is %v", val)
	auth.IsAunthetificated = true
	return c.Redirect("/")
}

func registerHandler(c *fiber.Ctx) error {
	var user User
	// Get User input from Form
	name := c.FormValue("name")
	password := c.FormValue("passwd")
	// Validate user input
	val, err := validateRegister(name, password)
	if err != nil {
		auth.Message = err.Error()
		return c.Redirect("/auth/register")
	}
	// Hash the password
	if val {
		log.Println("VALID")
		hashPasswd, err := hashPassword(password)
		if err != nil {
			return c.SendString(err.Error())
		}
		user = createUser(name, string(hashPasswd))

		// Store user data in the datatbase
		err = saveUserToDB(user)
		if err != nil {
			log.Println(err)
			return err
		}
		auth.Message = "Registrasi Berhasil"
		auth.IsAunthetificated = false

		return c.Redirect("/auth/login", fiber.StatusFound)
	}

	log.Printf("Name: %v,  pass : %v  \n", name, password)

	log.Printf("is not valid : %v", val)

	auth.Message = ""
	// Return a success message or error response
	return c.Redirect("/auth/register")
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		isLogin := auth.IsAunthetificated
		auth.IsAunthetificated = false

		return c.Render("front", fiber.Map{
			"Title":  "Hellow Wolrd",
			"Logged": isLogin,
		}, "layouts/main")
	})
	app.Get("/auth/login", func(c *fiber.Ctx) error {
		message := auth.Message
		auth.Message = ""
		return c.Render("login", fiber.Map{
			"Logged":  false,
			"Message": message,
		}, "layouts/main")
	})
	app.Post("/login", loginHandler)

	app.Get("/auth/register", func(c *fiber.Ctx) error {
		message := auth.Message
		auth.Message = ""
		return c.Render("register", fiber.Map{
			"Logging": false,
			"Message": message,
		}, "layouts/main")
	})
	app.Post("/register", registerHandler)
	log.Fatal(app.Listen(":6969"))
}
