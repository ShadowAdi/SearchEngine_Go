package db

import "time"

type SearchSettings struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	SearchOn  bool      `json:"searchOn"`
	AddNew    bool      `json:"addNew"`
	Amount    uint      `json:"amount"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *SearchSettings) Get() error {
	err := DBConn.Where("id = 1").First(s).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *SearchSettings) Update() error {
	tx := DBConn.Select("search_on", "add_new", "amount", "updatedAt").Where("id = 1").Updates(s)
	if tx.Error != nil {
		return tx.Error
	}
	return nil

}
