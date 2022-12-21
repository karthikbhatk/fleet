/**
 * VantaConnection.js
 *
 * @description :: An organization who is a customer of Vanta.
 * @docs        :: https://sailsjs.com/docs/concepts/models-and-orm/models
 */

module.exports = {

  attributes: {

    //  ╔═╗╦═╗╦╔╦╗╦╔╦╗╦╦  ╦╔═╗╔═╗
    //  ╠═╝╠╦╝║║║║║ ║ ║╚╗╔╝║╣ ╚═╗
    //  ╩  ╩╚═╩╩ ╩╩ ╩ ╩ ╚╝ ╚═╝╚═╝
    emailAddress: {
      description: 'The email address provided when this Vanta connection was created.',
      type: 'string',
      required: true,
      isEmail: true,
    },

    vantaSourceId: {
      description: 'The generated source ID for this Vanta Connection.',
      type: 'string',
      unique: true,
      required: true,
    },

    fleetInstanceUrl: {
      description: 'The full URL of the Fleet instance that will be connected to Vanta.',
      type: 'string',
      required: true,
      unique: true,
    },

    fleetApiKey: {
      type: 'string',
      required: true,
      description: 'The token used to authenticate requests to the user\'s Fleet instance.',
      extendedDescription: 'This token must be for an API-only user and needs to have admin privileges on the user\'s Fleet instance'
    },

    vantaAuthToken: {
      type: 'string',
      description: 'A token used to authorize requests to Vanta on behalf of this Vanta customer.'
    },

    vantaAuthTokenExpiresAt: {
      type: 'number',
      description: 'A JS timestamp of when this connection\'s authorization token will expire.'
    },

    vantaRefreshToken: {
      type: 'string',
      description: 'The token used to request new authorization tokens for this Vanta connection.'
    },

    isConnectedToVanta: {
      type: 'boolean',
      defaultsTo: false,
      description: 'whether this Vanta connection is authorized to send data to Vanta on behalf of the user.',
    }

    //  ╔═╗╔╦╗╔╗ ╔═╗╔╦╗╔═╗
    //  ║╣ ║║║╠╩╗║╣  ║║╚═╗
    //  ╚═╝╩ ╩╚═╝╚═╝═╩╝╚═╝


    //  ╔═╗╔═╗╔═╗╔═╗╔═╗╦╔═╗╔╦╗╦╔═╗╔╗╔╔═╗
    //  ╠═╣╚═╗╚═╗║ ║║  ║╠═╣ ║ ║║ ║║║║╚═╗
    //  ╩ ╩╚═╝╚═╝╚═╝╚═╝╩╩ ╩ ╩ ╩╚═╝╝╚╝╚═╝

  },

};

