// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/mdm/nanomdm/mdm"
)

var _ fleet.MDMAppleStore = (*MDMAppleStore)(nil)

type StoreAuthenticateFunc func(r *mdm.Request, msg *mdm.Authenticate) error

type StoreTokenUpdateFunc func(r *mdm.Request, msg *mdm.TokenUpdate) error

type StoreUserAuthenticateFunc func(r *mdm.Request, msg *mdm.UserAuthenticate) error

type DisableFunc func(r *mdm.Request) error

type StoreCommandReportFunc func(r *mdm.Request, report *mdm.CommandResults) error

type RetrieveNextCommandFunc func(r *mdm.Request, skipNotNow bool) (*mdm.Command, error)

type ClearQueueFunc func(r *mdm.Request) error

type StoreBootstrapTokenFunc func(r *mdm.Request, msg *mdm.SetBootstrapToken) error

type RetrieveBootstrapTokenFunc func(r *mdm.Request, msg *mdm.GetBootstrapToken) (*mdm.BootstrapToken, error)

type RetrievePushInfoFunc func(p0 context.Context, p1 []string) (map[string]*mdm.Push, error)

type IsPushCertStaleFunc func(ctx context.Context, topic string, staleToken string) (bool, error)

type RetrievePushCertFunc func(ctx context.Context, topic string) (cert *tls.Certificate, staleToken string, err error)

type StorePushCertFunc func(ctx context.Context, pemCert []byte, pemKey []byte) error

type EnqueueCommandFunc func(ctx context.Context, id []string, cmd *mdm.Command) (map[string]error, error)

type HasCertHashFunc func(r *mdm.Request, hash string) (bool, error)

type EnrollmentHasCertHashFunc func(r *mdm.Request, hash string) (bool, error)

type IsCertHashAssociatedFunc func(r *mdm.Request, hash string) (bool, error)

type AssociateCertHashFunc func(r *mdm.Request, hash string, certNotValidAfter time.Time) error

type RetrieveMigrationCheckinsFunc func(p0 context.Context, p1 chan<- interface{}) error

type RetrieveTokenUpdateTallyFunc func(ctx context.Context, id string) (int, error)

type GetAllMDMConfigAssetsByNameFunc func(ctx context.Context, assetNames []fleet.MDMAssetName) (map[fleet.MDMAssetName]fleet.MDMConfigAsset, error)

type GetABMTokenByOrgNameFunc func(ctx context.Context, orgName string) (*fleet.ABMToken, error)

type EnqueueDeviceLockCommandFunc func(ctx context.Context, host *fleet.Host, cmd *mdm.Command, pin string) error

type EnqueueDeviceWipeCommandFunc func(ctx context.Context, host *fleet.Host, cmd *mdm.Command) error

type MDMAppleStore struct {
	StoreAuthenticateFunc        StoreAuthenticateFunc
	StoreAuthenticateFuncInvoked bool

	StoreTokenUpdateFunc        StoreTokenUpdateFunc
	StoreTokenUpdateFuncInvoked bool

	StoreUserAuthenticateFunc        StoreUserAuthenticateFunc
	StoreUserAuthenticateFuncInvoked bool

	DisableFunc        DisableFunc
	DisableFuncInvoked bool

	StoreCommandReportFunc        StoreCommandReportFunc
	StoreCommandReportFuncInvoked bool

	RetrieveNextCommandFunc        RetrieveNextCommandFunc
	RetrieveNextCommandFuncInvoked bool

	ClearQueueFunc        ClearQueueFunc
	ClearQueueFuncInvoked bool

	StoreBootstrapTokenFunc        StoreBootstrapTokenFunc
	StoreBootstrapTokenFuncInvoked bool

	RetrieveBootstrapTokenFunc        RetrieveBootstrapTokenFunc
	RetrieveBootstrapTokenFuncInvoked bool

	RetrievePushInfoFunc        RetrievePushInfoFunc
	RetrievePushInfoFuncInvoked bool

	IsPushCertStaleFunc        IsPushCertStaleFunc
	IsPushCertStaleFuncInvoked bool

	RetrievePushCertFunc        RetrievePushCertFunc
	RetrievePushCertFuncInvoked bool

	StorePushCertFunc        StorePushCertFunc
	StorePushCertFuncInvoked bool

	EnqueueCommandFunc        EnqueueCommandFunc
	EnqueueCommandFuncInvoked bool

	HasCertHashFunc        HasCertHashFunc
	HasCertHashFuncInvoked bool

	EnrollmentHasCertHashFunc        EnrollmentHasCertHashFunc
	EnrollmentHasCertHashFuncInvoked bool

	IsCertHashAssociatedFunc        IsCertHashAssociatedFunc
	IsCertHashAssociatedFuncInvoked bool

	AssociateCertHashFunc        AssociateCertHashFunc
	AssociateCertHashFuncInvoked bool

	RetrieveMigrationCheckinsFunc        RetrieveMigrationCheckinsFunc
	RetrieveMigrationCheckinsFuncInvoked bool

	RetrieveTokenUpdateTallyFunc        RetrieveTokenUpdateTallyFunc
	RetrieveTokenUpdateTallyFuncInvoked bool

	GetAllMDMConfigAssetsByNameFunc        GetAllMDMConfigAssetsByNameFunc
	GetAllMDMConfigAssetsByNameFuncInvoked bool

	GetABMTokenByOrgNameFunc        GetABMTokenByOrgNameFunc
	GetABMTokenByOrgNameFuncInvoked bool

	EnqueueDeviceLockCommandFunc        EnqueueDeviceLockCommandFunc
	EnqueueDeviceLockCommandFuncInvoked bool

	EnqueueDeviceWipeCommandFunc        EnqueueDeviceWipeCommandFunc
	EnqueueDeviceWipeCommandFuncInvoked bool

	mu sync.Mutex
}

