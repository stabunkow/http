package user_repository

import (
	"stabunkow/http/pkg/mongodb"
	"stabunkow/http/pkg/redis"
	"stabunkow/http/pkg/setting"
	"testing"
	"time"
)

var repo *UserRepository

func setup() {
	setting.Setup("../../configs/app.ini")
	mongodb.Setup()
	redis.Setup()

	repo = NewUserRepository(mongodb.GetDefaultDb(), redis.GetDefaultCache())
}

// func TestCreateUser(t *testing.T) {
// 	setup()
// 	t.Log("insert 50000 user.")
// 	t1 := time.Now()
// 	for i := 1; i < 50000; i++ {
// 		repo.CreateUser("test", "test")
// 	}
// 	t.Log("finished, processing time:", time.Since(t1))
// }

// func TestSyncUserRedis(t *testing.T) {
// 	setup()
// 	users, err := repo.GetAllUser()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log("sync users")
// 	t1 := time.Now()
// 	for _, user := range users {
// 		repo.SyncUserRedis(user.GetId())
// 	}
// 	t.Log("finished, processing time", time.Since(t1))
// }

func TestFindUserById(t *testing.T) {
	setup()

	id := "5d492d94e78d81e2362a367b"
	t.Log("find user by id:", id)
	t1 := time.Now()
	_, err := repo.FindUserByIdThourghMongodb(id)
	if err != nil {
		t.Error(err)
	}
	t.Log("finished, processing time", time.Since(t1))
}

func TestFindUserIdByIdThroughRedis(t *testing.T) {
	setup()

	id := "5d492d94e78d81e2362a367b"
	repo.SyncUserRedis(id)
	t.Log("find user by id(redis):", id)
	t1 := time.Now()
	_, err := repo.FindUserByIdThourghRedis(id)
	if err != nil {
		t.Error(err)
	}
	t.Log("finished, processing time", time.Since(t1))
}
