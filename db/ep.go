package db

import (
	"poemcrawler/util"
)

func SaveErrorPage(ep util.ErrorPage) error {
	db, err := NewManager()
	defer db.Close()

	c := db.Session.DB(Name).C(ErrorPages)

	err = c.Insert(ep)
	if err != nil {
		return err
	}

	return nil
}
