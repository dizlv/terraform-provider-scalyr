// Test Parser
resource "scalyr_file" "terraform_parser" {
  path = "/logParsers/terraform"
  content = <<-EOF
  {
   formats: [
     "$who$ $said$",
   ]
  }
  EOF
}

resource "random_pet" "identifier" {
  keepers = {
    # New Pet on Every Run (causes us to validate on every run)
    id = uuid()
    # New Pet when Parser updated (causes us to validate on parser change)
    #id = scalyr_file.terraform_parser.id
  }
}

resource "scalyr_event" "validation_event" {
  depends_on = [scalyr_file.terraform_parser]
  message = "Anthony Test"
  attributes = {
    parser = "terraform"
    a_number = 42
    id = random_pet.identifier.id
  }
}

// Look for it formatted correctly
data "scalyr_query" "validate_parser" {
	depends_on = [scalyr_event.validation_event]
	expected_count = 1
	retry_count = 10
	retry_wait = 10
	start_time = "4 hours"
	query = <<-EOF
	who = "Anthony" said = "Test" id = "${random_pet.identifier.id}" | columns timestamp, message | limit 1
	EOF
}

output "good" {
  value = element(data.scalyr_query.validate_parser.results,0)
}
