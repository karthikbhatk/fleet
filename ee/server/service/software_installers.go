package service

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/fleetdm/fleet/v4/pkg/file"
	"github.com/fleetdm/fleet/v4/pkg/fleethttp"
	"github.com/fleetdm/fleet/v4/server/contexts/ctxerr"
	hostctx "github.com/fleetdm/fleet/v4/server/contexts/host"
	"github.com/fleetdm/fleet/v4/server/contexts/viewer"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/go-kit/log/level"
	"golang.org/x/sync/errgroup"
)

func (svc *Service) UploadSoftwareInstaller(ctx context.Context, payload *fleet.UploadSoftwareInstallerPayload) error {
	if err := svc.authz.Authorize(ctx, &fleet.SoftwareInstaller{TeamID: payload.TeamID}, fleet.ActionWrite); err != nil {
		return err
	}

	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return fleet.ErrNoContext
	}

	if err := svc.addMetadataToSoftwarePayload(ctx, payload); err != nil {
		return ctxerr.Wrap(ctx, err, "adding metadata to payload")
	}

	if err := svc.storeSoftware(ctx, payload); err != nil {
		return ctxerr.Wrap(ctx, err, "storing software installer")
	}

	// TODO: basic validation of install and post-install script (e.g., supported interpreters)?
	// TODO: any validation of pre-install query?

	installerID, err := svc.ds.MatchOrCreateSoftwareInstaller(ctx, payload)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "matching or creating software installer")
	}
	level.Debug(svc.logger).Log("msg", "software installer uploaded", "installer_id", installerID)

	// TODO: QA what breaks when you have a software title with no versions?

	var teamName *string
	if payload.TeamID != nil {
		t, err := svc.ds.Team(ctx, *payload.TeamID)
		if err != nil {
			return err
		}
		teamName = &t.Name
	}

	// Create activity
	if err := svc.ds.NewActivity(ctx, vc.User, fleet.ActivityTypeAddedSoftware{
		SoftwareTitle:   payload.Title,
		SoftwarePackage: payload.Filename,
		TeamName:        teamName,
		TeamID:          payload.TeamID,
	}); err != nil {
		return ctxerr.Wrap(ctx, err, "creating activity for added software")
	}

	return nil
}

func (svc *Service) DeleteSoftwareInstaller(ctx context.Context, id uint) error {
	// get the software installer to have its team id
	meta, err := svc.ds.GetSoftwareInstallerMetadata(ctx, id)
	if err != nil {
		if fleet.IsNotFound(err) {
			// couldn't get the metadata to have its team, authorize with a no-team
			// as a fallback - the requested installer does not exist so there's
			// no way to know what team it would be for, and returning a 404 without
			// authorization would leak the existing/non existing ids.
			if err := svc.authz.Authorize(ctx, &fleet.SoftwareInstaller{}, fleet.ActionWrite); err != nil {
				return err
			}
			return ctxerr.Wrap(ctx, err, "getting software installer metadata")
		}
	}

	// do the actual authorization with the software installer's team id
	if err := svc.authz.Authorize(ctx, &fleet.SoftwareInstaller{TeamID: meta.TeamID}, fleet.ActionWrite); err != nil {
		return err
	}

	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return fleet.ErrNoContext
	}

	if err := svc.ds.DeleteSoftwareInstaller(ctx, id); err != nil {
		return ctxerr.Wrap(ctx, err, "deleting software installer")
	}

	var teamName *string
	if meta.TeamID != nil {
		t, err := svc.ds.Team(ctx, *meta.TeamID)
		if err != nil {
			return ctxerr.Wrap(ctx, err, "getting team name for deleted software")
		}
		teamName = &t.Name
	}

	if err := svc.ds.NewActivity(ctx, vc.User, fleet.ActivityTypeDeletedSoftware{
		SoftwareTitle:   meta.SoftwareTitle,
		SoftwarePackage: meta.Name,
		TeamName:        teamName,
		TeamID:          meta.TeamID,
	}); err != nil {
		return ctxerr.Wrap(ctx, err, "creating activity for deleted software")
	}

	return nil
}

