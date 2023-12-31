package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zekurio/kikuri/internal/embedded"
	"github.com/zekurio/kikuri/internal/models"

	"github.com/charmbracelet/log"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/zekurio/kikuri/internal/services/database"
	"github.com/zekurio/kikuri/internal/services/database/dberr"
	"github.com/zekurio/kikuri/internal/util"
	"github.com/zekurio/kikuri/pkg/perms"
)

type Postgres struct {
	db *sql.DB
}

var (
	_           database.Database = (*Postgres)(nil)
	guildTables                   = []string{"guilds", "permissions"}
)

func NewPostgres(c models.DatabaseCreds) (*Postgres, error) {
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

// PERMISSIONS

func (p *Postgres) GetPermissions(guildID string) (permissions map[string]perms.Array, err error) {
	permissions = make(map[string]perms.Array)
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

func (p *Postgres) SetPermissions(guildID, roleID string, permissions perms.Array) error {
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

func (p *Postgres) GetVotes() (votes map[string]models.Vote, err error) {
	votes = make(map[string]models.Vote)
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
		vote, err := util.Unmarshal[models.Vote](rawData)
		if err != nil {
			p.DeleteVote(rawData)
		} else {
			votes[vote.ID] = vote
		}

	}

	return
}

func (p *Postgres) AddUpdateVote(v models.Vote) error {
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

// API TOKENS

func (p *Postgres) SetAPIToken(token models.APITokenEntry) error {
	res, err := p.db.Exec(`UPDATE apitokens 
	SET salt = $1, created = $2, expires = $3, lastaccess = $4, hits = $5 WHERE user_id = $6`,
		token.Salt, token.Created, token.Expires, token.LastAccess, token.Hits, token.UserID)
	if err != nil {
		return err
	}

	ar, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if ar == 0 {
		_, err := p.db.Exec(`INSERT INTO apitokens (user_id, salt, created, expires, lastaccess, hits) VALUES ($1, $2, $3, $4, $5, $6)`,
			token.UserID, token.Salt, token.Created, token.Expires, token.LastAccess, token.Hits)
		return err
	}

	return nil
}

func (p *Postgres) GetAPIToken(userID string) (models.APITokenEntry, error) {
	var token models.APITokenEntry
	err := p.db.QueryRow(`SELECT user_id, salt, created, expires, lastaccess, hits FROM apitokens WHERE user_id = $1`, userID).
		Scan(&token.UserID, &token.Salt, &token.Created, &token.Expires, &token.LastAccess, &token.Hits)
	return token, p.wrapErr(err)
}

func (p *Postgres) DeleteAPIToken(userID string) error {
	_, err := p.db.Exec(`DELETE FROM apitokens WHERE user_id = $1`, userID)
	return p.wrapErr(err)
}

// DATA MANAGEMENT

func (p *Postgres) FlushGuildData(guildID string) error {
	return p.tx(func(tx *sql.Tx) error {
		var (
			err          error
			failedTables []string
		)

		for _, table := range guildTables {
			_, err = tx.Exec(fmt.Sprintf(`DELETE FROM %s WHERE guild_id = $1`, table), guildID)
			if err != nil {
				failedTables = append(failedTables, table)
			}
		}

		if len(failedTables) > 0 {
			return fmt.Errorf("failed to flush tables: %s", strings.Join(failedTables, ", "))
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
