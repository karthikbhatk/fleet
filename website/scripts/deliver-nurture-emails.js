module.exports = {


  friendlyName: 'Deliver nurture emails',


  description: 'Sends nurture emails to users who have been at psychological stage 3 & 4 for more than a day, and users who have been stage five for six weeks.',


  fn: async function () {

    sails.log('Running custom shell script... (`sails run deliver-nurture-emails`)');

    let nowAt = Date.now();
    let nurtureCampaignStartedAt = new Date('06-22-2024').getTime();
    let oneHourAgoAt = nowAt - (1000 * 60 * 60);
    let oneDayAgoAt = nowAt - (1000 * 60 * 60);
    let sixWeeksAgoAt = nowAt - (1000 * 60 * 60 * 24 * 7 * 6);
    // Find user records that are over an hour old that were created after July 22nd.
    let usersWithMdmBuyingSituation = await User.find({
      createdAt: {
        '>=': nurtureCampaignStartedAt,
        '<=': oneHourAgoAt,
      },
      primaryBuyingSituation: 'mdm',
    });

    // Only send emails to stage 3 users who have not received a nurture email for this stage, and that have been stage 3 for at least one day.
    let stageThreeMdmFocusedUsersWhoHaveNotReceivedAnEmail = _.filter(usersWithMdmBuyingSituation, (user)=>{
      return user.stageThreeNurtureEmailSentAt === 0
      && user.psychologicalStage === '3 - Intrigued';
    });

    // Only send emails to stage 4 users who have not received a a nurture email for this stage, and that have been stage 4 for at least one day.
    let stageFourMdmFocusedUsersWhoHaveNotReceivedAnEmail = _.filter(usersWithMdmBuyingSituation, (user)=>{
      return user.stageFourNurtureEmailSentAt === 0
      && user.psychologicalStage === '4 - Has use case';
    });

    // Only send emails to stage 5 users who have not received a nurture email for this stage, and that have been stage 5 for at least six weeks.
    let stageFiveMdmFocusedUsersWhoHaveNotReceivedAnEmail = _.filter(usersWithMdmBuyingSituation, (user)=>{
      return user.stageFiveNurtureEmailSentAt === 0
      && user.psychologicalStage === '5 - Personally confident';
    });

    for(let user of stageThreeMdmFocusedUsersWhoHaveNotReceivedAnEmail) {
      if(user.psychologicalStageLastChangedAt > oneDayAgoAt) {
        continue;
      } else {
        await sails.helpers.sendTemplateEmail.with({
          template: 'email-nurture-stage-three',
          layout: 'layout-nurture-email',
          templateData: {
            firstName: user.firstName
          },
          to: user.emailAddress,
          toName: `${user.firstName} ${user.lastName}`,
          subject: 'Was it any good?',
          bcc: [sails.config.custom.activityCaptureEmailForNutureEmails],
          from: sails.config.custom.contactEmailForNutureEmails,
        });
      }
    }

    await User.update({id: {in: _.pluck(stageThreeMdmFocusedUsersWhoHaveNotReceivedAnEmail, 'id')}})
    .set({
      stageThreeNurtureEmailSentAt: nowAt,
    });

    for(let user of stageFourMdmFocusedUsersWhoHaveNotReceivedAnEmail) {
      if(user.psychologicalStageLastChangedAt > oneDayAgoAt) {
        continue;
      } else {
        await sails.helpers.sendTemplateEmail.with({
          template: 'email-nurture-stage-four',
          layout: 'layout-nurture-email',
          templateData: {
            firstName: user.firstName
          },
          to: user.emailAddress,
          toName: `${user.firstName} ${user.lastName}`,
          subject: 'Deploy open-source MDM',
          bcc: [sails.config.custom.activityCaptureEmailForNutureEmails],
          from: sails.config.custom.contactEmailForNutureEmails,
        });
      }
    }

    await User.update({id: {in: _.pluck(stageFourMdmFocusedUsersWhoHaveNotReceivedAnEmail, 'id')}})
    .set({
      stageFourNurtureEmailSentAt: nowAt,
    });

    for(let user of stageFiveMdmFocusedUsersWhoHaveNotReceivedAnEmail) {
      if(user.psychologicalStageLastChangedAt > sixWeeksAgoAt) {
        continue;
      } else {
        await sails.helpers.sendTemplateEmail.with({
          template: 'email-nurture-stage-five',
          layout: 'layout-nurture-email',
          templateData: {
            firstName: user.firstName
          },
          to: user.emailAddress,
          toName: `${user.firstName} ${user.lastName}`,
          subject: 'Update',
          bcc: [sails.config.custom.activityCaptureEmailForNutureEmails],
          from: sails.config.custom.contactEmailForNutureEmails,
        });
      }
    }

    await User.update({id: {in: _.pluck(stageFiveMdmFocusedUsersWhoHaveNotReceivedAnEmail, 'id')}})
    .set({
      stageFiveNurtureEmailSentAt: nowAt,
    });

    // Pause for 10 seconds to allow for all of the emails to be sent. (The sendTemplateEmail helper is called with await, but it queues a background action to send the emails)
    await sails.helpers.flow.pause(10000);

  }


};