func (svc *Service) GetSoftwareInstallerMetadata(ctx context.Context, installerID uint) (*fleet.SoftwareInstaller, error) {
	// first do a basic authorization check, any logged in user can read teams
	if err := svc.authz.Authorize(ctx, &fleet.Team{}, fleet.ActionRead); err != nil {
		return nil, err
	}

	// get the installer's metadata
	meta, err := svc.ds.GetSoftwareInstallerMetadata(ctx, installerID)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "getting software installer metadata")
	}

	// authorize with the software installer's team id
	if err := svc.authz.Authorize(ctx, &fleet.SoftwareInstaller{TeamID: meta.TeamID}, fleet.ActionRead); err != nil {
		return nil, err
	}

	return meta, nil
}

func (svc *Service) DownloadSoftwareInstaller(ctx context.Context, installerID uint) (*fleet.DownloadSoftwareInstallerPayload, error) {
	meta, err := svc.GetSoftwareInstallerMetadata(ctx, installerID)
	if err != nil {
		return nil, err
	}

	return svc.getSoftwareInstallerBinary(ctx, meta.StorageID, meta.Name)
}

func (svc *Service) OrbitDownloadSoftwareInstaller(ctx context.Context, installerID uint) (*fleet.DownloadSoftwareInstallerPayload, error) {
	// this is not a user-authenticated endpoint
	svc.authz.SkipAuthorization(ctx)

	// TODO: confirm error handling

	host, ok := hostctx.FromContext(ctx)
	if !ok {
		return nil, fleet.OrbitError{Message: "internal error: missing host from request context"}
	}

	// get the installer's metadata
	meta, err := svc.ds.GetSoftwareInstallerMetadata(ctx, installerID)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "getting software installer metadata")
	}

	// ensure it cannot get access to a different team's installer
	var hTeamID uint
	if host.TeamID != nil {
		hTeamID = *host.TeamID
	}
	if (meta.TeamID != nil && *meta.TeamID != hTeamID) || (meta.TeamID == nil && hTeamID != 0) {
		return nil, ctxerr.Wrap(ctx, fleet.OrbitError{}, "host team does not match installer team")
	}

	return svc.getSoftwareInstallerBinary(ctx, meta.StorageID, meta.Name)
}

func (svc *Service) getSoftwareInstallerBinary(ctx context.Context, storageID string, filename string) (*fleet.DownloadSoftwareInstallerPayload, error) {
	// check if the installer exists in the store
	exists, err := svc.softwareInstallStore.Exists(ctx, storageID)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "checking if installer exists")
	}
	if !exists {
		return nil, ctxerr.Wrap(ctx, err, "does not exist in software installer store")
	}

	// get the installer from the store
	installer, size, err := svc.softwareInstallStore.Get(ctx, storageID)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "getting installer from store")
	}

	return &fleet.DownloadSoftwareInstallerPayload{
		Filename:  filename,
		Installer: installer,
		Size:      size,
	}, nil
}

func (svc *Service) InstallSoftwareTitle(ctx context.Context, hostID uint, softwareTitleID uint) error {
	// we need to use ds.Host because ds.HostLite doesn't return the orbit
	// node key
	host, err := svc.ds.Host(ctx, hostID)
	if err != nil {
		// if error is because the host does not exist, check first if the user
		// had access to install software (to prevent leaking valid host ids).
		if fleet.IsNotFound(err) {
			if err := svc.authz.Authorize(ctx, &fleet.HostSoftwareInstallerResultAuthz{}, fleet.ActionWrite); err != nil {
				return err
			}
		}
		svc.authz.SkipAuthorization(ctx)
		return ctxerr.Wrap(ctx, err, "get host")
	}

	if host.OrbitNodeKey == nil || *host.OrbitNodeKey == "" {
		// fleetd is required to install software so if the host is
		// enrolled via plain osquery we return an error
		svc.authz.SkipAuthorization(ctx)
		// TODO(roberto): for cleanup task, confirm with product error message.
		return fleet.NewUserMessageError(errors.New("Host doesn't have fleetd installed"), http.StatusUnprocessableEntity)
	}

	// authorize with the host's team
	if err := svc.authz.Authorize(ctx, &fleet.HostSoftwareInstallerResultAuthz{HostTeamID: host.TeamID}, fleet.ActionWrite); err != nil {
		return err
	}

	installer, err := svc.ds.GetSoftwareInstallerForTitle(ctx, softwareTitleID, host.TeamID)
	if err != nil {
		if fleet.IsNotFound(err) {
			return &fleet.BadRequestError{
				Message: "Software title has no package added. Please add software package to install.",
				InternalErr: ctxerr.WrapWithData(
					ctx, err, "couldn't find an installer for software title",
					map[string]any{"host_id": host.ID, "team_id": host.TeamID, "title_id": softwareTitleID},
				),
			}
		}

		return ctxerr.Wrap(ctx, err, "finding software installer for title")
	}

	ext := filepath.Ext(installer.Name)
	var requiredPlatform string
	switch ext {
	case ".msi", ".exe":
		requiredPlatform = "windows"
	case ".pkg":
		requiredPlatform = "darwin"
	case ".deb":
		requiredPlatform = "linux"
	default:
		// this should never happen
		return ctxerr.Errorf(ctx, "software installer has unsupported type %s", ext)
	}

	if host.FleetPlatform() != requiredPlatform {
		return &fleet.BadRequestError{
			Message: fmt.Sprintf("Package (%s) can be installed only on %s hosts.", ext, requiredPlatform),
			InternalErr: ctxerr.WrapWithData(
				ctx, err, "invalid host platform for requested installer",
				map[string]any{"host_id": host.ID, "team_id": host.TeamID, "title_id": softwareTitleID},
			),
		}
	}

	err = svc.ds.InsertSoftwareInstallRequest(ctx, hostID, installer.InstallerID)
	return ctxerr.Wrap(ctx, err, "inserting software install request")
}

