package db

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/xmdhs/creditget/model"
)

type MysqlDb struct {
	DB *sqlx.DB
}

func NewMysql(url string) (*MysqlDb, error) {
	db, err := sqlx.Connect("mysql", url)
	if err != nil {
		return nil, fmt.Errorf("NewMysql: %w", err)
	}
	db.Exec(`CREATE TABLE credit (
		uid int(11) NOT NULL,
		name text NOT NULL,
		credits int(11) NOT NULL,
		extcredits1 int(11) NOT NULL,
		extcredits2 int(11) NOT NULL,
		extcredits3 int(11) NOT NULL,
		extcredits4 int(11) NOT NULL,
		extcredits5 int(11) NOT NULL,
		extcredits6 int(11) NOT NULL,
		extcredits7 int(11) NOT NULL,
		extcredits8 int(11) NOT NULL,
		oltime int(11) NOT NULL,
		groupname text NOT NULL,
		posts int(11) NOT NULL,
		threads int(11) NOT NULL,
		friends int(11) NOT NULL,
		views int(11) NOT NULL,
		medal int(11) NOT NULL,
		lastview bigint(20) NOT NULL,
		extgroupids text NOT NULL,
		sex TINYINT NOT NULL,
		PRIMARY KEY (uid) 
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
	  
	  CREATE TABLE config (
		id int(11) NOT NULL,
		value text NOT NULL,
		PRIMARY KEY (ID)
	  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;
	  `)

	return &MysqlDb{DB: db}, nil
}

func (m *MysqlDb) InsterCreditInfo(cxt context.Context, c *model.CreditInfo) error {
	_, err := m.DB.NamedExecContext(cxt, `INSERT or update INTO credit
	(uid, name, credits, extcredits1, extcredits2, extcredits3, extcredits4, extcredits5, extcredits6, extcredits7, extcredits8, oltime, groupname, posts, threads, friends, views, medal, lastview, extgroupids, sex)
	VALUES(:uid, :name, :credits, :extcredits1, :extcredits2, :extcredits3, :extcredits4, :extcredits5, :extcredits6, :extcredits7, :extcredits8, :oltime, :groupname, :posts, :threads, :friends, :views, :medal, :lastview, :extgroupids, :sex);
	`, c)
	if err != nil {
		return fmt.Errorf("InsterCreditInfo: %w", err)
	}
	return nil
}

func (m *MysqlDb) InsterConfig(cxt context.Context, c *model.Confing) error {
	_, err := m.DB.NamedExecContext(cxt, `INSERT or update INTO config
	(id, value)
	VALUES(:id, :value);
	`, c)
	if err != nil {
		return fmt.Errorf("InsterConfig: %w", err)
	}
	return nil
}

func (m *MysqlDb) SelectConfig(cxt context.Context, id int) (*model.Confing, error) {
	c := model.Confing{}
	err := m.DB.Get(&c, m.DB.Rebind(`SELECT id, value FROM config where id = ?;`), id)
	if err != nil {
		return nil, fmt.Errorf("SelectConfig: %w", err)
	}
	return &c, nil
}
