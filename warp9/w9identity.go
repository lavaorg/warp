// Copyright 2018 Larry Rau. All rights reserved
// See Apache2 LICENSE

package warp9

import (
	"strconv"
	"sync"
)

// provide an interface to an external source of user identity.
// The warp9 protocol will only concern itself with the numeric ID of the user/group
// The implementation may request a symbolic name of a user-id in the form of a string
// only for debugging outputs.
// user-id will not be emitted to any logs

// a simple file based default implementation will be used if external interfaces
// are not provided.

// Interface for accessing users and groups
type Users interface {
	User(uid uint32) User
	Group(gid uint32) Group
}

// Represents a user
type User interface {
	Name() string
	Id() uint32            // user id
	Groups() []Group       // groups the user belongs to (can return nil)
	IsMember(g Group) bool // returns true if the user is member of the specified group
}

// Represents a group of users
type Group interface {
	Name() string
	Id() uint32      // group id
	Members() []User // list of members that belong to the group (can return nil)
}

// default, simple implementation

var once sync.Once

type w9identity struct {
	users  map[uint32]*w9user
	groups map[uint32]*w9group
	sync.Mutex
}

type w9user struct {
	uid  uint32
	gid  uint32
	uids string
	note string
}

type w9group struct {
	gid  uint32
	gids string
	note string
}

// Simple Users implementation that fakes looking up users and groups
// by uid only. The names and groups memberships are empty
var Identity *w9identity

func (u *w9user) Name() string { return u.uids }

func (u *w9user) Id() uint32 { return u.uid }

func (u *w9user) Groups() []Group {
	return []Group{Identity.groups[u.uid]}
}

func (u *w9user) IsMember(g Group) bool { return u.gid == g.Id() }

func (g *w9group) Name() string { return g.gids }

func (g *w9group) Id() uint32 { return g.gid }

func (g *w9group) Members() []User {
	m := []User{}
	for _, u := range Identity.users {
		if u.gid == g.gid {
			var newu w9user = *u //make copy
			m = append(m, &newu)
		}
	}
	return m
}

func init() {
	Identity = new(w9identity)
	Identity.users = make(map[uint32]*w9user)
	Identity.groups = make(map[uint32]*w9group)

	Identity.addGroup(&w9group{1, "sys", "system users"})
	Identity.addGroup(&w9group{2, "nobody", "no users"})
	Identity.addGroup(&w9group{20, "staff", "normal users"})

	Identity.addUser(&w9user{1, 1, "sys", "system user"})
	Identity.addUser(&w9user{501, 20, "larry", "a good citizen"})

}

func (id *w9identity) User(uid uint32) User {
	//once.Do(initIdentity)
	Identity.Lock()
	defer Identity.Unlock()
	user, present := Identity.users[uid]
	if present {
		return user
	}

	user = new(w9user)
	user.uid = uid
	Identity.users[uid] = user
	return user
}

func (id *w9identity) Group(gid uint32) Group {
	//once.Do(initIdentity)
	Identity.Lock()
	defer Identity.Unlock()
	group, present := Identity.groups[gid]
	if present {
		return group
	}

	group = new(w9group)
	group.gid = gid
	Identity.groups[gid] = group
	return group
}

func (id *w9identity) DeprecatedUnameToUser(uname string) User {
	u, e := strconv.ParseUint(uname, 10, 32)
	if e != nil {
		return nil
	}
	return id.User(uint32(u))
}

func (id *w9identity) addUser(user *w9user) {
	id.users[user.uid] = user
}

func (id *w9identity) addGroup(grp *w9group) {
	id.groups[grp.gid] = grp
}