func (svc *Service) GetSoftwareInstallResults(ctx context.Context, resultUUID string) (*fleet.HostSoftwareInstallerResult, error) {
	// Basic auth check
	if err := svc.authz.Authorize(ctx, &fleet.Host{}, fleet.ActionList); err != nil {
		return nil, err
	}

	res, err := svc.ds.GetSoftwareInstallResults(ctx, resultUUID)
	if err != nil {
		return nil, err
	}

	// Team specific auth check
	if err := svc.authz.Authorize(ctx, &fleet.HostSoftwareInstallerResultAuthz{HostTeamID: res.HostTeamID}, fleet.ActionRead); err != nil {
		return nil, err
	}

	return res, nil
}

func (svc *Service) storeSoftware(ctx context.Context, payload *fleet.UploadSoftwareInstallerPayload) error {
	// check if exists in the installer store
	exists, err := svc.softwareInstallStore.Exists(ctx, payload.StorageID)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "checking if installer exists")
	}
	if !exists {
		if err := svc.softwareInstallStore.Put(ctx, payload.StorageID, payload.InstallerFile); err != nil {
			return ctxerr.Wrap(ctx, err, "storing installer")
		}
	}

	return nil
}

func (svc *Service) addMetadataToSoftwarePayload(ctx context.Context, payload *fleet.UploadSoftwareInstallerPayload) error {
	if payload == nil {
		return ctxerr.New(ctx, "payload is required")
	}

	if payload.InstallerFile == nil {
		return ctxerr.New(ctx, "installer file is required")
	}

	title, vers, ext, hash, err := file.ExtractInstallerMetadata(payload.InstallerFile)
	if err != nil {
		if errors.Is(err, file.ErrUnsupportedType) {
			return &fleet.BadRequestError{
				Message:     "The file should be .pkg, .msi, .exe or .deb.",
				InternalErr: ctxerr.Wrap(ctx, err, "extracting metadata from installer"),
			}
		}
		return ctxerr.Wrap(ctx, err, "extracting metadata from installer")
	}
	payload.Title = title
	payload.Version = vers
	payload.StorageID = hex.EncodeToString(hash)

	// reset the reade (it was consumed to extract metadata)
	if _, err := payload.InstallerFile.Seek(0, 0); err != nil {
		return ctxerr.Wrap(ctx, err, "resetting installer file reader")
	}

	if payload.InstallScript == "" {
		payload.InstallScript = file.GetInstallScript(ext)
	}

	source, err := fleet.SofwareInstallerSourceFromExtension(ext)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "determining source from extension")
	}
	payload.Source = source

	return nil
}

const maxInstallerSizeBytes int64 = 1024 * 1024 * 500

