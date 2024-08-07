parasails.registerPage('profiles', {
  //  ╦╔╗╔╦╔╦╗╦╔═╗╦    ╔═╗╔╦╗╔═╗╔╦╗╔═╗
  //  ║║║║║ ║ ║╠═╣║    ╚═╗ ║ ╠═╣ ║ ║╣
  //  ╩╝╚╝╩ ╩ ╩╩ ╩╩═╝  ╚═╝ ╩ ╩ ╩ ╩ ╚═╝
  data: {
    sortDirection: 'DESC',
    teamFilter: undefined,
    profilesToDisplay: [],
    platformFriendlyNames: {
      'darwin': 'macOS, iOS, and ipadOS',
      'windows': 'Windows',
      'linux': 'Linux'
    },
    selectedTeam: {},
    modal: '',
    syncing: false,

  },

  //  ╦  ╦╔═╗╔═╗╔═╗╦ ╦╔═╗╦  ╔═╗
  //  ║  ║╠╣ ║╣ ║  ╚╦╝║  ║  ║╣
  //  ╩═╝╩╚  ╚═╝╚═╝ ╩ ╚═╝╩═╝╚═╝
  beforeMount: function() {
  },
  mounted: async function() {
    this.profilesToDisplay = this.profiles;
    console.log(this.teams)
    //…
  },

  //  ╦╔╗╔╔╦╗╔═╗╦═╗╔═╗╔═╗╔╦╗╦╔═╗╔╗╔╔═╗
  //  ║║║║ ║ ║╣ ╠╦╝╠═╣║   ║ ║║ ║║║║╚═╗
  //  ╩╝╚╝ ╩ ╚═╝╩╚═╩ ╩╚═╝ ╩ ╩╚═╝╝╚╝╚═╝
  methods: {
    clickChangeSortDirection: async function() {
      if(this.sortDirection === 'ASC') {
        this.sortDirection = 'DESC';
        this.profilesToDisplay = _.sortByOrder(this.profiles, 'name', 'desc');
      } else {
        this.sortDirection = 'ASC';
        this.profilesToDisplay = _.sortByOrder(this.profiles, 'name', 'asc');
      }
      await this.forceRender();
    },
    changeTeamFilter: async function() {
      console.log(this.teamFilter);
      if(this.teamFilter){
        this.selectedTeam = _.find(this.teams, {fleetApid: this.teamFilter});
        console.log(this.selectedTeam);
        let profilesOnThisTeam = _.filter(this.profiles, (profile)=>{
          // console.log(profile.profiles);
          return _.where(profile.teams, {'fleetApid': this.selectedTeam.fleetApid}).length > 0
        })
        this.profilesToDisplay = profilesOnThisTeam;
      } else {
        this.profilesToDisplay = this.profiles;
      }
    },
    clickChangeTeamFilter: async function(teamApid) {
      this.teamFilter = teamApid;
      console.log(teamApid);
      this.selectedTeam = _.find(this.teams, {'fleetApid': teamApid});
      console.log(this.selectedTeam);
      let profilesOnThisTeam = _.filter(this.profiles, (profile)=>{
        // console.log(profile.profiles);
        return _.where(profile.teams, {'fleetApid': this.selectedTeam.fleetApid}).length > 0
      })
      this.profilesToDisplay = profilesOnThisTeam;
      console.log(profilesOnThisTeam);
    },
    clickDownloadProfile: async function(profile) {
      // Call the download profile cloud action
      // Return the downloaded profile with the correct filename
      // Or possible make these just open the download endpoint in a new tab to download it.
    },
    clickOpenEditModal: async function(profile) {
      this.modal = 'edit-profile';
    },
    clickOpenDeleteModal: async function(profile) {
      this.modal = 'delete-profile';
    },
    clickOpenAddProfileModal: async function() {
      this.modal = 'add-profile';
    }
  }
});
