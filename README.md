# terraform-provider-scalyr

The Scalyr Terraform Provider allows you to provision assets within Scalyr.

# Build and Install

[Golang is required](https://go.dev/doc/install).

	git clone git@github.com:ansoni/terraform-provider-scalyr
	go build 
	mkdir -p ~/.terraform.d/plugins
	cp terraform-provider-scalyr ~/.terraform.d/plugins/
	
# Provider Configuration

Different tokens are required for different resources. Having a `write_config_token` and a `write_log_token` will allow every resource.

## Schema

### Optional

- `endpoint` (String) Scalyr Server
- `read_config_token` (String) Scalyr ReadConfig API Token
- `read_log_token` (String) Scalyr ReadLog API Token
- `region` (String) Scalyr Region
- `write_config_token` (String) Scalyr WriteConfig API Token
- `write_log_token` (String) Scalyr WriteLog API Token


# Resources

## scalyr_file

Create a Scalyr Configuration file.

### Schema

#### Required

- `content` (String)
- `path` (String)

#### Read-Only

- `create_date` (String)
- `id` (String) The ID of this resource.
- `mod_date` (String)
- `version` (Number)

## scalyr_event

Send an event to Scalyr

### Schema

#### Optional

- `attributes` (Map of String)
- `message` (String)
- `parser` (String)

#### Read-Only

- `id` (String) The ID of this resource.

# Data Sources

## data scalyr_file

Read contents of a Scalyr Configuration.

### Schema

#### Required

- `path` (String)

#### Read-Only

- `content` (String)
- `create_date` (String)
- `id` (String) The ID of this resource.
- `mod_date` (String)
- `version` (Number)

## data scalyr_query

Perform a query and assert on results.

### Schema

#### Required

- `query` (String)

#### Optional

- `end_time` (String)
- `expected_count` (Number)
- `max_count` (Number)
- `query_type` (String)
- `retry_count` (Number)
- `retry_wait` (Number)
- `start_time` (String)

#### Read-Only

- `id` (String) The ID of this resource.
- `results` (List of Map of String)

## data scalyr_teams

Read all teams available to the configured token.

### Schema

#### Read-Only

- `id` (String) The ID of this resource.
- `teams` (Set of String)


# Examples

Available in [examples](./example)
