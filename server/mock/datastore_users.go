// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import "github.com/fleetdm/fleet/v4/server/fleet"

var _ fleet.UserStore = (*UserStore)(nil)

type NewUserFunc func(user *fleet.User) (*fleet.User, error)

type ListUsersFunc func(opt fleet.UserListOptions) ([]*fleet.User, error)

type UserByEmailFunc func(email string) (*fleet.User, error)

type UserByIDFunc func(id uint) (*fleet.User, error)

type SaveUserFunc func(user *fleet.User) error

type SaveUsersFunc func(users []*fleet.User) error

type DeleteUserFunc func(id uint) error

type PendingEmailChangeFunc func(userID uint, newEmail string, token string) error

type ConfirmPendingEmailChangeFunc func(userID uint, token string) (string, error)

type UserStore struct {
	NewUserFunc        NewUserFunc
	NewUserFuncInvoked bool

	ListUsersFunc        ListUsersFunc
	ListUsersFuncInvoked bool

	UserByEmailFunc        UserByEmailFunc
	UserByEmailFuncInvoked bool

	UserByIDFunc        UserByIDFunc
	UserByIDFuncInvoked bool

	SaveUserFunc        SaveUserFunc
	SaveUserFuncInvoked bool

	SaveUsersFunc        SaveUsersFunc
	SaveUsersFuncInvoked bool

	DeleteUserFunc        DeleteUserFunc
	DeleteUserFuncInvoked bool

	PendingEmailChangeFunc        PendingEmailChangeFunc
	PendingEmailChangeFuncInvoked bool

	ConfirmPendingEmailChangeFunc        ConfirmPendingEmailChangeFunc
	ConfirmPendingEmailChangeFuncInvoked bool
}

func (s *UserStore) NewUser(user *fleet.User) (*fleet.User, error) {
	s.NewUserFuncInvoked = true
	return s.NewUserFunc(user)
}

func (s *UserStore) ListUsers(opt fleet.UserListOptions) ([]*fleet.User, error) {
	s.ListUsersFuncInvoked = true
	return s.ListUsersFunc(opt)
}

func (s *UserStore) UserByEmail(email string) (*fleet.User, error) {
	s.UserByEmailFuncInvoked = true
	return s.UserByEmailFunc(email)
}

func (s *UserStore) UserByID(id uint) (*fleet.User, error) {
	s.UserByIDFuncInvoked = true
	return s.UserByIDFunc(id)
}

func (s *UserStore) SaveUser(user *fleet.User) error {
	s.SaveUserFuncInvoked = true
	return s.SaveUserFunc(user)
}

func (s *UserStore) SaveUsers(users []*fleet.User) error {
	s.SaveUsersFuncInvoked = true
	return s.SaveUsersFunc(users)
}

func (s *UserStore) DeleteUser(id uint) error {
	s.DeleteUserFuncInvoked = true
	return s.DeleteUserFunc(id)
}

func (s *UserStore) PendingEmailChange(userID uint, newEmail string, token string) error {
	s.PendingEmailChangeFuncInvoked = true
	return s.PendingEmailChangeFunc(userID, newEmail, token)
}

func (s *UserStore) ConfirmPendingEmailChange(userID uint, token string) (string, error) {
	s.ConfirmPendingEmailChangeFuncInvoked = true
	return s.ConfirmPendingEmailChangeFunc(userID, token)
}
