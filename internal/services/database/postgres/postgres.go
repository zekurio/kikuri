package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
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

var _ database.Database = (*Postgres)(nil)
var guildTables = []string{
	"guilds",
	"permissions",
	"reddit",
	"redditsettings",
	"redditblocklist",
	"redditrules",
	"redditemotes"}

type tableColumn struct {
	Table  string
	Column string
}

var userTables = []tableColumn{
	{"apitokens", "userid"},
	{"refreshTokens", "userid"},
	{"users", "userid"},
}

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

func (p *Postgres) GetGuildAutoVoice(guildID string) (autovoices []string, err error) {
	chStr, err := GetValue[string](p, "guilds", "autovoice_ids", "guildid", guildID)
	if chStr == "" {
		return []string{}, err
	}

	return strings.Split(chStr, ";"), nil
}

func (p *Postgres) SetGuildAutoVoice(guildID string, channelIDs []string) error {
	return SetValue(p, "guilds", "autovoice_ids", strings.Join(channelIDs, ","), "guildid", guildID)
}

// PERMISSIONS

func (p *Postgres) GetPermissions(guildID string) (permissions map[string]perms.Array, err error) {
	permissions = make(map[string]perms.Array)
	rows, err := p.db.Query(`SELECT roleid, perms FROM permissions WHERE guildid = $1`, guildID)
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
		_, err := p.db.Exec(`DELETE FROM permissions WHERE guildid = $1 AND roleid = $2`, guildID, roleID)
		return err
	}

	pStr := strings.Join(permissions, ";")
	res, err := p.db.Exec(`UPDATE permissions SET perms = $1 WHERE guildid = $2 AND roleid = $3`, pStr, guildID, roleID)
	if err != nil {
		return err
	}
	ar, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ar == 0 {
		_, err := p.db.Exec(`INSERT INTO permissions (guildid, roleid, perms) VALUES ($1, $2, $3)`, guildID, roleID, pStr)
		return err
	}

	return nil
}

// VOTES

func (p *Postgres) GetVotes() (votes map[string]models.Vote, err error) {
	votes = make(map[string]models.Vote)
	rows, err := p.db.Query(`SELECT id, jsondata FROM votes`)
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
	_, err = p.db.Exec(`INSERT INTO votes (id, jsondata) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET jsondata = $2`, v.ID, rawData)
	return err
}

func (p *Postgres) DeleteVote(voteID string) error {
	_, err := p.db.Exec(`DELETE FROM votes WHERE id = $1`, voteID)
	return p.wrapErr(err)
}

// OAUTH2

func (p *Postgres) SetUserRefreshToken(ident, token string, expires time.Time) error {
	_, err := p.db.Exec(`INSERT INTO refreshtokens (ident, token, expires) VALUES ($1, $2, $3) ON CONFLICT (ident) DO UPDATE SET token = $2, expires = $3`, ident, token, expires)
	return p.wrapErr(err)
}

func (p *Postgres) GetUserByRefreshToken(token string) (ident string, expires time.Time, err error) {
	rows, err := p.db.Query(`SELECT ident, expires FROM refreshtokens WHERE token = $1`, token)
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
	_, err := p.db.Exec(`DELETE FROM refreshtokens WHERE ident = $1`, ident)
	return p.wrapErr(err)
}

// API TOKENS

func (p *Postgres) SetAPIToken(token models.APITokenEntry) error {
	res, err := p.db.Exec(`UPDATE apitokens 
	SET salt = $1, created = $2, expires = $3, lastaccess = $4, hits = $5 WHERE ident = $6`,
		token.Salt, token.Created, token.Expires, token.LastAccess, token.Hits, token.UserID)
	if err != nil {
		return err
	}

	ar, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if ar == 0 {
		_, err := p.db.Exec(`INSERT INTO apitokens (ident, salt, created, expires, lastaccess, hits) VALUES ($1, $2, $3, $4, $5, $6)`,
			token.UserID, token.Salt, token.Created, token.Expires, token.LastAccess, token.Hits)
		return err
	}

	return nil
}

func (p *Postgres) GetAPIToken(userID string) (models.APITokenEntry, error) {
	var token models.APITokenEntry
	err := p.db.QueryRow(`SELECT ident, salt, created, expires, lastaccess, hits FROM apitokens WHERE ident = $1`, userID).
		Scan(&token.UserID, &token.Salt, &token.Created, &token.Expires, &token.LastAccess, &token.Hits)
	return token, p.wrapErr(err)
}

func (p *Postgres) DeleteAPIToken(userID string) error {
	_, err := p.db.Exec(`DELETE FROM apitokens WHERE ident = $1`, userID)
	return p.wrapErr(err)
}

// REDDIT

func (p *Postgres) GetRedditKarma(userID, guildID string) (int, error) {
	return GetValue[int, string](p, "reddit", "karma", "userid", userID)
}

func (p *Postgres) GetRedditKarmaSum(userID string) (int, error) {
	var sum int
	err := p.db.QueryRow(`SELECT SUM(karma) FROM reddit WHERE userid = $1`, userID).Scan(&sum)
	return sum, p.wrapErr(err)
}