func (svc *Service) BatchSetSoftwareInstallers(ctx context.Context, tmName string, payloads []fleet.SoftwareInstallerPayload, dryRun bool) error {
	if tmName == "" {
		svc.authz.SkipAuthorization(ctx) // so that the error message is not replaced by "forbidden"
		return ctxerr.Wrap(ctx, fleet.NewInvalidArgumentError("team_name", "must not be empty"))
	}

	tm, err := svc.ds.TeamByName(ctx, tmName)
	if err != nil {
		// If this is a dry run, the team may not have been created yet
		if dryRun && fleet.IsNotFound(err) {
			return nil
		}
		return err
	}

	if err := svc.authz.Authorize(ctx, &fleet.SoftwareInstaller{TeamID: &tm.ID}, fleet.ActionWrite); err != nil {
		return ctxerr.Wrap(ctx, err, "validating authorization")
	}

	g, workerCtx := errgroup.WithContext(ctx)
	g.SetLimit(3)
	installers := make([]*fleet.UploadSoftwareInstallerPayload, len(payloads))
	// any duplicate name in the provided set results in an error
	// byName := make(map[string]bool, len(payloads))

	for i, p := range payloads {
		i, p := i, p

		g.Go(func() error {
			client := fleethttp.NewClient()
			client.Transport = fleethttp.NewSizeLimitTransport(maxInstallerSizeBytes)
			req, err := http.NewRequestWithContext(workerCtx, http.MethodGet, p.URL, nil)
			if err != nil {
				return err
			}

			resp, err := client.Do(req)
			if err != nil {
				if errors.Is(err, fleethttp.ErrMaxSizeExceeded) {
					return fleet.NewInvalidArgumentError(
						"software.url",
						fmt.Sprintf("Couldn't edit software. URL (%q) doesn't exist. The maximum file size is %d MB", p.URL, maxInstallerSizeBytes/(1024*1024)),
					)
				}

				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusNotFound {
				return fleet.NewInvalidArgumentError(
					"software.url",
					fmt.Sprintf("Couldn't edit software. URL (%q) doesn't exist. Please make sure that URLs are publicy accessible to the internet.", p.URL),
				)
			}

			// Allow all 2xx and 3xx status codes in this pass.
			if resp.StatusCode > 400 {
				return fleet.NewInvalidArgumentError(
					"software.url",
					fmt.Sprintf("Couldn't edit software. URL (%q) received response status code %d.", p.URL, resp.StatusCode),
				)
			}

			// TODO(roberto): this reads the entire body, it's not
			// that bad since it's limited to
			// maxInstallerSizeBytes, but this could be changed so
			// `fleet.UploadSoftwarePayload` takes an io.Reader (vs
			// an io.ReadSeeker.) That requires changes in other
			// downstream methods that we use to extract metadata,
			// store the installer, etc.
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ctxerr.Wrapf(ctx, err, "reading installer %q contents", p.URL)
			}

			var filename string
			cdh, ok := resp.Header["Content-Disposition"]
			if ok && len(cdh) > 0 {
				_, params, err := mime.ParseMediaType(cdh[0])
				if err == nil {
					filename = params["filename"]
				}
			}

			// TODO: use payload fields to figure out the right extension
			// if it fails, try to extract it from the URL
			if filename == "" {
				filename = file.ExtractFilenameFromURLPath(p.URL, ".pkg")
			}

			installer := &fleet.UploadSoftwareInstallerPayload{
				Filename:          filename,
				TeamID:            &tm.ID,
				InstallScript:     p.InstallScript,
				PreInstallQuery:   p.PreInstallQuery,
				PostInstallScript: p.PostInstallScript,
				InstallerFile:     bytes.NewReader(bodyBytes),
			}

			if err := svc.addMetadataToSoftwarePayload(ctx, installer); err != nil {
				return err
			}

			installers[i] = installer

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		// NOTE: intentionally not wrapping to avoid polluting user
		// errors.
		return err
	}

	if dryRun {
		return nil
	}

	for _, payload := range installers {
		if err := svc.storeSoftware(ctx, payload); err != nil {
			return err
		}
	}

	if err := svc.ds.BatchSetSoftwareInstallers(ctx, &tm.ID, installers); err != nil {
		return ctxerr.Wrap(ctx, err, "batch set software installers")
	}

	// Note: per @noahtalerman we don't want activity items for CLI actions
	// anymore, so that's intentionally skipped.

	return nil
}
