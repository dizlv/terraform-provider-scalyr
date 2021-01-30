/*
# Query a specific configuration
data "scalyr_file" "test" {
  path = "/logParsers/test123"
}

# A New Configuration
resource "scalyr_file" "test" {
  path = "/test/test"
  content = <<-EOF
Hello there
  EOF
}

# Output the content
output "file_content" {
  value = data.scalyr_file.test.content
}
*/
