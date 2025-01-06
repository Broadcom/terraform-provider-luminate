resource "luminate_rdp_access_policy" "new-rdp-access-policy" {
  name =  "my-rdp-access-policy"

  identity_provider_id = "identity_provider_id"
  user_ids = ["user1_id", "user2_id"]
  group_ids = ["group1_id", "group2_id"]

  applications = ["application1_id","application2_id"]

  validators = {
    web_verification = true
  }

  conditions = {
    source_ip = ["127.0.0.1/24", "1.1.1.1/16"]
    location = ["Wallis and Futuna"]
  }
}
