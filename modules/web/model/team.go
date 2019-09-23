package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	. "github.com/710leo/urlooker/modules/web/store"
)

type Team struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Resume      string    `json:"resume"`
	Creator     int64     `json:"creator"`
	Created     time.Time `json:"-" xorm:"<-"`
	CreatorName string    `json:"-" xorm:"-"`
}

func GetTeamById(id int64) (*Team, error) {
	obj := new(Team)
	has, err := Orm.Where("id=?", id).Get(obj)
	if err != nil {
		return obj, err
	}
	if !has {
		return obj, nil
	}
	return obj, nil
}

func AddTeam(name, resume string, creator int64, uids []int64) (int64, error) {
	if !(len(name) > 0 && creator > 0) {
		return 0, fmt.Errorf("团队信息有误")
	}

	session := Orm.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return 0, err
	}

	has, err := session.Where("name=?", name).Get(new(Team))
	if err != nil {
		return 0, err
	}
	if has {
		return 0, fmt.Errorf("团队名称已经被占用了")
	}

	team := Team{Name: name, Resume: resume, Creator: creator}
	_, err = session.Insert(&team)
	if err != nil {
		session.Rollback()
		return 0, err
	}

	err = addUsersIntoTeam(team.Id, uids, session)
	if err != nil {
		session.Rollback()
		return 0, err
	}

	return team.Id, session.Commit()
}

func RemoveTeamById(tid int64) error {
	session := Orm.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}

	_, err = session.Id(tid).Delete(new(Team))
	if err != nil {
		session.Rollback()
		return err
	}

	err = removeAllUsersFromTeam(tid, session)
	if err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (this *Team) Update(uids []int64) error {
	session := Orm.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}

	_, err = session.Id(this.Id).Update(this)
	if err != nil {
		session.Rollback()
		return err
	}

	err = updateUsersOfTeam(this.Id, uids, session)
	if err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	if err != nil {
		return err
	}

	return nil
}

func QueryTeams(query string, limit int) ([]*Team, error) {
	teams := make([]*Team, 0)
	if query == "" {
		return teams, nil
	}

	err := Orm.Where("name LIKE ?", "%"+query+"%").Limit(limit).Find(&teams)
	return teams, err
}

func GetTeamsByIds(ids string) ([]*Team, error) {
	teams := make([]*Team, 0)

	teamIdSlice := strings.Split(ids, ",")
	for _, tidStr := range teamIdSlice {
		tid, err := strconv.ParseInt(tidStr, 10, 64)
		if err != nil {
			continue
		}

		team, err := GetTeamById(tid)
		if err != nil {
			return teams, err
		}

		teams = append(teams, team)
	}
	return teams, nil
}
