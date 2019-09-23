package model

import (
	"github.com/go-xorm/xorm"

	. "github.com/710leo/urlooker/modules/web/store"
)

type RelTeamUser struct {
	Id  int64 `json:"id"`
	Tid int64 `json:"tid"`
	Uid int64 `json:"uid"`
}

func TeamsOfUser(query string, uid int64, limit, offset int) ([]*Team, error) {
	teams := make([]*Team, 0)

	var err error
	if query != "" {
		tn := "%" + query + "%"
		err = Orm.Sql("SELECT * FROM team WHERE name LIKE ? AND ( id IN (SELECT DISTINCT tid FROM rel_team_user WHERE uid=? ) OR creator= ?)  ORDER BY name LIMIT ?,?", tn, uid, uid, offset, limit).Find(&teams)
	} else {
		err = Orm.Sql("SELECT * FROM team WHERE (id IN (SELECT DISTINCT tid FROM rel_team_user WHERE uid=?)) OR creator = ?  ORDER BY name LIMIT ?,?", uid, uid, offset, limit).Find(&teams)
	}

	return teams, err
}

func TeamCountOfUser(query string, uid int64) (int64, error) {
	if query != "" {
		tn := "%" + query + "%"
		return Orm.Where("name LIKE ? AND ( id IN (SELECT DISTINCT tid FROM rel_team_user WHERE uid=? ))", tn, uid).Count(new(Team))
	} else {
		return Orm.Where("id IN (SELECT tid FROM rel_team_user WHERE uid=? )", uid).Count(new(Team))
	}
}

func UsersOfTeam(tid int64) ([]*User, error) {
	users := make([]*User, 0)

	err := Orm.Cols("id", "name", "cnname", "email", "phone", "wechat").Sql("SELECT * FROM user WHERE id IN ( SELECT uid FROM rel_team_user WHERE tid=? )", tid).Find(&users)
	if err != nil {
		return users, err
	}

	return users, nil
}

func UsersInfoOfTeam(tid int64) ([]*User, error) {
	users := make([]*User, 0)

	err := Orm.Cols("name", "cnname", "email", "phone", "wechat").Sql("SELECT * FROM user WHERE id IN ( SELECT uid FROM rel_team_user WHERE tid=? )", tid).Find(&users)
	if err != nil {
		return users, err
	}

	return users, nil
}

func IsMemberOfTeam(uid, tid int64) (bool, error) {
	users, err := UsersOfTeam(tid)
	if err != nil {
		return false, err
	}

	for _, user := range users {
		if user != nil && user.Id == uid {
			return true, nil
		}
	}

	return false, nil
}

func IsCreatorOfTeam(uid, tid int64) (bool, error) {
	team, err := GetTeamById(tid)
	if err != nil {
		return false, err
	}

	if team.Creator == uid {
		return true, nil
	}
	return false, nil
}

func updateUsersOfTeam(tid int64, uids []int64, session *xorm.Session) error {
	err := removeAllUsersFromTeam(tid, session)
	if err != nil {
		return err
	}

	err = addUsersIntoTeam(tid, uids, session)
	if err != nil {
		return err
	}

	return err
}

func addUsersIntoTeam(tid int64, uids []int64, session *xorm.Session) error {
	for _, uid := range uids {
		relTeamUser := &RelTeamUser{Tid: tid, Uid: uid}
		_, err := session.Insert(relTeamUser)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeAllUsersFromTeam(tid int64, session *xorm.Session) error {
	_, err := session.Where("tid=?", tid).Delete(new(RelTeamUser))
	return err
}
