package authz

import input.action
import input.object
import input.subject

# Default deny
default allow = false

# team_role gets the role that the subject has for the team, returning undefined
# if the user has no explicit role for that team.
team_role(subject, team_id) = role {
	subject_team := subject.teams[_]
	subject_team.team_id == team_id
	role := subject_team.role
}

##
# Enroll Secrets
##

# Admins can read/write all
allow {
	object.type == "enroll_secret"
	subject.global_role == "admin"
}

# Global maintainers can read all
allow {
	object.type == "enroll_secret"
	subject.global_role == "maintainer"
	action == "read"
}

# Team maintainers can read for appropriate teams
allow {
	object.type == "enroll_secret"
	action == "read"
	team_role(subject, object.team_id) == "maintainer"
}

# (Observers are not granted read for enroll secrets)

##
# Hosts
##

# Allow read/write for global admin
allow {
	object.type == "host"
	subject.global_role = "admin"
	action == ["read", "write"][_]
}

# Allow read/write for global maintainer
allow {
	object.type == "host"
	subject.global_role = "maintainer"
	action == ["read", "write"][_]
}

# Allow read for global observer
allow {
	object.type == "host"
	subject.global_role = "maintainer"
	action == "read"
}

# Allow read/write for matching team maintainer
allow {
	object.type == "host"
	subject.global_role = "maintainer"
	action == ["read", "write"][_]
	team_role(subject, object.team) == "maintainer"
}

##
# Labels
##

# All users can read labels
allow {
	object.type == "label"
	action == "read"
}

# Only global admins and maintainers can write labels
allow {
	object.type == "label"
	write_roles := {"admin", "maintainer"}
	write_roles[subject.global_role]
}

##
# Organization
##

# Only global admins can access organization
allow {
	object.type == "organization"
	subject.global_role == "admin"
}
