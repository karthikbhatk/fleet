# frozen_string_literal: true

require 'puppet/util/fleet_client'

Puppet::Functions.create_function(:"fleetdm::preassign_profile") do
  dispatch :preassign_profile do
    param 'String', :profile_identifier
    param 'String', :host_uuid
    param 'String', :template
    optional_param 'String', :group
    optional_param 'Enum[absent, present]', :ensure
  end

  def preassign_profile(profile_identifier, host_uuid, template, group = 'default', ensure_profile = 'present')
    client = Puppet::Util::FleetClient.instance
    run_identifier = "#{closure_scope.catalog.catalog_uuid}-#{Puppet[:node_name_value]}"
    response = client.preassign_profile(run_identifier, host_uuid, template, group, ensure_profile)

    if response['error'].empty?
      base64_checksum = Digest::MD5.base64digest(template)
      host = client.get_host_by_identifier(host_uuid)
      host_profiles = client.get_host_profiles(host['body']['host']['id'])

      if host_profiles['error'].empty?
        unless host_profiles['body']['profiles'].any? { |p| p['checksum'] == base64_checksum }
          Puppet.info("successfully pre-assigned profile #{profile_identifier}")
          response['changed'] = true
        end
      end
    else
      Puppet.err("error pre-assigning profile #{profile_identifier}: #{response['error']} \n\n #{template}")
    end

    response
  end
end
