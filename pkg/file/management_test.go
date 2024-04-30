package file

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetInstallAndRemoveScript(t *testing.T) {
	supportedTypes := []InstallerType{
		InstallerTypeMsi,
		InstallerTypePkg,
		InstallerTypeDeb,
		InstallerTypeExe,
	}

	for _, st := range supportedTypes {
		fileName := "foo/bar baz.f"
		script := GetInstallScript(st, fileName)
		require.NotEmpty(t, script)
		require.Contains(t, script, fileName)

		script = GetRemoveScript(st, fileName)
		require.NotEmpty(t, script)
		require.Contains(t, script, fileName)
	}
}
