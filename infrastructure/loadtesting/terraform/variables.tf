variable "tag" {
  description = "The tag to deploy. This would be the same as the branch name"
}

variable "fleet_config" {
  description = "The configuration to use for fleet itself, gets translated as environment variables"
  type        = map(string)
  default     = {}
}

variable "loadtest_containers" {
  description = "The number of containers to loadtest with"
  type        = number
  default     = 0
}
