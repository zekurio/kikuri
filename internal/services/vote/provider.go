package vote

import (
	"github.com/zekrotja/ken"
	"github.com/zekurio/kikuri/internal/models"
)

// VotesProvider is an interface for the vote service
// it handles the storage of votes, syncing with
// the redis cache and database and provides functions
// to create, update and delete votes
// it also handles the vote expiration
type VotesProvider interface {
	// Populate populates the votes map with the data from the database
	Populate() error

	// Create creates a new vote and adds it to the votes map
	Create(ctx ken.SubCommandContext, vote models.Vote) error

	// Get returns a vote by its id
	Get(id string) (models.Vote, error)

	// GetAllFromGuild returns all votes from a guild
	GetAllFromGuild(guildID string) map[string]models.Vote

	// Update updates a vote
	Update(vote models.Vote) error

	// Close closes a vote, this can also be called by the expiration
	// handler but needs to pass the state of the vote
	Close(id string, state models.VoteState) error

	// CloseAll closes all votes from a guild, returns amount of closed votes
	CloseAll(guildID string) (int, error)
}
