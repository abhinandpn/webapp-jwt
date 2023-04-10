package models

type User struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	Email         string `gorm:"unique"`
	Password      string
	Block_status  bool
	delete_status bool
}

/*

If use gorm.model for the user struct then we can able to use these data

 GORM MODEL

  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`

*/
