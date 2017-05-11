package db

import (
	"gopkg.in/mgo.v2"
)

// 数据库信息
const (
	Name       = "poemcrawler"                                         // 数据库名
	Poems      = "poems"                                               // 诗歌数据表
	Poets      = "poets"                                               // 诗人数据表
	ErrorPages = "errorpages"                                          // 解析有问题的页面数据表
	uri        = "mongodb://lelvAdmin:lelvAdmin@localhost:27017/admin" // 数据库地址uri
)

// DBManager 数据库管理器
type Manager struct {
	Session *mgo.Session
}

// NewManager 创建数据库管理器对象
func NewManager() (*Manager, error) {
	Session, err := mgo.Dial(uri)
	if err != nil {
		return nil, err
	}
	return &Manager{Session}, nil
}

// SetDB 根据数据库名字，创建数据库连接
func (m *Manager) SetDB(name string) *mgo.Database {
	return m.Session.DB(name)
}

// Coll 根据数据库表名，返回表对象
func (m *Manager) Coll(name string) *mgo.Collection {
	return m.Session.DB(Name).C(name)
}

// Close 关闭数据库
func (m *Manager) Close() {
	m.Session.Close()
}
