package user_repository

import (
	"stabunkow/http/models"
	"stabunkow/http/pkg/mongodb"
	"stabunkow/http/pkg/redis"
	"strconv"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"gopkg.in/mgo.v2/bson"
)

type UserRepository struct {
	db    *mongodb.Db
	cache *redis.Cache
}

func NewUserRepository(db *mongodb.Db, cache *redis.Cache) *UserRepository {
	return &UserRepository{
		db:    db,
		cache: cache,
	}
}

func (repo *UserRepository) CreateUser(openId string) (*models.User, error) {
	conn := repo.db.Conn()
	defer conn.Close()

	stoves := make([]*models.Stove, 0)
	for i := 0; i < 12; i++ {
		stove := &models.Stove{Idx: i, Level: 1, ChefIdx: -1}
		stoves = append(stoves, stove)
	}
	chefs := make([]*models.Chef, 0)
	for i := 0; i < 5; i++ {
		chef := &models.Chef{Idx: i, No: "1", Level: 1, StoveIdx: -1}
		chefs = append(chefs, chef)
	}

	user := &models.User{
		WechatOpenId: openId,
		Coins:        "0",
		Diamonds:     "0",
		Restaurant: &models.Restaurant{
			Level:           1,
			Scene:           1,
			ConveyBeltLevel: 1,
			Stoves:          stoves,
		},
		Chefs:     chefs,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	return user, conn.Collection("users").Insert(user)
}

func userKey(id string) string {
	return "user:" + id
}

func userRestaurantKey(id string) string {
	return userKey(id) + ":restaurant"
}

func userRestaurantStoveKey(id, idx string) string {
	return userRestaurantKey(id) + ":stove:" + idx
}

func userRestaurantStovesKey(id string) string {
	return userRestaurantKey(id) + ":stoves"
}

func userChefKey(id, idx string) string {
	return userKey(id) + ":chef:" + idx
}

func userChefsKey(id string) string {
	return userKey(id) + ":chefs"
}

func (repo *UserRepository) ExistsStore(id string) bool {
	conn := repo.cache.Conn()
	defer conn.Close()

	bol, _ := redigo.Bool(conn.Do("EXISTS", userKey(id)))
	return bol
}

func (repo *UserRepository) GetAllUser() ([]*models.User, error) {
	conn := repo.db.Conn()
	defer conn.Close()

	u := make([]*models.User, 0)

	err := conn.Collection("users").Find(nil).All(&u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (repo *UserRepository) FindUserById(id string) (*models.User, error) {
	if repo.ExistsStore(id) {
		return repo.FindUserByIdThourghRedis(id)
	} else {
		return repo.FindUserByIdThourghMongodb(id)
	}
}

func (repo *UserRepository) FindUserByIdThourghMongodb(id string) (*models.User, error) {
	conn := repo.db.Conn()
	defer conn.Close()

	user := &models.User{}
	err := conn.Collection("users").Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) FindUserByIdThourghRedis(id string) (*models.User, error) {
	conn := repo.cache.Conn()
	defer conn.Close()

	user := &models.User{}
	v, err := redigo.Values(conn.Do("HGETALL", userKey(id)))
	if err != nil {
		return nil, err
	}

	err = redigo.ScanStruct(v, user)
	if err != nil {
		return nil, err
	}

	chefs := make([]*models.Chef, 0)
	v, err = redigo.Values(conn.Do("SORT", userChefsKey(id),
		"BY", userChefKey(id, "*->idx"),
		"GET", userChefKey(id, "*->idx"),
		"GET", userChefKey(id, "*->no"),
		"GET", userChefKey(id, "*->level"),
		"GET", userChefKey(id, "*->stove_idx"),
	))
	if err != nil {
		return nil, err
	}
	err = redigo.ScanSlice(v, &chefs)
	if err != nil {
		return nil, err
	}
	user.Chefs = chefs

	v, err = redigo.Values(conn.Do("HGETALL", userRestaurantKey(id)))
	if err != nil {
		return nil, err
	}
	rat := &models.Restaurant{}
	err = redigo.ScanStruct(v, rat)
	if err != nil {
		return nil, err
	}
	user.Restaurant = rat

	stoves := make([]*models.Stove, 0)
	v, err = redigo.Values(conn.Do("SORT", userRestaurantStovesKey(id),
		"BY", userRestaurantStoveKey(id, "*->idx"),
		"GET", userRestaurantStoveKey(id, "*->idx"),
		"GET", userRestaurantStoveKey(id, "*->level"),
		"GET", userRestaurantStoveKey(id, "*->chef_idx"),
	))
	if err != nil {
		return nil, err
	}
	err = redigo.ScanSlice(v, &stoves)
	if err != nil {
		return nil, err
	}
	user.Restaurant.Stoves = stoves

	return user, nil
}

func (repo *UserRepository) FindUserByWechatOpenId(openid string) (*models.User, error) {
	conn := repo.db.Conn()
	defer conn.Close()

	user := &models.User{}
	err := conn.Collection("users").Find(bson.M{"wechat_open_id": openid}).Select(bson.M{"wechat_open_id": 1}).One(user)
	if err != nil {
		return nil, err
	}

	return repo.FindUserById(user.GetId())
}

func (repo *UserRepository) UpdateUserSidById(id, sid string) error {
	conn := repo.db.Conn()
	defer conn.Close()

	err := conn.Collection("users").Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{
		"$set": bson.M{
			"sid":        sid,
			"updated_at": time.Now().Unix(),
		},
	})

	if err != nil {
		return err
	}

	if repo.ExistsStore(id) {
		rconn := repo.cache.Conn()
		defer rconn.Close()

		_, err = rconn.Do("HSET", userKey(id), "sid", sid)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *UserRepository) UpdateUserSessionKeyById(id, sessionKey string) error {
	conn := repo.db.Conn()
	defer conn.Close()

	err := conn.Collection("users").Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{
		"$set": bson.M{
			"wechat_session_key": sessionKey,
			"updated_at":         time.Now().Unix(),
		},
	})

	if err != nil {
		return err
	}

	if repo.ExistsStore(id) {
		rconn := repo.cache.Conn()
		defer rconn.Close()

		_, err = rconn.Do("HSET", userKey(id), "wechat_session_key", sessionKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *UserRepository) SyncUserRedis(id string) error {
	if repo.ExistsStore(id) {
		return nil
	}

	user, err := repo.FindUserById(id)
	if err != nil {
		return err
	}

	conn := repo.cache.Conn()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HMSET", redigo.Args{}.Add(userKey(id)).AddFlat(user)...)
	conn.Send("HMSET", redigo.Args{}.Add(userRestaurantKey(id)).AddFlat(user.Restaurant)...)

	chefs := user.Chefs
	chefSets := make([]interface{}, 0)
	chefSets = append(chefSets, userChefsKey(id))
	for i := 0; i < len(chefs); i++ {
		idx := strconv.Itoa(i)
		conn.Send("HMSET", redigo.Args{}.Add(userChefKey(id, idx)).AddFlat(chefs[i])...)
		chefSets = append(chefSets, i)
	}
	conn.Send("SADD", chefSets...)

	stoves := user.Restaurant.Stoves
	stoveSets := make([]interface{}, 0)
	stoveSets = append(stoveSets, userRestaurantStovesKey(id))
	for i := 0; i < 12; i++ {
		idx := strconv.Itoa(i)
		conn.Send("HMSET", redigo.Args{}.Add(userRestaurantStoveKey(id, idx)).AddFlat(stoves[i])...)
		stoveSets = append(stoveSets, i)
	}
	conn.Send("SADD", stoveSets...)
	_, err = conn.Do("EXEC")

	return err
}