func (p *Postgres) GetRedditGuildEntries(guildID string, limit int) ([]models.GuildReddit, error) {
	var entries []models.GuildReddit
	rows, err := p.db.Query(`SELECT userid, karma FROM reddit WHERE guildid = $1 ORDER BY karma DESC LIMIT $2`, guildID, limit)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	for rows.Next() {
		var entry models.GuildReddit
		err := rows.Scan(&entry.UserID, &entry.Karma)
		if err != nil {
			return nil, p.wrapErr(err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (p *Postgres) SetRedditKarma(userID, guildID string, val int) error {
	return SetValue(p, "reddit", "karma", val, "userid", userID)
}

func (p *Postgres) UpdateRedditKarma(userID, guildID string, diff int) error {
	return p.tx(func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO reddit (userid, guildid, karma) VALUES ($1, $2, $3) ON CONFLICT (userid, guildid) DO UPDATE SET karma = reddit.karma + $3`, userID, guildID, diff)
		return err
	})
}

func (p *Postgres) SetRedditState(guildID string, state bool) error {
	return SetValue(p, "reddit", "state", state, "guildid", guildID)
}

func (p *Postgres) GetRedditState(guildID string) (bool, error) {
	return GetValue[bool, string](p, "reddit", "state", "guildid", guildID)
}

func (p *Postgres) SetRedditEmotes(guildID, emotesInc, emotesDec string) error {
	err := SetValue(p, "reddit", "emotesinc", emotesInc, "guildid", guildID)
	if err != nil {
		return err
	}

	return SetValue(p, "reddit", "emotesdec", emotesDec, "guildid", guildID)
}

func (p *Postgres) GetRedditEmotes(guildID string) (emotesInc, emotesDec string, err error) {
	emotesInc, err = GetValue[string, string](p, "reddit", "emotesinc", "guildid", guildID)

	if err != nil {
		return "", "", p.wrapErr(err)
	}

	emotesDec, err = GetValue[string, string](p, "reddit", "emotesdec", "guildid", guildID)
	if err != nil {
		return "", "", p.wrapErr(err)
	}

	return emotesInc, emotesDec, err
}

func (p *Postgres) SetRedditTokens(guildID string, tokens int) error {
	return SetValue(p, "reddit", "tokens", tokens, "guildid", guildID)
}

func (p *Postgres) GetRedditTokens(guildID string) (int, error) {
	return GetValue[int, string](p, "reddit", "tokens", "guildid", guildID)
}

func (p *Postgres) SetRedditPenalty(guildID string, state bool) error {
	return SetValue(p, "reddit", "penalty", state, "guildid", guildID)
}

func (p *Postgres) GetRedditPenalty(guildID string) (bool, error) {
	return GetValue[bool, string](p, "reddit", "penalty", "guildid", guildID)
}

func (p *Postgres) GetRedditBlockList(guildID string) ([]string, error) {
	var blockList []string
	rows, err := p.db.Query(`SELECT userid FROM redditblocklist WHERE guildid = $1`, guildID)

	if err != nil {
		return nil, p.wrapErr(err)
	}

	for rows.Next() {
		var userID string
		err := rows.Scan(&userID)
		if err != nil {
			return nil, p.wrapErr(err)
		}
		blockList = append(blockList, userID)
	}

	return blockList, nil
}

func (p *Postgres) IsRedditBlockListed(guildID, userID string) (bool, error) {
	return GetValue[bool, string](p, "redditblocklist", "guildid", "userid", guildID)
}

func (p *Postgres) AddRedditBlockList(guildID, userID string) error {
	_, err := p.db.Exec(`INSERT INTO redditblocklist (guildid, userid) VALUES ($1, $2)`, guildID, userID)
	return p.wrapErr(err)
}

func (p *Postgres) RemoveRedditBlockList(guildID, userID string) error {
	_, err := p.db.Exec(`DELETE FROM redditblocklist WHERE guildid = $1 AND userid = $2`, guildID, userID)
	return p.wrapErr(err)
}

func (p *Postgres) GetRedditRules(guildID string) ([]models.RedditRule, error) {
	var rules []models.RedditRule

	rows, err := p.db.Query(`SELECT id, trigger, value, action, argument FROM redditrules WHERE guildid = $1`, guildID)
	if err != nil {
		return nil, p.wrapErr(err)
	}

	for rows.Next() {
		var rule models.RedditRule
		err := rows.Scan(&rule.ID, &rule.Trigger, &rule.Value, &rule.Action, &rule.Argument)
		if err != nil {
			return nil, p.wrapErr(err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func (p *Postgres) CheckRedditRule(guildID, checksum string) (ok bool, err error) {
	err = p.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM redditrules WHERE guildid = $1 AND checksum = $2)`, guildID, checksum).Scan(&ok)
	return ok, p.wrapErr(err)
}

func (p *Postgres) AddOrUpdateRedditRule(rule models.RedditRule) error {
	return p.tx(func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO redditrules (id, guildid, trigger, value, action, argument, checksum) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (id) DO UPDATE SET trigger = $3, value = $4, action = $5, argument = $6, checksum = $7`,
			rule.ID, rule.GuildID, rule.Trigger, rule.Value, rule.Action, rule.Argument, rule.Checksum)
		return err
	})
}

func (p *Postgres) RemoveRedditRule(guildID string, id snowflake.ID) error {
	_, err := p.db.Exec(`DELETE FROM redditrules WHERE guildid = $1 AND id = $2`, guildID, id)
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
			_, err = tx.Exec(fmt.Sprintf(`DELETE FROM %s WHERE guildid = $1`, table), guildID)
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

func (p *Postgres) FlushUserData(userID string) error {
	return p.tx(func(tx *sql.Tx) error {
		var (
			err          error
			failedTables []string
		)

		for _, table := range guildTables {
			_, err = tx.Exec(fmt.Sprintf(`DELETE FROM %s WHERE userid = $1`, table), userID)
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
