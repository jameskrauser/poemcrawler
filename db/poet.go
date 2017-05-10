package db

import (
	"poemcrawler/util"
	"log"
)

func SavePoet(p util.Poet) error {
	db, err := NewManager()
	defer db.Close()

	c := db.Session.DB(Name).C(Poets)

	err = c.Insert(p)
	if err != nil {
		log.Println("保存诗人信息失败：", p.Name)
		return err
	}
	log.Println("保存诗人信息成功：", p.Name)

	return nil
}