package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"imgGo-Server/global"
	"imgGo-Server/model"
)

var orm *gorm.DB
var fOrm *gorm.DB

func GetDao() *gorm.DB {
	if orm != nil {
		return orm
	}
	return InitDao()
}

func InitDao() *gorm.DB {
	var err error

	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	conf := global.Conf.Database
	dsn := conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.DbName + "?charset=" + conf.CharSet + "&parseTime=True&loc=Local"
	orm, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败")
	}

	// Auto Migrate
	err = orm.AutoMigrate(
		&model.FileDbModel{},
		&model.ViewLogDbModel{},
		&model.UserDbModel{},
		&model.TokenDbModel{},
	)
	if err != nil {
		panic("Auto Migrate Failed!")
	}

	//// Set table options
	//err = orm.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&RecordItem{})
	//if err != nil {
	//	panic("Set table options Failed!")
	//}
	//
	//item := RecordItem{}
	//
	//// 插入
	//orm.Create(&item)
	//
	//// 查询
	//orm.Find(&item, "id = ?", 10)
	//
	//// 批量插入
	//var users = []RecordItem{item, item, item}
	//orm.Create(&users)

	return orm
}

func GetFilerunDao() *gorm.DB {
	if fOrm != nil {
		return fOrm
	}
	return InitFilerunDao()
}

func InitFilerunDao() *gorm.DB {
	var err error
	conf := global.Conf.FilerunDB
	dsn := conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.DbName + "?charset=" + conf.CharSet + "&parseTime=True&loc=Local"
	fOrm, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败")
	}
	return fOrm
}
