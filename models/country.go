// https://sql2gorm.mccode.info/
package model

type Country struct {
	Idx  int    `gorm:"column:idx;primary_key;AUTO_INCREMENT"`
	Name string `gorm:"column:name;NOT NULL"`
}

func (m *Country) TableName() string {
	return "country"
}
