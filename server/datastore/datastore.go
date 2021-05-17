// Package datastore provides testcases for Datastore implementations.
package datastore

import (
	"testing"

	"github.com/fleetdm/fleet/server/kolide"
)

// TestFunctions are the test functions that a Datastore implementation should
// run to verify proper implementation.
var TestFunctions = [...]func(*testing.T, kolide.Datastore){
	testOrgInfo,
	testAdditionalQueries,
	testEnrollSecrets,
	testEnrollSecretsCaseSensitive,
	testEnrollSecretRoundtrip,
	testCreateInvite,
	testInviteByEmail,
	testInviteByToken,
	testListInvites,
	testDeleteInvite,
	testDeleteQuery,
	testDeleteQueries,
	testSaveQuery,
	testListQuery,
	testDeletePack,
	testNewPack,
	testEnrollHost,
	testAuthenticateHost,
	testAuthenticateHostCaseSensitive,
	testLabels,
	testSaveLabel,
	testManagingLabelsOnPacks,
	testPasswordResetRequests,
	testCreateUser,
	testSaveUser,
	testUserByID,
	testListUsers,
	testPasswordResetRequests,
	testSearchHosts,
	testSearchHostsLimit,
	testSearchLabels,
	testSearchLabelsLimit,
	testListHostsInLabel,
	testListUniqueHostsInLabels,
	testSaveHosts,
	testDeleteHost,
	testListHosts,
	testListHostsFilterAdditional,
	testListHostsStatus,
	testListHostsQuery,
	testListHostsInPack,
	testListPacksForHost,
	testHostIDsByName,
	testHostByIdentifier,
	testAddHostsToTeam,
	testListPacks,
	testDistributedQueryCampaign,
	testCleanupDistributedQueryCampaigns,
	testBuiltInLabels,
	testLoadPacksForQueries,
	testScheduledQuery,
	testDeleteScheduledQuery,
	testNewScheduledQuery,
	testListScheduledQueriesInPack,
	testCascadingDeletionOfQueries,
	testGetPackByName,
	testGetQueryByName,
	testAddLabelToPackTwice,
	testGenerateHostStatusStatistics,
	testMarkHostSeen,
	testMarkHostsSeen,
	testCleanupIncomingHosts,
	testDuplicateNewQuery,
	testChangeEmail,
	testChangeLabelDetails,
	testMigrationStatus,
	testUnicode,
	testCountHostsInTargets,
	testHostStatus,
	testHostIDsInTargets,
	testApplyQueries,
	testApplyPackSpecRoundtrip,
	testApplyPackSpecMissingQueries,
	testApplyPackSpecMissingName,
	testGetPackSpec,
	testApplyLabelSpecsRoundtrip,
	testGetLabelSpec,
	testLabelIDsByName,
	testListLabelsForPack,
	testHostAdditional,
	testCarveMetadata,
	testCarveBlocks,
	testCarveListCarves,
	testCarveCleanupCarves,
	testCarveUpdateCarve,
	testTeamGetSetDelete,
	testTeamUsers,
	testUserTeams,
	testUserCreateWithTeams,
	testSaveHostSoftware,
}
