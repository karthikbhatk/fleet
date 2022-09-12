parasails.registerPage('osquery-tables', {
  //  ╦╔╗╔╦╔╦╗╦╔═╗╦    ╔═╗╔╦╗╔═╗╔╦╗╔═╗
  //  ║║║║║ ║ ║╠═╣║    ╚═╗ ║ ╠═╣ ║ ║╣
  //  ╩╝╚╝╩ ╩ ╩╩ ╩╩═╝  ╚═╝ ╩ ╩ ╩ ╩ ╚═╝
  data: {
    //…
    selectedPlatform: 'all',
    search: '',
    showTableNav: false,
    userFriendlyPlatformNames: {
      'darwin': 'macOS',
      'linux': 'Linux',
      'windows': 'Windows',
      'all': 'All platforms'
    },
  },

  computed: {
    filteredTables: function () {
      return this.allTables.filter(
        (table) =>
          this._isIncluded(table.platforms, this.selectedPlatform) &&
          this._isIncluded(table.title, this.search)
      );
    },
    numberOfTablesDisplayed: function() {
      return this.filteredTables.length;
    },
  },

  //  ╦  ╦╔═╗╔═╗╔═╗╦ ╦╔═╗╦  ╔═╗
  //  ║  ║╠╣ ║╣ ║  ╚╦╝║  ║  ║╣
  //  ╩═╝╩╚  ╚═╝╚═╝ ╩ ╚═╝╩═╝╚═╝
  beforeMount: function() {

  },
  mounted: async function() {
    // sort the array of all tables
    this.allTables = this.allTables.sort((a, b)=>{
      if(a.name < b.name){
        return -1;
      }
    });
    let keywordsForThisTable = [];
    if(this.tableToDisplay.keywordsForSyntaxHighlighting){
      keywordsForThisTable = this.tableToDisplay.keywordsForSyntaxHighlighting;
    }
    (function highlightThatSyntax(){
      $('pre code').each((i, block) => {
        window.hljs.highlightBlock(block);
      });
      $('.hljs').each((i, el)=>{
        for(let keyword of keywordsForThisTable) {
          if(_.includes(el.innerHTML, keyword)) {
            let replacementHMTL = el.innerHTML.replaceAll(keyword, '<span class="hljs-attr">'+keyword+'</span>');
            $(el).html(replacementHMTL);
          }
        }
      });
    })();
    // Adjust the height of the sidebar navigation to match the height of the html partial
    (function adjustSideBarHeight(){
      let tablePartialHeight = $('[purpose="table-container"]').height();
      $('[purpose="table-of-contents"]').css({'max-height': tablePartialHeight - 120});
      // $('[purpose="left-sidebar"]').css({'height': tablePartialHeight});
    })();
  },
  //  ╦╔╗╔╔╦╗╔═╗╦═╗╔═╗╔═╗╔╦╗╦╔═╗╔╗╔╔═╗
  //  ║║║║ ║ ║╣ ╠╦╝╠═╣║   ║ ║║ ║║║║╚═╗
  //  ╩╝╚╝ ╩ ╚═╝╩╚═╩ ╩╚═╝ ╩ ╩╚═╝╝╚╝╚═╝
  methods: {
    //…
    clickFilterByPlatform: async function(platform) {
      this.selectedPlatform = platform;
    },

    _isIncluded: function (data, selectedOption) {
      if (selectedOption.startsWith('all') || selectedOption === '') {
        return true;
      }
      if (_.isArray(data)) {
        data = data.join(', ');
      }
      return (
        _.isString(data) && data.toLowerCase().includes(selectedOption.toLowerCase())
      );
    },

    toggleTableNav: function() {
      this.showTableNav = !this.showTableNav;
    },
  }
});
