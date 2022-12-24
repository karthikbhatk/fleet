parasails.registerPage('pricing', {
  //  ╦╔╗╔╦╔╦╗╦╔═╗╦    ╔═╗╔╦╗╔═╗╔╦╗╔═╗
  //  ║║║║║ ║ ║╠═╣║    ╚═╗ ║ ╠═╣ ║ ║╣
  //  ╩╝╚╝╩ ╩ ╩╩ ╩╩═╝  ╚═╝ ╩ ╩ ╩ ╩ ╚═╝
  data: {
    formData: {
      macos: 0,
      windows: 0,
      servers: 0,
      iot: 0,
      containers: 0
    },

    // For tracking client-side validation errors in our form.
    // > Has property set to `true` for each invalid property in `formData`.
    formErrors: { /* … */ },

    pricingCalculatorFormRules: {},

    // Syncing / loading state
    syncing: false,

    // Server error state
    cloudError: '',
    estimatedCost: '', // For pricing calculator
    displaySecurityPricingMode: true, // For pricing mode switch
  },

  //  ╦  ╦╔═╗╔═╗╔═╗╦ ╦╔═╗╦  ╔═╗
  //  ║  ║╠╣ ║╣ ║  ╚╦╝║  ║  ║╣
  //  ╩═╝╩╚  ╚═╝╚═╝ ╩ ╚═╝╩═╝╚═╝
  beforeMount: function() {
    //…
  },
  mounted: async function(){
    //…
  },

  //  ╦╔╗╔╔╦╗╔═╗╦═╗╔═╗╔═╗╔╦╗╦╔═╗╔╗╔╔═╗
  //  ║║║║ ║ ║╣ ╠╦╝╠═╣║   ║ ║║ ║║║║╚═╗
  //  ╩╝╚╝ ╩ ╚═╝╩╚═╩ ╩╚═╝ ╩ ╩╚═╝╝╚╝╚═╝
  methods: {
    clickChatButton: function() {
      // Temporary hack to open the chat
      // (there's currently no official API for doing this outside of React)
      //
      // > Alex: hey mike! if you're just trying to open the chat on load, we actually have a `defaultIsOpen` field
      // > you can set to `true` :) i haven't added the `Papercups.open` function to the global `Papercups` object yet,
      // > but this is basically what the functions look like if you want to try and just invoke them yourself:
      // > https://github.com/papercups-io/chat-widget/blob/master/src/index.tsx#L4-L6
      // > ~Dec 31, 2020
      window.dispatchEvent(new Event('papercups:open'));
    },
    updateEstimatedTotal: function() {
      let total = (7 * this.formData.macos) +
      (7 * this.formData.windows) +
      (7 * this.formData.servers) +
      (7 * this.formData.iot) +
      (7 * this.formData.containers);
      this.estimatedCost = total;
    },
  }
});
