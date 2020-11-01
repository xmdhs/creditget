package get

type Userinfo struct {
	Variables variables
}

type variables struct {
	Space space `json:"space"`
}

type space struct {
	Adminid      string      `json:"adminid"`
	Allowadmincp string      `json:"allowadmincp"`
	Avatarstatus string      `json:"avatarstatus"`
	Blacklist    string      `json:"blacklist"`
	Credits      string      `json:"credits"`
	Digestposts  string      `json:"digestposts"`
	Emailstatus  string      `json:"emailstatus"`
	Extcredits1  string      `json:"extcredits1"`
	Extcredits2  string      `json:"extcredits2"`
	Extcredits3  string      `json:"extcredits3"`
	Extcredits4  string      `json:"extcredits4"`
	Extcredits5  string      `json:"extcredits5"`
	Extcredits6  string      `json:"extcredits6"`
	Extcredits7  string      `json:"extcredits7"`
	Extcredits8  string      `json:"extcredits8"`
	Extgroupids  string      `json:"extgroupids"`
	Friends      string      `json:"friends"`
	Lastvisit    string      `json:"lastactivitydb"`
	Oltime       string      `json:"oltime"`
	Posts        string      `json:"posts"`
	Threads      string      `json:"threads"`
	UID          string      `json:"uid"`
	Username     string      `json:"username"`
	Views        string      `json:"views"`
	Medals       interface{} `json:"medals"`
	Group        group       `json:"group"`
}

type group struct {
	Grouptitle string `json:"grouptitle"`
}
