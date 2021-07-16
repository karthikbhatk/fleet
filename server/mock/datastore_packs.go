// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import "github.com/fleetdm/fleet/v4/server/fleet"

var _ fleet.PackStore = (*PackStore)(nil)

type ApplyPackSpecsFunc func(specs []*fleet.PackSpec) error

type GetPackSpecsFunc func() ([]*fleet.PackSpec, error)

type GetPackSpecFunc func(name string) (*fleet.PackSpec, error)

type NewPackFunc func(pack *fleet.Pack, opts ...fleet.OptionalArg) (*fleet.Pack, error)

type SavePackFunc func(pack *fleet.Pack) error

type DeletePackFunc func(name string) error

type PackFunc func(pid uint) (*fleet.Pack, error)

type ListPacksFunc func(opt fleet.ListOptions) ([]*fleet.Pack, error)

type PackByNameFunc func(name string, opts ...fleet.OptionalArg) (*fleet.Pack, bool, error)

type AddLabelToPackFunc func(lid uint, pid uint, opts ...fleet.OptionalArg) error

type RemoveLabelFromPackFunc func(lid uint, pid uint) error

type ListLabelsForPackFunc func(pid uint) ([]*fleet.Label, error)

type AddHostToPackFunc func(hid uint, pid uint) error

type RemoveHostFromPackFunc func(hid uint, pid uint) error

type ListPacksForHostFunc func(hid uint) (packs []*fleet.Pack, err error)

type ListHostsInPackFunc func(pid uint, opt fleet.ListOptions) ([]uint, error)

type ListExplicitHostsInPackFunc func(pid uint, opt fleet.ListOptions) ([]uint, error)

type EnsureGlobalPackFunc func() (*fleet.Pack, error)

type PackStore struct {
	ApplyPackSpecsFunc        ApplyPackSpecsFunc
	ApplyPackSpecsFuncInvoked bool

	GetPackSpecsFunc        GetPackSpecsFunc
	GetPackSpecsFuncInvoked bool

	GetPackSpecFunc        GetPackSpecFunc
	GetPackSpecFuncInvoked bool

	NewPackFunc        NewPackFunc
	NewPackFuncInvoked bool

	SavePackFunc        SavePackFunc
	SavePackFuncInvoked bool

	DeletePackFunc        DeletePackFunc
	DeletePackFuncInvoked bool

	PackFunc        PackFunc
	PackFuncInvoked bool

	ListPacksFunc        ListPacksFunc
	ListPacksFuncInvoked bool

	PackByNameFunc        PackByNameFunc
	PackByNameFuncInvoked bool

	AddLabelToPackFunc        AddLabelToPackFunc
	AddLabelToPackFuncInvoked bool

	RemoveLabelFromPackFunc        RemoveLabelFromPackFunc
	RemoveLabelFromPackFuncInvoked bool

	ListLabelsForPackFunc        ListLabelsForPackFunc
	ListLabelsForPackFuncInvoked bool

	AddHostToPackFunc        AddHostToPackFunc
	AddHostToPackFuncInvoked bool

	RemoveHostFromPackFunc        RemoveHostFromPackFunc
	RemoveHostFromPackFuncInvoked bool

	ListPacksForHostFunc        ListPacksForHostFunc
	ListPacksForHostFuncInvoked bool

	ListHostsInPackFunc        ListHostsInPackFunc
	ListHostsInPackFuncInvoked bool

	ListExplicitHostsInPackFunc        ListExplicitHostsInPackFunc
	ListExplicitHostsInPackFuncInvoked bool

	EnsureGlobalPackFunc        EnsureGlobalPackFunc
	EnsureGlobalPackFuncInvoked bool
}

func (s *PackStore) EnsureGlobalPack() (*fleet.Pack, error) {
	panic("implement me")
}

func (s *PackStore) ApplyPackSpecs(specs []*fleet.PackSpec) error {
	s.ApplyPackSpecsFuncInvoked = true
	return s.ApplyPackSpecsFunc(specs)
}

func (s *PackStore) GetPackSpecs() ([]*fleet.PackSpec, error) {
	s.GetPackSpecsFuncInvoked = true
	return s.GetPackSpecsFunc()
}

func (s *PackStore) GetPackSpec(name string) (*fleet.PackSpec, error) {
	s.GetPackSpecFuncInvoked = true
	return s.GetPackSpecFunc(name)
}

func (s *PackStore) NewPack(pack *fleet.Pack, opts ...fleet.OptionalArg) (*fleet.Pack, error) {
	s.NewPackFuncInvoked = true
	return s.NewPackFunc(pack, opts...)
}

func (s *PackStore) SavePack(pack *fleet.Pack) error {
	s.SavePackFuncInvoked = true
	return s.SavePackFunc(pack)
}

func (s *PackStore) DeletePack(name string) error {
	s.DeletePackFuncInvoked = true
	return s.DeletePackFunc(name)
}

func (s *PackStore) Pack(pid uint) (*fleet.Pack, error) {
	s.PackFuncInvoked = true
	return s.PackFunc(pid)
}

func (s *PackStore) ListPacks(opt fleet.ListOptions) ([]*fleet.Pack, error) {
	s.ListPacksFuncInvoked = true
	return s.ListPacksFunc(opt)
}

func (s *PackStore) PackByName(name string, opts ...fleet.OptionalArg) (*fleet.Pack, bool, error) {
	s.PackByNameFuncInvoked = true
	return s.PackByNameFunc(name, opts...)
}

func (s *PackStore) AddLabelToPack(lid uint, pid uint, opts ...fleet.OptionalArg) error {
	s.AddLabelToPackFuncInvoked = true
	return s.AddLabelToPackFunc(lid, pid, opts...)
}

func (s *PackStore) RemoveLabelFromPack(lid uint, pid uint) error {
	s.RemoveLabelFromPackFuncInvoked = true
	return s.RemoveLabelFromPackFunc(lid, pid)
}

func (s *PackStore) ListLabelsForPack(pid uint) ([]*fleet.Label, error) {
	s.ListLabelsForPackFuncInvoked = true
	return s.ListLabelsForPackFunc(pid)
}

func (s *PackStore) AddHostToPack(hid uint, pid uint) error {
	s.AddHostToPackFuncInvoked = true
	return s.AddHostToPackFunc(hid, pid)
}

func (s *PackStore) RemoveHostFromPack(hid uint, pid uint) error {
	s.RemoveHostFromPackFuncInvoked = true
	return s.RemoveHostFromPackFunc(hid, pid)
}

func (s *PackStore) ListPacksForHost(hid uint) (packs []*fleet.Pack, err error) {
	s.ListPacksForHostFuncInvoked = true
	return s.ListPacksForHostFunc(hid)
}

func (s *PackStore) ListHostsInPack(pid uint, opt fleet.ListOptions) ([]uint, error) {
	s.ListHostsInPackFuncInvoked = true
	return s.ListHostsInPackFunc(pid, opt)
}

func (s *PackStore) ListExplicitHostsInPack(pid uint, opt fleet.ListOptions) ([]uint, error) {
	s.ListExplicitHostsInPackFuncInvoked = true
	return s.ListExplicitHostsInPackFunc(pid, opt)
}
