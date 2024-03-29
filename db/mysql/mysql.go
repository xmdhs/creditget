package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/xmdhs/creditget/db"
	"github.com/xmdhs/creditget/model"
	"golang.org/x/exp/slices"
)

var _ db.DB = &MysqlDb{}

type MysqlDb struct {
	db *sqlx.DB
}

func NewMysql(cxt context.Context, url string) (*MysqlDb, error) {
	db, err := sqlx.ConnectContext(cxt, "mysql", url)
	if err != nil {
		return nil, fmt.Errorf("NewMysql: %w", err)
	}
	_, err = db.ExecContext(cxt, `create table if not exists credit (
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
		medal int(11) NOT NULL,
		lastview bigint(20) NOT NULL,
		extgroupids text NOT NULL,
		sex TINYINT NOT NULL,
		PRIMARY KEY (uid) 
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
	  `)
	if err != nil {
		return nil, fmt.Errorf("NewMysql: %w", err)
	}
	_, err = db.ExecContext(cxt, ` create table if not exists config (
		id int(11) NOT NULL,
		value text NOT NULL,
		PRIMARY KEY (ID)
	  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;
	  `)
	if err != nil {
		return nil, fmt.Errorf("NewMysql: %w", err)
	}
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	return &MysqlDb{db: db}, nil
}

func (m *MysqlDb) Begin(ctx context.Context, opts *sql.TxOptions) (*db.Tx, error) {
	tx, err := m.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("Begin: %w", err)
	}
	return db.NewTx(tx), nil
}

func (m *MysqlDb) BatchInsterCreditInfo(cxt context.Context, tx *sqlx.Tx, c []model.CreditInfo) error {
	_, err := tx.NamedExecContext(cxt, `REPLACE INTO credit
	(uid, name, credits, extcredits1, extcredits2, extcredits3, extcredits4, extcredits5, extcredits6, extcredits7, extcredits8, oltime, groupname, posts, threads, friends, medal, lastview, extgroupids, sex)
	VALUES(:uid, :name, :credits, :extcredits1, :extcredits2, :extcredits3, :extcredits4, :extcredits5, :extcredits6, :extcredits7, :extcredits8, :oltime, :groupname, :posts, :threads, :friends, :medal, :lastview, :extgroupids, :sex);
	`, c)
	if err != nil {
		return fmt.Errorf("InsterCreditInfo: %w", err)
	}
	return nil
}

func (m *MysqlDb) InsterConfig(cxt context.Context, tx *sqlx.Tx, c *model.Confing) error {
	_, err := tx.NamedExecContext(cxt, `REPLACE INTO config
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
	err := m.db.GetContext(cxt, &c, m.db.Rebind(`SELECT id, value FROM config where id = ?;`), id)
	if err != nil {
		return nil, fmt.Errorf("SelectConfig: %w", err)
	}
	return &c, nil
}

func (m *MysqlDb) GetCreditInfo(cxt context.Context, uid int) (*model.CreditInfo, error) {
	c := model.CreditInfo{}
	err := m.db.GetContext(cxt, &c, m.db.Rebind(`SELECT * from credit where uid = ?`), uid)
	if err != nil {
		return nil, fmt.Errorf("GetCreditInfo: %w", err)
	}
	return &c, nil
}

var ErrNotVaildFidld = errors.New("无效的字段")

func (m *MysqlDb) GetRank(cxt context.Context, uid int, field string) (int, error) {
	if !slices.Contains(model.CreditInfoFileds, field) {
		return 0, fmt.Errorf("GetRank: %w", ErrNotVaildFidld)
	}

	i := 0
	err := m.db.GetContext(cxt, &i, m.db.Rebind(`SELECT COUNT(*) FROM credit WHERE `+field+` > (SELECT `+field+` FROM credit WHERE uid = ?);`), uid)
	if err != nil {
		return 0, fmt.Errorf("GetRank: %w", err)
	}
	return i + 1, nil
}

func (m *MysqlDb) GetRanks(cxt context.Context, field string, limit, offset int, desc bool) ([]model.CreditInfo, error) {
	if !slices.Contains(model.CreditInfoFileds, field) {
		return nil, fmt.Errorf("GetRanks: %w", ErrNotVaildFidld)
	}
	word := "ASC"
	if desc {
		word = "DESC"
	}

	ml := make([]model.CreditInfo, 0, limit)
	err := m.db.SelectContext(cxt, &ml, m.db.Rebind(`SELECT * FROM credit ORDER BY `+field+` `+word+` LIMIT ? OFFSET ?`), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetRanks: %w", err)
	}
	return ml, nil
}

func (m *MysqlDb) GetSum(cxt context.Context) (int, error) {
	i := 0
	err := m.db.GetContext(cxt, &i, `SELECT COUNT(*) FROM credit`)
	if err != nil {
		return 0, fmt.Errorf("GetSum: %w", err)
	}
	return i, nil
}

func (m *MysqlDb) GetAvailableUserSum(cxt context.Context) (int, error) {
	i := 0
	err := m.db.GetContext(cxt, &i, `SELECT COUNT(*) FROM credit WHERE NOT lastview = 0`)
	if err != nil {
		return 0, fmt.Errorf("GetAvailableUserSum: %w", err)
	}
	return i, nil
}

func (m *MysqlDb) GetNilSum(cxt context.Context, field string) (int, error) {
	if !slices.Contains(model.CreditInfoFileds, field) {
		return 0, fmt.Errorf("GetNilSum: %w", ErrNotVaildFidld)
	}
	i := 0
	err := m.db.GetContext(cxt, &i, `SELECT COUNT(*) FROM credit WHERE `+field+` = 0 AND NOT lastview = 0`)
	if err != nil {
		return 0, fmt.Errorf("GetNilSum: %w", err)
	}
	return i, nil
}

func (m *MysqlDb) GetGroupname(cxt context.Context) ([]string, error) {
	sl := []string{}

	err := m.db.SelectContext(cxt, &sl, `SELECT DISTINCT groupname FROM credit`)
	if err != nil {
		return nil, fmt.Errorf("GetGroupname: %w", err)
	}
	return sl, nil
}

func (m *MysqlDb) GetGroupCount(cxt context.Context, groupname string) (int, error) {
	i := 0
	like := "%" + groupname + "%"
	err := m.db.GetContext(cxt, &i, m.db.Rebind(`SELECT COUNT(*) FROM credit WHERE groupname = ? OR extgroupids LIKE ?`), groupname, like)
	if err != nil {
		return 0, fmt.Errorf("GetGroupCount: %w", err)
	}
	return i, nil
}
