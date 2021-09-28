package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var inviteSearchColumns = []string{"name", "email"}

// NewInvite generates a new invitation.
func (d *Datastore) NewInvite(ctx context.Context, i *fleet.Invite) (*fleet.Invite, error) {
	if err := fleet.ValidateRole(i.GlobalRole.Ptr(), i.Teams); err != nil {
		return nil, err
	}

	err := d.withRetryTxx(ctx, func(tx sqlx.ExtContext) error {
		sqlStmt := `
	INSERT INTO invites ( invited_by, email, name, position, token, sso_enabled, global_role )
	  VALUES ( ?, ?, ?, ?, ?, ?, ?)
	`

		result, err := tx.ExecContext(ctx, sqlStmt, i.InvitedBy, i.Email,
			i.Name, i.Position, i.Token, i.SSOEnabled, i.GlobalRole)
		if err != nil && isDuplicate(err) {
			return alreadyExists("Invite", i.Email)
		} else if err != nil {
			return errors.Wrap(err, "create invite")
		}

		id, _ := result.LastInsertId()
		i.ID = uint(id)

		if len(i.Teams) == 0 {
			i.Teams = []fleet.UserTeam{}
			return nil
		}

		// Bulk insert teams
		const valueStr = "(?,?,?),"
		var args []interface{}
		for _, userTeam := range i.Teams {
			args = append(args, i.ID, userTeam.Team.ID, userTeam.Role)
		}
		sql := "INSERT INTO invite_teams (invite_id, team_id, role) VALUES " +
			strings.Repeat(valueStr, len(i.Teams))
		sql = strings.TrimSuffix(sql, ",")
		if _, err := tx.ExecContext(ctx, sql, args...); err != nil {
			return errors.Wrap(err, "insert teams")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return i, nil
}

// ListInvites lists all invites in the Fleet database. Supply query options
// using the opt parameter. See fleet.ListOptions
func (d *Datastore) ListInvites(ctx context.Context, opt fleet.ListOptions) ([]*fleet.Invite, error) {
	invites := []*fleet.Invite{}
	query := "SELECT * FROM invites WHERE true"
	query, params := searchLike(query, nil, opt.MatchQuery, inviteSearchColumns...)
	query = appendListOptionsToSQL(query, opt)

	err := sqlx.SelectContext(ctx, d.reader, &invites, query, params...)
	if err == sql.ErrNoRows {
		return nil, notFound("Invite")
	} else if err != nil {
		return nil, errors.Wrap(err, "select invite by ID")
	}

	if err := d.loadTeamsForInvites(ctx, invites); err != nil {
		return nil, errors.Wrap(err, "load teams")
	}

	return invites, nil
}

// Invite returns Invite identified by id.
func (d *Datastore) Invite(ctx context.Context, id uint) (*fleet.Invite, error) {
	var invite fleet.Invite
	err := sqlx.GetContext(ctx, d.reader, &invite, "SELECT * FROM invites WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return nil, notFound("Invite").WithID(id)
	} else if err != nil {
		return nil, errors.Wrap(err, "select invite by ID")
	}

	if err := d.loadTeamsForInvites(ctx, []*fleet.Invite{&invite}); err != nil {
		return nil, errors.Wrap(err, "load teams")
	}

	return &invite, nil
}

// InviteByEmail finds an Invite with a particular email, if one exists.
func (d *Datastore) InviteByEmail(ctx context.Context, email string) (*fleet.Invite, error) {
	var invite fleet.Invite
	err := sqlx.GetContext(ctx, d.reader, &invite, "SELECT * FROM invites WHERE email = ?", email)
	if err == sql.ErrNoRows {
		return nil, notFound("Invite").
			WithMessage(fmt.Sprintf("with email %s", email))
	} else if err != nil {
		return nil, errors.Wrap(err, "sqlx get invite by email")
	}

	if err := d.loadTeamsForInvites(ctx, []*fleet.Invite{&invite}); err != nil {
		return nil, errors.Wrap(err, "load teams")
	}

	return &invite, nil
}

// InviteByToken finds an Invite with a particular token, if one exists.
func (d *Datastore) InviteByToken(ctx context.Context, token string) (*fleet.Invite, error) {
	var invite fleet.Invite
	err := sqlx.GetContext(ctx, d.reader, &invite, "SELECT * FROM invites WHERE token = ?", token)
	if err == sql.ErrNoRows {
		return nil, notFound("Invite").
			WithMessage(fmt.Sprintf("with token %s", token))
	} else if err != nil {
		return nil, errors.Wrap(err, "sqlx get invite by token")
	}

	if err := d.loadTeamsForInvites(ctx, []*fleet.Invite{&invite}); err != nil {
		return nil, errors.Wrap(err, "load teams")
	}

	return &invite, nil
}

func (d *Datastore) DeleteInvite(ctx context.Context, id uint) error {
	return d.deleteEntity(ctx, invitesTable, id)
}

func (d *Datastore) loadTeamsForInvites(ctx context.Context, invites []*fleet.Invite) error {
	inviteIDs := make([]uint, 0, len(invites)+1)
	// Make sure the slice is never empty for IN by filling a nonexistent ID
	inviteIDs = append(inviteIDs, 0)
	idToInvite := make(map[uint]*fleet.Invite, len(invites))
	for _, u := range invites {
		// Initialize empty slice so we get an array in JSON responses instead
		// of null if it is empty
		u.Teams = []fleet.UserTeam{}
		// Track IDs for queries and matching
		inviteIDs = append(inviteIDs, u.ID)
		idToInvite[u.ID] = u
	}

	sql := `
		SELECT ut.team_id AS id, ut.invite_id, ut.role, t.name
		FROM invite_teams ut INNER JOIN teams t ON ut.team_id = t.id
		WHERE ut.invite_id IN (?)
		ORDER BY invite_id, team_id
	`
	sql, args, err := sqlx.In(sql, inviteIDs)
	if err != nil {
		return errors.Wrap(err, "sqlx.In loadTeamsForInvites")
	}

	var rows []struct {
		fleet.UserTeam
		InviteID uint `db:"invite_id"`
	}
	if err := sqlx.SelectContext(ctx, d.reader, &rows, sql, args...); err != nil {
		return errors.Wrap(err, "get loadTeamsForInvites")
	}

	// Map each row to the appropriate invite
	for _, r := range rows {
		invite := idToInvite[r.InviteID]
		invite.Teams = append(invite.Teams, r.UserTeam)
	}

	return nil
}
