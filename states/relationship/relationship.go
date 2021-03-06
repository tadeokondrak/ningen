package relationship

import (
	"sync"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/diamondburned/ningen/handlerrepo"
)

type State struct {
	mutex         sync.RWMutex
	relationships map[discord.Snowflake]discord.RelationshipType
}

func NewState(r handlerrepo.AddHandler) *State {
	relastate := &State{
		relationships: map[discord.Snowflake]discord.RelationshipType{},
	}

	r.AddHandler(func(r *gateway.ReadyEvent) {
		relastate.mutex.Lock()
		defer relastate.mutex.Unlock()

		for _, rl := range r.Relationships {
			relastate.relationships[rl.UserID] = rl.Type
		}
	})

	r.AddHandler(func(add *gateway.RelationshipAddEvent) {
		relastate.mutex.Lock()
		defer relastate.mutex.Unlock()

		relastate.relationships[add.UserID] = add.Type
	})

	r.AddHandler(func(rem *gateway.RelationshipRemoveEvent) {
		relastate.mutex.Lock()
		defer relastate.mutex.Unlock()

		delete(relastate.relationships, rem.UserID)
	})

	return relastate
}

// Relationship returns the relationship for the given user, or 0 if there is
// none.
func (r *State) Relationship(userID discord.Snowflake) discord.RelationshipType {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if t, ok := r.relationships[userID]; ok {
		return t
	}

	return 0
}

// Blocked returns if the user is blocked.
func (r *State) Blocked(userID discord.Snowflake) bool {
	return r.Relationship(userID) == discord.BlockedRelationship
}
