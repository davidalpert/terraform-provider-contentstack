resource "contentstack_global_field" "common_metadata" {
  uid         = "common_metadata"
  title       = "Common Metadata"
  description = "created by terraform"
  fields = [
    {
      uid          = "title"
      display_name = "Title"
      data_type    = "text"
    }
  ]
}
