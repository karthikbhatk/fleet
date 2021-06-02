package authz

import input.action
import input.object
import input.subject

# Constants

# Actions
read := "read"
write := "write"
run := "run"

# Roles
admin := "admin"
maintainer := "maintainer"
observer := "observer"

# Default deny
default allow = false

# team_role gets the role that the subject has for the team, returning undefined
# if the user has no explicit role for that team.
team_role(subject, team_id) = role {
	subject_team := subject.teams[_]
	subject_team.id == team_id
	role := subject_team.role
}

##
# Global config
##

# Any logged in user can read global config
allow {
  object.type == "app_config"
  not is_null(subject)
  action == read
}

# Admin can write global config
allow {
  object.type == "app_config"
  subject.global_role == admin
  action == write
}

##
# Teams
##

# Any logged in user can read teams (service must filter appropriately based on
# access).
allow {
  object.type == "team"
  not is_null(subject)
  action == read
}

# Admin can write teams
allow {
  object.type == "team"
  subject.global_role == admin
  action == write
}

##
# Users
##

# Any user can modify self
allow {
  object.type == "user"
  object.id == subject.id
}

# Any user can read other users
allow {
  object.type == "user"
  not is_null(subject)
  action == read
}

# Admins can write all users
allow {
  object.type == "user"
  subject.global_role == admin
}

##
# Invites
##

# Only global admins may read/write invites
allow {
  object.type == "invite"
  subject.global_role == admin
	action == read
}
allow {
  object.type == "invite"
  subject.global_role == admin
	action == write
}

##
# Sessions
##

# Any user can modify own session
allow {
  object.type == "session"
  object.user_id == subject.id
}

# Admins can write all users
allow {
  object.type == "session"
  subject.global_role == admin
}

##
# Enroll Secrets
##

# Admins can read/write all
allow {
	object.type == "enroll_secret"
	subject.global_role == admin
}

# Global maintainers can read all
allow {
	object.type == "enroll_secret"
	subject.global_role == maintainer
	action == read
}

# Team maintainers can read for appropriate teams
allow {
	object.type == "enroll_secret"
	team_role(subject, object.team_id) == maintainer
	action == read
}

# (Observers are not granted read for enroll secrets)

##
# Hosts
##

# Allow read/write for global admin
allow {
	object.type == "host"
	subject.global_role = admin
	action == [read, write][_]
}

# Allow read for global maintainer
allow {
	object.type == "host"
	subject.global_role = maintainer
	action == read
}

# Allow read for global observer
allow {
	object.type == "host"
	subject.global_role = observer
	action == read
}

# Allow read for matching team maintainer/observer
allow {
	object.type == "host"
	team_role(subject, object.team_id) == maintainer
	action == read
}
allow {
	object.type == "host"
	team_role(subject, object.team_id) == observer
	action == read
}

##
# Labels
##

# All users can read labels
allow {
  not is_null(subject)
	object.type == "label"
	action == read
}

# Only global admins and maintainers can write labels
allow {
	object.type == "label"
	subject.global_role == admin
	action == [read, write][_]
}
allow {
	object.type == "label"
	subject.global_role == maintainer
	action == [read, write][_]
}

##
# Queries
##

# All users can read queries
allow {
  not is_null(subject)
	object.type == "query"
	action == read
}

# Only admins and (team) maintainers can write queries
allow {
	object.type == "query"
	subject.global_role == admin
	action == write
}
allow {
	object.type == "query"
	subject.global_role == maintainer
	action == write
}
allow {
	object.type == "query"
	# If role is maintainer on any team
	team_role(subject, subject.teams[_].id) == maintainer
	action == write
}

# Global admins and (team) maintainers can run any
allow {
	object.type == "query"
	subject.global_role == admin
	action = run
}
allow {
	object.type == "query"
	subject.global_role == maintainer
	action = run
}
allow {
	object.type == "query"
	# If role is maintainer on any team
	team_role(subject, subject.teams[_].id) == maintainer
	action == run
}
# (Team) observers can run only if observers_can_run
allow {
	object.type == "query"
	object.observer_can_run == true
	subject.global_role == observer
	action = run
}
allow {
	object.type == "query"
	object.observer_can_run == true
	# If role is observer on any team
	team_role(subject, subject.teams[_].id) == observer
	action == run
}

##
# Targets
##

allow {
  not is_null(subject)
  object.type == "target"
	action == read
}

##
# Packs
##

# Only global admins and maintainers can read/write packs
allow {
	object.type == "pack"
	subject.global_role == admin
	action == [read, write][_]
}
allow {
	object.type == "pack"
	subject.global_role == maintainer
	action == [read, write][_]
}

##
# File Carves
##

# Only global adminscan read/write carves
allow {
	object.type == "carve"
	subject.global_role == admin
	action == [read, write][_]
}
