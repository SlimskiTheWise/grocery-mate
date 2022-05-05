package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"groceryMate/api/controllers"
	"groceryMate/api/models"
)

var server = controllers.Server{}
var userInstance = models.User{}
var followInstance = models.Follow{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	refreshUserTable()

	user := models.User{
		Nickname:   "Steven victor",
		Email:      "steven@gmail.com",
		Password:   "password",
		Name:       "steve",
		Bio:        "here to make friends!",
		ProfilePic: "/img.jpeg",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() error {

	users := []models.User{
		models.User{
			Nickname:   "Steven victor",
			Email:      "steven@gmail.com",
			Password:   "password",
			Name:       "steve",
			Bio:        "here to make friends!",
			ProfilePic: "/img.jpeg",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		models.User{
			Nickname:   "Kenny Morris",
			Email:      "kenny@gmail.com",
			Password:   "password",
			Name:       "kenny",
			Bio:        "here to make friends!",
			ProfilePic: "/img.jpeg",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func refreshUserAndFollowTable() error {

	err := server.DB.DropTableIfExists(&models.User{}, &models.Follow{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Follow{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOneFollow() (models.Follow, error) {

	err := refreshUserAndFollowTable()
	if err != nil {
		return models.Follow{}, err
	}
	users := []models.User{
		models.User{
			Nickname:   "Steven victor",
			Email:      "steven@gmail.com",
			Password:   "password",
			Name:       "steve",
			Bio:        "here to make friends!",
			ProfilePic: "/img.jpeg",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		models.User{
			Nickname:   "Kenny Morris",
			Email:      "kenny@gmail.com",
			Password:   "password",
			Name:       "kenny",
			Bio:        "here to make friends!",
			ProfilePic: "/img.jpeg",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return followInstance, err
		}
	}

	follow := models.Follow{
		FollowerID: 1,
		FolloweeID: 2,
	}
	err = server.DB.Model(&models.Follow{}).Create(&follow).Error
	if err != nil {
		return models.Follow{}, err
	}
	return follow, nil
}

//func seedUsersAndFollows() ([]models.User, []models.Follow, error) {
//
//	var err error
//
//	if err != nil {
//		return []models.User{}, []models.Follow{}, err
//	}
//	var users = []models.User{
//		models.User{
//			Nickname: "Steven victor",
//			Email:    "steven@gmail.com",
//			Password: "password",
//		},
//		models.User{
//			Nickname: "Magu Frank",
//			Email:    "magu@gmail.com",
//			Password: "password",
//		},
//	}
//	var follows = []models.Follow{
//		models.Follow{
//			Title:   "Title 1",
//			Content: "Hello world 1",
//		},
//		models.Follow{
//			Title:   "Title 2",
//			Content: "Hello world 2",
//		},
//	}
//
//	for i, _ := range users {
//		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
//		if err != nil {
//			log.Fatalf("cannot seed users table: %v", err)
//		}
//		follows[i].AuthorID = users[i].ID
//
//		err = server.DB.Model(&models.Follow{}).Create(&follows[i]).Error
//		if err != nil {
//			log.Fatalf("cannot seed follows table: %v", err)
//		}
//	}
//	return users, follows, nil
//}
