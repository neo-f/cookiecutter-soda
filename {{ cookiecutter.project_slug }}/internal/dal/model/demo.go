package model

type Demo struct {
	ID   string `gorm:"column:id;primaryKey;type:char(24);index;notnull"`
	Name string `gorm:"column:name;primaryKey;type:varchar(24);notnull;comment:name"`
}
