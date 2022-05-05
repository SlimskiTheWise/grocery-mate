package modeltests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/go-playground/assert.v1"
	"groceryMate/api/models"
)

//func TestFindAllFollows(t *testing.T) {
//
//	err := refreshUserAndFollowTable()
//	if err != nil {
//		log.Fatalf("Error refreshing user and post table %v\n", err)
//	}
//	_, _, err = seedUsersAndFollows()
//	if err != nil {
//		log.Fatalf("Error seeding user and post  table %v\n", err)
//	}
//	posts, err := postInstance.FindAllFollows(server.DB)
//	if err != nil {
//		t.Errorf("this is the error getting the posts: %v\n", err)
//		return
//	}
//	assert.Equal(t, len(*posts), 2)
//}

func TestCreateFollow(t *testing.T) {

	err := refreshUserAndFollowTable()
	if err != nil {
		log.Fatalf("Error user and follow refreshing table %v\n", err)
	}

	_, err = seedOneUserAndOneFollow()

	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newFollow := models.Follow{
		ID:         2,
		FollowerID: 1,
		FolloweeID: 2,
	}
	createdFollow, err := newFollow.CreateFollow(server.DB, int32(newFollow.FolloweeID))
	if err != nil {
		t.Errorf("this is the error getting the post: %v\n", err)
		return
	}
	assert.Equal(t, newFollow.ID, createdFollow.ID)
	assert.Equal(t, newFollow.FollowerID, createdFollow.FollowerID)
	assert.Equal(t, newFollow.FolloweeID, createdFollow.FolloweeID)
}

func TestGetFollowByID(t *testing.T) {

	err := refreshUserAndFollowTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	follow, err := seedOneUserAndOneFollow()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundFollow, err := followInstance.FindFollowByID(server.DB, follow.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundFollow.ID, follow.ID)

}

func TestDeleteAFollow(t *testing.T) {

	err := refreshUserAndFollowTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	follow, err := seedOneUserAndOneFollow()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := followInstance.DeleteAFollow(server.DB, follow.ID, follow.FollowerID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