func (fs *MDMAppleStore) StoreAuthenticate(r *mdm.Request, msg *mdm.Authenticate) error {
	fs.mu.Lock()
	fs.StoreAuthenticateFuncInvoked = true
	fs.mu.Unlock()
	return fs.StoreAuthenticateFunc(r, msg)
}

func (fs *MDMAppleStore) StoreTokenUpdate(r *mdm.Request, msg *mdm.TokenUpdate) error {
	fs.mu.Lock()
	fs.StoreTokenUpdateFuncInvoked = true
	fs.mu.Unlock()
	return fs.StoreTokenUpdateFunc(r, msg)
}

func (fs *MDMAppleStore) StoreUserAuthenticate(r *mdm.Request, msg *mdm.UserAuthenticate) error {
	fs.mu.Lock()
	fs.StoreUserAuthenticateFuncInvoked = true
	fs.mu.Unlock()
	return fs.StoreUserAuthenticateFunc(r, msg)
}

func (fs *MDMAppleStore) Disable(r *mdm.Request) error {
	fs.mu.Lock()
	fs.DisableFuncInvoked = true
	fs.mu.Unlock()
	return fs.DisableFunc(r)
}

func (fs *MDMAppleStore) StoreCommandReport(r *mdm.Request, report *mdm.CommandResults) error {
	fs.mu.Lock()
	fs.StoreCommandReportFuncInvoked = true
	fs.mu.Unlock()
	return fs.StoreCommandReportFunc(r, report)
}

func (fs *MDMAppleStore) RetrieveNextCommand(r *mdm.Request, skipNotNow bool) (*mdm.Command, error) {
	fs.mu.Lock()
	fs.RetrieveNextCommandFuncInvoked = true
	fs.mu.Unlock()
	return fs.RetrieveNextCommandFunc(r, skipNotNow)
}

func (fs *MDMAppleStore) ClearQueue(r *mdm.Request) error {
	fs.mu.Lock()
	fs.ClearQueueFuncInvoked = true
	fs.mu.Unlock()
	return fs.ClearQueueFunc(r)
}

func (fs *MDMAppleStore) StoreBootstrapToken(r *mdm.Request, msg *mdm.SetBootstrapToken) error {
	fs.mu.Lock()
	fs.StoreBootstrapTokenFuncInvoked = true
	fs.mu.Unlock()
	return fs.StoreBootstrapTokenFunc(r, msg)
}

func (fs *MDMAppleStore) RetrieveBootstrapToken(r *mdm.Request, msg *mdm.GetBootstrapToken) (*mdm.BootstrapToken, error) {
	fs.mu.Lock()
	fs.RetrieveBootstrapTokenFuncInvoked = true
	fs.mu.Unlock()
	return fs.RetrieveBootstrapTokenFunc(r, msg)
}

func (fs *MDMAppleStore) RetrievePushInfo(p0 context.Context, p1 []string) (map[string]*mdm.Push, error) {
	fs.mu.Lock()
	fs.RetrievePushInfoFuncInvoked = true
	fs.mu.Unlock()
	return fs.RetrievePushInfoFunc(p0, p1)
}

func (fs *MDMAppleStore) IsPushCertStale(ctx context.Context, topic string, staleToken string) (bool, error) {
	fs.mu.Lock()
	fs.IsPushCertStaleFuncInvoked = true
	fs.mu.Unlock()
	return fs.IsPushCertStaleFunc(ctx, topic, staleToken)
}

