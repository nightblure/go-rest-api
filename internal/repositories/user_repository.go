package repos

import (
	"goapi/internal/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Name string
	Age  int
}

type UserRepository struct {
	url string
}

func CreateUserRepository() UserRepository {
	return UserRepository{url: "file::memory:?cache=shared"}
}

func (r *UserRepository) GetConn() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(r.url), &gorm.Config{})
	if err != nil {
		panic("error while try open db connection")
	}

	return db
}

type PaginationResponse struct {
	Page       int64         `json:"page"`
	Size       int64         `json:"size"`
	TotalCount int64         `json:"totalCount"`
	PagesCount int64         `json:"pagesCount"`
	Items      []entity.User `json:"items"`
}

func (r *UserRepository) GetUsers(page, size int64) PaginationResponse {
	var users []UserModel

	if page == 0 {
		page = 1
	}

	if size == 0 {
		size = 5
	}

	var pagesCount int64 = 0
	var totalCount int64
	r.GetConn().Model(&UserModel{}).Count(&totalCount)

	if page > 0 && size > 0 {
		pagesCount = totalCount / size

		if totalCount%size != 0 {
			pagesCount = pagesCount + 1
		}

		offset := int((size * page) - size)

		r.GetConn().Limit(int(size)).Offset(offset).Find(&users)
	} else {
		r.GetConn().Find(&users)
	}

	var usersOut []entity.User

	for _, u := range users {
		usersOut = append(usersOut, entity.User{Name: u.Name, Age: u.Age})
	}

	return PaginationResponse{
		Page:       page,
		Size:       size,
		PagesCount: pagesCount,
		TotalCount: totalCount,
		Items:      usersOut,
	}
}

func (r *UserRepository) GetUser(userName string) entity.User {
	var user UserModel
	r.GetConn().Where(&UserModel{Name: userName}).Find(&user)
	return entity.User{Name: user.Name, Age: user.Age}
}

func FillDb() {
	var r = CreateUserRepository()
	var db = r.GetConn()
	db.AutoMigrate(&UserModel{})

	var users = entity.CreateMockBatch()
	var dbUsers []UserModel

	for _, user := range users {
		dbUsers = append(dbUsers, UserModel{Name: user.Name, Age: user.Age})
	}

	db.Create(&dbUsers)
}
