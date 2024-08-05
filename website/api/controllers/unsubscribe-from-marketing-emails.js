module.exports = {


  friendlyName: 'Unsubscribe from marketing emails',


  description: 'Unsubscribes a specified email address from the nurture email automation.',


  inputs: {
    emailAddress: {
      type: 'string',
      description: 'The email address of the user who wants to unsubscribe from marketing emails.',
      required: true,
    }
  },


  exits: {
    userNotFound: {
      description: 'The provided email address could not be matched to a Fleet user account',
      responseType: 'badRequest',
    },
    success: {
      description: 'The user has opted out of markering emails',
    }
  },


  fn: async function ({emailAddress}) {

    let userRecord = await User.findOne({emailAddress: emailAddress});

    if(!userRecord){
      throw 'userNotFound';
    }
    // Update the user record for this email address to set their nurture email timestamps to 1
    // so they are excluded them from future runs of the deliver-nurture-emails script.
    // FUTURE: update the user model to have a subscribedToNurtureEmails attribute.
    await User.updateOne({emailAddress: emailAddress}).set({
      stageThreeNurtureEmailSentAt: 1,
      stageFourNurtureEmailSentAt: 1,
      stageFiveNurtureEmailSentAt: 1,
    });

    // FUTURE: update the users contact record in salesforce to indicate that they do not want to receive automated marketing emails. (If we update the stage, it might be changed by an action the user takes on the website)
    // if(sails.config.environment === 'production'){
    //   require('assert')(sails.config.custom.salesforceIntegrationUsername);
    //   require('assert')(sails.config.custom.salesforceIntegrationPasskey);

    //   // Log in to Salesforce.
    //   let jsforce = require('jsforce');
    //   let salesforceConnection = new jsforce.Connection({
    //     loginUrl : 'https://fleetdm.my.salesforce.com'
    //   });
    //   await salesforceConnection.login(sails.config.custom.salesforceIntegrationUsername, sails.config.custom.salesforceIntegrationPasskey);

    //   let existingContactRecord = await salesforceConnection.sobject('Contact')
    //   .findOne({
    //     Email:  emailAddress,
    //   });

    //   if(existingContactRecord) {
    //     //If we found an existing contact record in salesforce, update its status to be "Do not contact"
    //     let salesforceContactId = existingContactRecord.Id;
    //     await salesforceConnection.sobject('Contact')
    //     .update({
    //       Id: salesforceContactId,
    //       Stage__c: 'Do not contact',// eslint-disable-line camelcase
    //     });
    //   }
    // }
    // All done.
    return this.res.redirect('/#unsubscribed');

  }


};
