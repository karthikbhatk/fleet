module.exports = {


  friendlyName: 'Save questionnaire progress and continue',


  description: 'Saves the user\'s current progress in the get started questionnaire',


  inputs: {
    currentStep: {
      type: 'string',
      // isIn: [TODO]
    },
    formData: {
      type: {},
    }
  },


  exits: {

  },


  fn: async function ({currentStep, formData}) {
    console.log(formData, currentStep);
    let userRecord = await User.findOne({id: this.req.me.id});
    let questionnaireProgress;
    console.log(userRecord.getStartedQuestionnarieAnswers);
    if(!userRecord.currentGetStartedQuestionnarieStep || _.isEmpty(userRecord.getStartedQuestionnarieAnswers)) {
      questionnaireProgress = {};
    } else {
      questionnaireProgress = _.clone(userRecord.getStartedQuestionnarieAnswers);
    }
    questionnaireProgress[currentStep] = formData;
    console.log(questionnaireProgress);
    // this.req.session.getStartedProgress.previouslyAnsweredQuestions[currentStep] = formData;
    let previouslyAnsweredQuestions = _.clone(questionnaireProgress);
    // this.req.session.getStartedProgress.currentStep = currentStep;
    let nextStepInForm;

    switch(currentStep) {
      case 'start':
        nextStepInForm = 'what-are-you-using-fleet-for';
        break;
      case 'what-are-you-using-fleet-for':
        nextStepInForm = 'have-you-ever-used-fleet';
        break;
      case 'have-you-ever-used-fleet':
        let fleetUseStatus = previouslyAnsweredQuestions['have-you-ever-used-fleet'].fleetUseStatus;
        if(fleetUseStatus === 'yes-recently-deployed' || fleetUseStatus === 'yes-deployed') {
          nextStepInForm = 'how-many-hosts';
        } else if(fleetUseStatus === 'no' && previouslyAnsweredQuestions['what-are-you-using-fleet-for'].primaryBuyingSituation === 'eo-security') {
          nextStepInForm = 'what-are-you-working-on-eo-security';
        } else {
          nextStepInForm = 'welcome-to-fleet';
        }
        break;
      case 'how-many-hosts':
        if(previouslyAnsweredQuestions['how-many-hosts'].numberOfHosts === '1 to 100' ||
          previouslyAnsweredQuestions['how-many-hosts'].numberOfHosts === '100 to 700') {
          nextStepInForm = 'will-you-be-self-hosting';
        } else {
          nextStepInForm = 'lets-talk-to-your-team';
        }
        break;
      case 'will-you-be-self-hosting':
        if(previouslyAnsweredQuestions['will-you-be-self-hosting'].willSelfHost === 'true'){
          nextStepInForm = 'self-hosted-deploy';
        } else {
          nextStepInForm = 'managed-cloud-for-growing-deployments';
        }
        break;
      case 'what-are-you-working-on-eo-security':
        nextStepInForm = 'is-it-any-good';
        break;
      case 'is-it-any-good':
        nextStepInForm = 'what-did-you-think';
        break;
      case 'what-did-you-think':
        if(previouslyAnsweredQuestions['what-did-you-think'] === 'let-me-think-about-it'){
          nextStepInForm = 'is-it-any-good';
        } else {
          nextStepInForm = 'deploy-fleet-in-your-environment-eo-security';
        }
        break;
    }
    let currentGetStartedQuestionnarieStep = nextStepInForm;
    // this.req.session.getStartedProgress.currentStep = nextStepInForm;
    // All done.
    console.log('»»»»');
    console.log(this.req.session.getStartedProgress);
    let getStartedProgress = _.clone(questionnaireProgress);
    await User.updateOne({id: userRecord.id}).set({getStartedQuestionnarieAnswers: questionnaireProgress, currentGetStartedQuestionnarieStep});
    return {getStartedProgress, nextStepInForm};
  }


};
