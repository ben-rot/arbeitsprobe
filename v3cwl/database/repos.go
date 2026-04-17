package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/qustavo/dotsql"

	"cwl/types"
)



type repository struct {
	db *sqlx.DB
	dot *dotsql.DotSql
}



func (r *repository) loadQuery(queryName string) string {
	query, err := r.dot.Raw(queryName)
	if err != nil {
		fmt.Printf("Could not load query: %v", err)
		return ""
	}

	return query
}




type UserRepo struct {
	repository
}




func (r *UserRepo) Upsert(discordUser *types.DiscordUser) (int, error) {

	query, err := r.dot.Raw("upsert-user")
	if err != nil {
		return 0, err
	}


	var userId int
	stmt, err := r.db.PrepareNamed(query)
    if err != nil {
        return 0, err
    }
	defer stmt.Close()


	err = stmt.QueryRow(discordUser).Scan(&userId)

	return userId, err
}




func (r *UserRepo) GetById(userId int) (*types.User, error) {

	query, err := r.dot.Raw("user-by-id")
	if err != nil {
		return nil, err
	}


	var user types.User
	err = r.db.Get(&user, query, userId)

	
	return &user, err
}



func (r *UserRepo) GetManagedFamilyId(userId int) (*int, error) {

	query, err := r.dot.Raw("check-manager")
	if err != nil {
        return nil, err
    }

	var familyId int
	err = r.db.Get(&familyId, query, userId)
	if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil 
        }
        return nil, err 
    }

	return &familyId, err
}




func (r *UserRepo) Manages(userId int) ([]int, error) {

	query := `SELECT family_id FROM managers WHERE user_id = $1` // JOIN for family details later

	var families []int

	err := r.db.Select(&families, query, userId)
	if err != nil {
		return nil, err
	}

	return families, nil
}

















type FamilyRepo struct {
	repository
}



func (r *FamilyRepo) GetById(familyId int) (*types.Family, error) {

	query := `SELECT id, discord_server_id, name, server_icon FROM families WHERE id = $1`

	var family types.Family
	err := r.db.Get(&family, query, familyId)
	if err != nil {
		return nil, err
	}
	
	return &family, nil
}


func (r *FamilyRepo) UpdateProfile(familyData *types.Family) (*types.Family, error) {
	
	return nil, nil
}





func (r *FamilyRepo) GetClans(familyId int) ([]types.Clan, error) {

	query := `SELECT id, tag, name, league_id FROM clans WHERE family_id = $1`

	var clans []types.Clan
	err := r.db.Select(&clans, query, familyId)
	if err != nil {
		return nil, err
	}

	return clans, nil
}



func (r *FamilyRepo) GetAccounts(familyId int) ([]types.Account, error) {
	
	query := `SELECT id, tag, name, townhall, COALESCE(clan_id, 0) AS clan_id, owner_id FROM accounts WHERE family_id = $1`

	var accounts []types.Account
	err := r.db.Select(&accounts, query, familyId)
	if err != nil {
		return nil, err
	}	

	return accounts, nil
}



func (r *FamilyRepo) GetManagers(familyId int, userId int) (int, error) {
	// Managed + Manager struct non existent yet

	return 0, nil
}




func (r *FamilyRepo) GetStatsForDashboard(familyId int) (*types.FamilyStats, error) {

	query, _ := r.dot.Raw("family-stats-for-dashboard")

	var familyStats types.FamilyStats
	err := r.db.Get(&familyStats, query, familyId)
	if err != nil {
		return nil, err
	}
	
	return &familyStats, nil
}




















type ClanRepo struct {
	repository
}



func (r *ClanRepo) GetById(clanId int) (*types.Clan, error) {

	query := `SELECT id, tag, name, league_id FROM clans WHERE id = $1`

	var clan types.Clan
	err := r.db.Get(&clan, query, clanId)
	if err != nil {
		return nil, err
	}

	return &clan, nil
}



func (r *ClanRepo) GetAccounts(clanData *types.Clan) ([]types.Account, error) {

	query := `SELECT id, tag, name, townhall FROM accounts WHERE clan_id = $1`

	var accounts []types.Account
	err := r.db.Select(&accounts, query, clanData.ID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}




func (r *ClanRepo) GetStatsForDashboard(familyId int) ([]types.ClanStats, error) {


	query, _ := r.dot.Raw("clan-stats-for-dashboard")

	var clansWithStats = make([]types.ClanStats, 0)
	err := r.db.Select(&clansWithStats, query, familyId)
	if err != nil {
		return nil, err
	}

	return clansWithStats, nil
}


















type AccountRepo struct {
	repository
}



func (r *AccountRepo) SanitizeTag(tag string) (string) {

	tag = strings.ToUpper(strings.TrimSpace(tag))
	return strings.TrimPrefix(tag, "#")
}


func (r *AccountRepo) Exists(tag string) (bool, error) {
	
	var exists bool
	query,_ := r.dot.Raw("exists")

	err := r.db.QueryRow(query, tag).Scan(&exists)
	if err != nil {
        fmt.Printf("Unique check failed for %s: %v", tag, err)
        return false, err
    }

	return exists, nil
}



func (r *AccountRepo) Save(accountData *types.Account) error {

	query,_ := r.dot.Raw("save")

	_, err := r.db.NamedExec(query, accountData)
	return err
}





func (r *AccountRepo) Delete(tag string) error {

	query := `DELETE FROM accounts WHERE tag = $1`

	_, err := r.db.Exec(query, tag)
	return err
}














func (r *AccountRepo) OfOwnerById(userId int) ([]types.Account, error) {

	query, err := r.dot.Raw("get-accounts-owner-by-id")
	if err != nil {
        return nil, err
    }

	var accounts = make([]types.Account, 0)
	err = r.db.Select(&accounts, query, userId)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}




func (r *AccountRepo) DashboardDataByOwnerId(userId int) ([]types.DashboardAccount, error) {

	query := r.loadQuery("get-accounts-dashboard-by-owner-id")

	var accounts []types.DashboardAccount
	err := r.db.Select(&accounts, query, userId)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}