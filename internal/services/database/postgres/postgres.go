package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/services/database/dberr"
	"github.com/zekurio/daemon/internal/util"
	"github.com/zekurio/daemon/internal/util/embedded"
	"github.com/zekurio/daemon/internal/util/vote"
	"github.com/zekurio/daemon/pkg/perms"
)

type Postgres struct {
	db *sql.DB
}

var (
	_           database.Database = (*Postgres)(nil)
	guildTables                   = []string{"guilds", "permissions"}
)

func NewPostgres(c config.DatabaseCreds) (*Postgres, error) {
	var (
		p   Postgres
		err error
	)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Username, c.Password, c.Database)
	p.db, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = p.db.Ping()
	if err != nil {
		return nil, err
	}

	goose.SetBaseFS(embedded.Migrations)
	err = goose.SetDialect("postgres")
	if err != nil {
		return nil, err
	}
	goose.SetLogger(log.StandardLog())
	err = goose.Up(p.db, "migrations")
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

// GUILDS

func (p *Postgres) GetGuildAutoRoles(guildID string) (autoroles []string, err error) {
	roleStr, err := GetValue[string](p, "guilds", "autorole_ids", "guild_id", guildID)
	if roleStr == "" {
		return []string{}, err
	}

	return strings.Split(roleStr, ";"), nil
}

func (p *Postgres) SetGuildAutoRoles(guildID string, roleIDs []string) error {
	return SetValue(p, "guilds", "autorole_ids", strings.Join(roleIDs, ";"), "guild_id", guildID)
}

func (p *Postgres) GetGuildAutoVoice(guildID string) (autovoices []string, err error) {
	chStr, err := GetValue[string](p, "guilds", "autovoice_ids", "guild_id", guildID)
	if chStr == "" {
		return []string{}, err
	}

	return strings.Split(chStr, ";"), nil
}

func (p *Postgres) SetGuildAutoVoice(guildID string, channelIDs []string) error {
	return SetValue(p, "guilds", "autovoice_ids", strings.Join(channelIDs, ","), "guild_id", guildID)
}

func (p *Postgres) GetGuildAPIEnabled(guildID string) (enabled bool, err error) {
	return GetValue[bool](p, "guilds", "api_enabled", "guild_id", guildID)
}

func (p *Postgres) SetGuildAPIEnabled(guildID string, enabled bool) error {
	return SetValue(p, "guilds", "api_enabled", enabled, "guild_id", guildID)
}

// PERMISSIONS

func (p *Postgres) GetPermissions(guildID string) (permissions map[string]perms.PermsArray, err error) {
	permissions = make(map[string]perms.PermsArray)
	rows, err := p.db.Query(`SELECT role_id, perms FROM permissions WHERE guild_id = $1`, guildID)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	for rows.Next() {
		var roleID string
		var permStr string

		err := rows.Scan(&roleID, &permStr)
		if err != nil {
			return nil, p.wrapErr(err)
		}

		permissions[roleID] = strings.Split(permStr, ";")
	}

	return
}

func (p *Postgres) SetPermissions(guildID, roleID string, permissions perms.PermsArray) error {

	if len(permissions) == 0 {
		_, err := p.db.Exec(`DELETE FROM permissions WHERE guild_id = $1 AND role_id = $2`, guildID, roleID)
		return err
	}

	pStr := strings.Join(permissions, ";")
	res, err := p.db.Exec(`UPDATE permissions SET perms = $1 WHERE guild_id = $2 AND role_id = $3`, pStr, guildID, roleID)
	if err != nil {
		return err
	}
	ar, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ar == 0 {
		_, err := p.db.Exec(`INSERT INTO permissions (guild_id, role_id, perms) VALUES ($1, $2, $3)`, guildID, roleID, pStr)
		return err
	}

	return nil
}

// VOTES

