package votes

import "github.com/zekurio/daemon/internal/models"

// Provider is an interface for the vote service
// it handles the storage of votes, syncing with
// the redis cache and database and provides functions
// to create, update and delete votes
// it also handles the vote expiration
type Provider interface {
	// Populate populates the votes map with the data from the database
	Populate() error

	// Create creates a new vote and adds it to the votes map
	// TODO decide if this should use a vote or take in the data as parameters
	Create(vote models.Vote) error

	// Get returns a vote by its id
	Get(id string) (models.Vote, error)

	// Update updates a vote
	Update(vote models.Vote) error

	// Close closes a vote, this can also be called by the expiration
	// handler but needs to pass the state of the vote
	Close(vote models.Vote, state models.VoteState) error
}
