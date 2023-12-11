package roleutils

import (
	"errors"
	"github.com/zekrotja/dgrs"

	"github.com/bwmarrin/discordgo"
)

// GetRoleByID returns the role with the given ID
func GetRoleByID(state *dgrs.State, guildID, roleID string) (*discordgo.Role, error) {
	roles, err := state.Roles(guildID, true)
	if err != nil {
		return nil, err
	}

	for _, r := range roles {
		if r.ID == roleID {
			return r, nil
		}
	}

	return nil, errors.New("role not found")
}

// Sort sorts the roles either ascending or descending
func Sort(roles []*discordgo.Role, reversed bool) []*discordgo.Role {
	if reversed {
		for i := 0; i < len(roles)/2; i++ {
			roles[i], roles[len(roles)-1-i] = roles[len(roles)-1-i], roles[i]
		}
	} else {
		for i := 0; i < len(roles)/2; i++ {
			roles[i], roles[len(roles)-1-i] = roles[len(roles)-1-i], roles[i]
		}
	}

	return roles
}

// GetSortedMemberRoles returns the guilds roles sorted either ascending or
// descending. Can also include the @everyone role.
func GetSortedMemberRoles(state *dgrs.State, guildID, memberID string, includeEveryone, reversed bool) ([]*discordgo.Role, error) {
	member, err := state.Member(guildID, memberID)
	if err != nil {
		return nil, err
	}

	roles, err := state.Roles(guildID, true)
	if err != nil {
		return nil, err
	}

	rolesMap := make(map[string]*discordgo.Role)
	for _, r := range roles {
		rolesMap[r.ID] = r
	}

	mRoles := make([]*discordgo.Role, len(member.Roles)+1)
	applied := 0
	for _, rID := range member.Roles {
		if r, ok := rolesMap[rID]; ok {
			mRoles[applied] = r
			applied++
		}
	}

	if includeEveryone {
		mRoles[applied] = rolesMap[guildID]
		applied++
	}

	mRoles = mRoles[:applied]

	return Sort(mRoles, reversed), nil
}

// GetSortedGuildRoles returns the guilds roles sorted either ascending or
// descending.
func GetSortedGuildRoles(state *dgrs.State, guildID string, reversed bool) ([]*discordgo.Role, error) {
	roles, err := state.Roles(guildID, true)
	if err != nil {
		return nil, err
	}

	return Sort(roles, reversed), nil
}
