/*

# Get All Teams
data "scalyr_teams" "all" {
}

# Dump all tokens from All Accounts
data "scalyr_tokens" "all" {
  for_each = toset(data.scalyr_teams.all.teams)
}

output "teams" {
  value = data.scalyr_teams.all.teams
}

output "keys" {
  value = data.scalyr_tokens.all
}

*/