func (p *Postgres) GetVotes() (votes map[string]vote.Vote, err error) {
	votes = make(map[string]vote.Vote)
	rows, err := p.db.Query(`SELECT id, json_data FROM votes`)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	for rows.Next() {
		var voteID, rawData string
		err := rows.Scan(&voteID, &rawData)
		if err != nil {
			continue
		}
		vote, err := util.Unmarshal[vote.Vote](rawData)
		if err != nil {
			p.DeleteVote(rawData)
		} else {
			votes[vote.ID] = vote
		}

	}

	return
}

func (p *Postgres) AddUpdateVote(v vote.Vote) error {
	rawData, err := util.Marshal(v)
	if err != nil {
		return err
	}
	_, err = p.db.Exec(`INSERT INTO votes (id, json_data) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET json_data = $2`, v.ID, rawData)
	return err
}

func (p *Postgres) DeleteVote(voteID string) error {
	_, err := p.db.Exec(`DELETE FROM votes WHERE id = $1`, voteID)
	return p.wrapErr(err)
}

// OAUTH2

func (p *Postgres) SetUserRefreshToken(ident, token string, expires time.Time) error {
	_, err := p.db.Exec(`INSERT INTO refresh_tokens (ident, token, expires) VALUES ($1, $2, $3) ON CONFLICT (ident) DO UPDATE SET token = $2, expires = $3`, ident, token, expires)
	return p.wrapErr(err)
}

func (p *Postgres) GetUserByRefreshToken(token string) (ident string, expires time.Time, err error) {
	rows, err := p.db.Query(`SELECT ident, expires FROM refresh_tokens WHERE token = $1`, token)
	if err != nil {
		return "", time.Time{}, p.wrapErr(err)
	}

	if !rows.Next() {
		return "", time.Time{}, p.wrapErr(sql.ErrNoRows)
	}

	err = rows.Scan(&ident, &expires)
	if err != nil {
		return "", time.Time{}, p.wrapErr(err)
	}

	return ident, expires, nil
}

func (p *Postgres) RevokeUserRefreshToken(ident string) error {
	_, err := p.db.Exec(`DELETE FROM refresh_tokens WHERE ident = $1`, ident)
	return p.wrapErr(err)
}

// DATA MANAGEMENT

func (p *Postgres) FlushGuildData(guildID string) error {
	return p.tx(func(tx *sql.Tx) error {
		var (
			err          error
			failedGuilds []string
		)

		for _, table := range guildTables {
			_, err := tx.Exec(fmt.Sprintf(`DELETE FROM %s WHERE guild_id = $1`, table), guildID)
			if err != nil {
				failedGuilds = append(failedGuilds, guildID)
			}
		}

		if len(failedGuilds) > 0 || err != nil {
			return fmt.Errorf("failed to flush guild data for guilds: %v", failedGuilds)
		}

		return nil
	})
}

//
// HELPERS
//

func (p *Postgres) tx(f func(*sql.Tx) error) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	if err = f(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *Postgres) wrapErr(err error) error {
	if err != nil && err == sql.ErrNoRows {
		return dberr.ErrNotFound
	}
	return err
}

// GetValue retrieves a specific value from a PostgresSQL table.
func GetValue[TVal, TWv any](t *Postgres, table, valueKey, whereKey string, whereValue TWv) (TVal, error) {
	var value TVal
	query := fmt.Sprintf(`SELECT "%s" FROM %s WHERE "%s" = $1`, valueKey, table, whereKey)
	err := t.db.QueryRow(query, whereValue).Scan(&value)
	return value, t.wrapErr(err)
}

// SetValue updates a specific value in a PostgresSQL table, or inserts a new row if none is found.
func SetValue[TVal, TWv any](t *Postgres, table, valueKey string, value TVal, whereKey string, whereValue TWv) error {
	updateQuery := fmt.Sprintf(`UPDATE %s SET "%s" = $1 WHERE "%s" = $2`, table, valueKey, whereKey)
	result, err := t.db.Exec(updateQuery, value, whereValue)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		insertQuery := fmt.Sprintf(`INSERT INTO %s ("%s", "%s") VALUES ($1, $2)`, table, whereKey, valueKey)
		_, err = t.db.Exec(insertQuery, whereValue, value)
		if err == nil {
			return err
		}
	}

	return err
}
