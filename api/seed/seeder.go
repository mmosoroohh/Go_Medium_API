package seed

import (
	"github.com/jinzhu/gorm"
	"github.com/mmosoroohh/Go_Medium_API/api/models"
	"log"
)

var users = []models.User{
	models.User{
		Username: "mmosoroohh",
		Email:    "arnoldosoro@gmail.com",
		Password: "Password",
	},
	models.User{
		Username: "lutherjunior",
		Email:    "lutherjunior@gmail.com",
		Password: "Password",
	},
	models.User{
		Username: "osorobrian",
		Email:    "brianosoro@gmail.com",
		Password: "Password",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "Title 1",
		Content: "Content 1",
	},
	models.Post{
		Title:   "Title 2",
		Content: "Content 2",
	},
	models.Post{
		Title:   "Title 3",
		Content: "Content 3",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("Can't drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("Can't migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("Can't seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("Can't seed posts table: %v", err)
		}
	}
}
