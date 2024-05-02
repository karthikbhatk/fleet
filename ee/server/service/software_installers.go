package service

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/fleetdm/fleet/v4/pkg/file"
	"github.com/fleetdm/fleet/v4/server/contexts/ctxerr"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/go-kit/kit/log/level"
)

func (svc *Service) UploadSoftwareInstaller(ctx context.Context, payload *fleet.UploadSoftwareInstallerPayload) error {
	if err := svc.authz.Authorize(ctx, &fleet.SoftwareInstaller{TeamID: payload.TeamID}, fleet.ActionWrite); err != nil {
		return err
	}
	if payload == nil {
		return ctxerr.New(ctx, "payload is required")
	}

	if payload.InstallerFile == nil {
		return ctxerr.New(ctx, "installer file is required")
	}

	title, vers, hash, err := file.ExtractInstallerMetadata(payload.Filename, payload.InstallerFile)
	if err != nil {
		// TODO: confirm error handling
		if strings.Contains(err.Error(), "unsupported file type") {
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

	// checck if exists in the installer store
	exists, err := svc.softwareInstallStore.Exists(ctx, payload.StorageID)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "checking if installer exists")
	}
	if !exists {
		// reset the reader before storing (it was consumed to extract metadata)
		if _, err := payload.InstallerFile.Seek(0, 0); err != nil {
			return ctxerr.Wrap(ctx, err, "resetting installer file reader")
		}
		if err := svc.softwareInstallStore.Put(ctx, payload.StorageID, payload.InstallerFile); err != nil {
			return ctxerr.Wrap(ctx, err, "storing installer")
		}
	}

	if payload.InstallScript == "" {
		payload.InstallScript = file.GetInstallScript(payload.Filename)
	}

	// TODO: basic validation of install and post-install script (e.g., supported interpreters)?

	// TODO: any validation of pre-install query?

	source, err := fleet.SofwareInstallerSourceFromFilename(payload.Filename)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "determining source from filename")
	}
	payload.Source = source

	installerID, err := svc.ds.MatchOrCreateSoftwareInstaller(ctx, payload)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "matching or creating software installer")
	}
	level.Debug(svc.logger).Log("msg", "software installer uploaded", "installer_id", installerID)

	// TODO: QA what breaks when you have a software title with no versions?

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

	if err := svc.ds.DeleteSoftwareInstaller(ctx, id); err != nil {
		return ctxerr.Wrap(ctx, err, "deleting software installer")
	}

	return nil
}
