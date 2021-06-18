// Package datastore provides testcases for Datastore implementations.
package datastore

import (
	"testing"

	"github.com/fleetdm/fleet/server/fleet"
)

// TestFunctions are the test functions that a Datastore implementation should
// run to verify proper implementation.
var TestFunctions = [...]func(*testing.T, fleet.Datastore){
	testOrgInfo,
	testAdditionalQueries,
	testEnrollSecrets,
	testEnrollSecretsCaseSensitive,
	testEnrollSecretRoundtrip,
	testEnrollSecretUniqueness,
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
	testSavePack,
	testEnrollHost,
	testAuthenticateHost,
	testAuthenticateHostCaseSensitive,
	testLabels,
	testSaveLabel,
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
	testSaveHostPackStats,
	testDeleteHost,
	testListHosts,
	testListHostsFilterAdditional,
	testListHostsStatus,
	testListHostsQuery,
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
	testHostAdditional,
	testCarveMetadata,
	testCarveBlocks,
	testCarveListCarves,
	testCarveCleanupCarves,
	testCarveUpdateCarve,
	testTeamGetSetDelete,
	testTeamUsers,
	testTeamListTeams,
	testTeamSearchTeams,
	testUserTeams,
	testUserCreateWithTeams,
	testSaveHostSoftware,
}
