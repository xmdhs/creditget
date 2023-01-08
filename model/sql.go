package model

type CreditInfo struct {
	Uid         int32  `json:"uid" db:"uid"`
	Name        string `json:"name" db:"name"`
	Credits     int32  `json:"credits" db:"credits"`
	Extcredits1 int32  `json:"extcredits1" db:"extcredits1"`
	Extcredits2 int32  `json:"extcredits2" db:"extcredits2"`
	Extcredits3 int32  `json:"extcredits3" db:"extcredits3"`
	Extcredits4 int32  `json:"extcredits4" db:"extcredits4"`
	Extcredits5 int32  `json:"extcredits5" db:"extcredits5"`
	Extcredits6 int32  `json:"extcredits6" db:"extcredits6"`
	Extcredits7 int32  `json:"extcredits7" db:"extcredits7"`
	Extcredits8 int32  `json:"extcredits8" db:"extcredits8"`
	Oltime      int32  `json:"oltime" db:"oltime"`
	Groupname   string `json:"groupname" db:"groupname"`
	Posts       int32  `json:"posts" db:"posts"`
	Threads     int32  `json:"threads" db:"threads"`
	Friends     int32  `json:"friends" db:"friends"`
	Medal       int32  `json:"medal" db:"medal"`
	Lastview    int64  `json:"lastview" db:"lastview"`
	Extgroupids string `json:"extgroupids" db:"extgroupids"`
	Sex         int32  `json:"sex" db:"sex"`
}

type Confing struct {
	ID    int    `db:"id"`
	VALUE string `db:"value"`
}

var CreditInfoFileds = []string{"uid", "name", "credits", "extcredits1", "extcredits2", "extcredits3", "extcredits4", "extcredits5", "extcredits6", "extcredits7", "extcredits8", "oltime", "groupname", "posts", "threads", "friends", "medal", "lastview", "extgroupids", "sex"}
