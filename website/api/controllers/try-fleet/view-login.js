module.exports = {


  friendlyName: 'View login',


  description: 'Display "Login" page.',


  exits: {

    success: {
      viewTemplatePath: 'pages/try-fleet/login'
    },

    redirect: {
      description: 'The requesting user is already logged in.',
      responseType: 'redirect'
    }


  },


  fn: async function () {

    // If the user is logged in, redirect them to the Fleet sandbox page.
    if (this.req.me) {
      throw {redirect: '/try-fleet/sandbox'};
    }

    // Respond with view.
    return {};

  }


};