func (fs *MDMAppleStore) RetrievePushCert(ctx context.Context, topic string) (cert *tls.Certificate, staleToken string, err error) {
	fs.mu.Lock()
	fs.RetrievePushCertFuncInvoked = true
	fs.mu.Unlock()
	return fs.RetrievePushCertFunc(ctx, topic)
}

func (fs *MDMAppleStore) StorePushCert(ctx context.Context, pemCert []byte, pemKey []byte) error {
	fs.mu.Lock()
	fs.StorePushCertFuncInvoked = true
	fs.mu.Unlock()
	return fs.StorePushCertFunc(ctx, pemCert, pemKey)
}

func (fs *MDMAppleStore) EnqueueCommand(ctx context.Context, id []string, cmd *mdm.Command) (map[string]error, error) {
	fs.mu.Lock()
	fs.EnqueueCommandFuncInvoked = true
	fs.mu.Unlock()
	return fs.EnqueueCommandFunc(ctx, id, cmd)
}

func (fs *MDMAppleStore) HasCertHash(r *mdm.Request, hash string) (bool, error) {
	fs.mu.Lock()
	fs.HasCertHashFuncInvoked = true
	fs.mu.Unlock()
	return fs.HasCertHashFunc(r, hash)
}

func (fs *MDMAppleStore) EnrollmentHasCertHash(r *mdm.Request, hash string) (bool, error) {
	fs.mu.Lock()
	fs.EnrollmentHasCertHashFuncInvoked = true
	fs.mu.Unlock()
	return fs.EnrollmentHasCertHashFunc(r, hash)
}

func (fs *MDMAppleStore) IsCertHashAssociated(r *mdm.Request, hash string) (bool, error) {
	fs.mu.Lock()
	fs.IsCertHashAssociatedFuncInvoked = true
	fs.mu.Unlock()
	return fs.IsCertHashAssociatedFunc(r, hash)
}

func (fs *MDMAppleStore) AssociateCertHash(r *mdm.Request, hash string, certNotValidAfter time.Time) error {
	fs.mu.Lock()
	fs.AssociateCertHashFuncInvoked = true
	fs.mu.Unlock()
	return fs.AssociateCertHashFunc(r, hash, certNotValidAfter)
}

func (fs *MDMAppleStore) RetrieveMigrationCheckins(p0 context.Context, p1 chan<- interface{}) error {
	fs.mu.Lock()
	fs.RetrieveMigrationCheckinsFuncInvoked = true
	fs.mu.Unlock()
	return fs.RetrieveMigrationCheckinsFunc(p0, p1)
}

func (fs *MDMAppleStore) RetrieveTokenUpdateTally(ctx context.Context, id string) (int, error) {
	fs.mu.Lock()
	fs.RetrieveTokenUpdateTallyFuncInvoked = true
	fs.mu.Unlock()
	return fs.RetrieveTokenUpdateTallyFunc(ctx, id)
}

func (fs *MDMAppleStore) GetAllMDMConfigAssetsByName(ctx context.Context, assetNames []fleet.MDMAssetName) (map[fleet.MDMAssetName]fleet.MDMConfigAsset, error) {
	fs.mu.Lock()
	fs.GetAllMDMConfigAssetsByNameFuncInvoked = true
	fs.mu.Unlock()
	return fs.GetAllMDMConfigAssetsByNameFunc(ctx, assetNames)
}

func (fs *MDMAppleStore) GetABMTokenByOrgName(ctx context.Context, orgName string) (*fleet.ABMToken, error) {
	fs.mu.Lock()
	fs.GetABMTokenByOrgNameFuncInvoked = true
	fs.mu.Unlock()
	return fs.GetABMTokenByOrgNameFunc(ctx, orgName)
}

func (fs *MDMAppleStore) EnqueueDeviceLockCommand(ctx context.Context, host *fleet.Host, cmd *mdm.Command, pin string) error {
	fs.mu.Lock()
	fs.EnqueueDeviceLockCommandFuncInvoked = true
	fs.mu.Unlock()
	return fs.EnqueueDeviceLockCommandFunc(ctx, host, cmd, pin)
}

func (fs *MDMAppleStore) EnqueueDeviceWipeCommand(ctx context.Context, host *fleet.Host, cmd *mdm.Command) error {
	fs.mu.Lock()
	fs.EnqueueDeviceWipeCommandFuncInvoked = true
	fs.mu.Unlock()
	return fs.EnqueueDeviceWipeCommandFunc(ctx, host, cmd)
}
