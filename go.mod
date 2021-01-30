module ansoni/terraform-provider-scalyr

go 1.14

require (
	github.com/google/uuid v1.2.0
	github.com/graphql-go/graphql v0.7.9
	github.com/hashicorp/terraform v0.12.28
	github.com/ansoni/terraform-provider-scalyr/scalyr-go v0.0.0-00010101000000-000000000000
)

replace github.com/ansoni/terraform-provider-scalyr/scalyr-go => ./scalyr-go
