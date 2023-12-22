module.exports = {


  friendlyName: 'View pricing',


  description: 'Display "Pricing" page.',


  exits: {

    success: {
      viewTemplatePath: 'pages/pricing'
    },

    badConfig: {
      responseType: 'badConfig'
    },

  },


  fn: async function () {

    if(!_.isObject(sails.config.builtStaticContent) || !_.isArray(sails.config.builtStaticContent.pricingTable)) {
      throw {badConfig: 'builtStaticContent.pricingTable'};
    }
    let pricingTableFeatures = sails.config.builtStaticContent.pricingTable;

    let pricingTable = [];
    let pricingTableCategories = ['Support', 'Deployment', 'Integrations', 'Endpoint operations', 'Device management', 'Vulnerability management'];
    for(let category of pricingTableCategories) {
      // Get all the features in that have a pricingTableFeatures array that contains this category.
      let featuresInThisCategory = _.filter(pricingTableFeatures, (feature)=>{
        return _.contains(feature.pricingTableCategories, category);
      });
      // Build a dictionary containing the category name, and all features in the category, sorting premium features to the bottom of the list.
      let allFeaturesInThisCategory = {
        categoryName: category,
        features: _.sortBy(featuresInThisCategory, (feature)=>{
          return feature.tier !== 'Free';
        })
      };
      // Add the dictionaries to the arrays that we'll use to build the features table.
      pricingTable.push(allFeaturesInThisCategory);
    }
    // Respond with view.
    return {
      pricingTable
    };

  }


};
