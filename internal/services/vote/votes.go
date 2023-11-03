package vote

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/ken"
	"github.com/zekurio/kikuri/internal/models"
	"github.com/zekurio/kikuri/internal/services/database"
	"github.com/zekurio/kikuri/internal/services/scheduler"
	"github.com/zekurio/kikuri/internal/util/static"
	"time"
)

// VotesHandler is the implementation of the votes service provider
type VotesHandler struct {
	db    database.Database
	sched scheduler.Provider
	ken   ken.IKen

	runningVotes map[string]models.Vote
}

// NewVotesHandler returns a new votes handler instance
func NewVotesHandler(ctn di.Container) *VotesHandler {
	return &VotesHandler{
		db:           ctn.Get(static.DiDatabase).(database.Database),
		sched:        ctn.Get(static.DiScheduler).(scheduler.Provider),
		ken:          ctn.Get(static.DiCommandHandler).(ken.IKen),
		runningVotes: make(map[string]models.Vote),
	}
}

// Populate populates the votes map with the data from the database
func (v *VotesHandler) Populate() error {
	// get all votes from database
	votes, err := v.db.GetVotes()
	if err != nil {
		return err
	}

	// iterate over all votes
	for _, vote := range votes {
		// add vote to map
		v.runningVotes[vote.ID] = vote

		// schedule closing
		// TODO check if vote is set to expire
		if (vote.Expires != time.Time{}) {
			_, err = v.sched.Schedule(scheduler.FormatCronJobSpec(vote.Expires, time.Second*5), func() {
				now := time.Now()
				if (vote.Expires != time.Time{} && vote.Expires.Before(now)) {
					err := v.Close(vote.ID, models.StateExpired)
					if err != nil {
						log.Error("Failed closing vote", err)
					}
				}
			})
			if err != nil {
				return err
			}
		}

		b := v.ken.Components().Add(vote.MessageID, vote.ChannelID)

		_, err = vote.AddButtons(b)
		if err != nil {
			log.Error("Failed adding buttons to vote")
		}
	}

	return nil
}

// Create handles the vote creation
func (v *VotesHandler) Create(ctx ken.SubCommandContext, vote models.Vote) error {

	// handle followup message
	emb, err := vote.AsEmbed(ctx.GetSession())
	if err != nil {
		return err
	}

	fum := ctx.FollowUpEmbed(emb).Send()
	err = fum.Error
	if err != nil {
		return err
	}

	b := fum.AddComponents()

	vote.MessageID = fum.Message.ID

	_, err = vote.AddButtons(b)
	if err != nil {
		return err
	}

	// add vote to map
	v.runningVotes[vote.ID] = vote

	// save
	err = v.db.AddUpdateVote(vote)

	// schedule closing
	if (vote.Expires != time.Time{}) {
		_, err = v.sched.Schedule(scheduler.FormatCronJobSpec(vote.Expires, time.Second*5), func() {
			now := time.Now()
			if (vote.Expires != time.Time{} && vote.Expires.Before(now)) {
				err := v.Close(vote.ID, models.StateExpired)
				if err != nil {
					log.Error("Failed closing vote", err)
				}
			}
		})
	}

	return err
}

// Get returns a vote by its id
func (v *VotesHandler) Get(id string) (models.Vote, error) {
	if vote, ok := v.runningVotes[id]; ok {
		return vote, nil
	}

	return models.Vote{}, errors.New("vote not found")
}

// GetAllFromGuild returns all votes from a guild
func (v *VotesHandler) GetAllFromGuild(guildID string) map[string]models.Vote {
	var votes = make(map[string]models.Vote)
	for _, vote := range v.runningVotes {
		if vote.GuildID == guildID {
			votes[vote.ID] = vote
		}
	}

	return votes
}

// Close handles the closing and deletion of a vote
func (v *VotesHandler) Close(id string, state models.VoteState) error {
	vote, err := v.Get(id)

	// execute the close function of the vote itself
	err = vote.Close(v.ken, state)
	if err != nil {
		return err
	}

	// delete from map
	delete(v.runningVotes, vote.ID)

	// delete from database
	err = v.db.DeleteVote(vote.ID)

	return err
}

// CloseAll closes all votes from a guild
func (v *VotesHandler) CloseAll(guildID string) (int, error) {
	// iterate over all votes
	var votesClosed int
	for _, vote := range v.runningVotes {
		// check if vote is from guild
		if vote.GuildID == guildID {
			// close vote
			err := v.Close(vote.ID, models.StateClosed)
			if err != nil {
				return votesClosed, err
			}
			votesClosed++
		}
	}

	return votesClosed, nil
}

// Update updates a vote
func (v *VotesHandler) Update(vote models.Vote) error {
	// update vote in map
	v.runningVotes[vote.ID] = vote

	// update vote in database
	return v.db.AddUpdateVote(vote)
}
