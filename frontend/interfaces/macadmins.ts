export interface IDataTableMDMFormat {
  status: string;
  hosts: number;
}

export interface IMunkiVersionsAggregate {
  version: string;
  hosts_count: number;
}

export interface IMunkiIssuesAggregate {
  id: number;
  name: string;
  type: string;
  hosts_count: number;
}
export interface IMDMAggregateStatus {
  enrolled_manual_hosts_count: number;
  enrolled_automated_hosts_count: number;
  unenrolled_hosts_count: number;
}

export interface IMDMSolution {
  id: number;
  name: string | null;
  server_url: string;
  hosts_count: number;
}

export interface IMacadminAggregate {
  macadmins: {
    counts_updated_at: string;
    munki_versions: IMunkiVersionsAggregate[];
    munki_issues: IMunkiIssuesAggregate[];
    mobile_device_management_enrollment_status: IMDMAggregateStatus;
    mobile_device_management_solution: IMDMSolution[] | null;
  };
}
