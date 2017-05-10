package db

import (
	"poemcrawler/util"
	"log"
)

func SavePoem(p util.Poem) error {
	db, err := NewManager()
	defer db.Close()

	c := db.Session.DB(Name).C(Poems)

	err = c.Insert(p)
	if err != nil {
		log.Println("保存诗歌失败：", p.Title)
		return err
	}
	log.Println("保存诗歌成功：", p.Title)

	return nil
}

func SavePoems(ps []util.Poem) (n int, err error) {
	db, err := NewManager()
	defer db.Close()

	c := db.Session.DB(Name).C(Poems)

	n = 0
	for _, p := range ps {
		err = c.Insert(p)
		if err != nil {
			log.Println("保存诗歌失败：", p.Title)
			continue
		}
		n++
		log.Println("保存诗歌成功：", p.Title)
	}

	return n, err
}
