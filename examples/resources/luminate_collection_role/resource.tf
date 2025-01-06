resource "luminate_collection_role" "policy-owner" {
  role_type = "PolicyOwner"
  identity_provider_id =  "local"
  entity_id = "a8a48219-835f-4183-a2a9-bbba8cad8eb8"
  entity_type = "User"
  collection_id = "${luminate_collection.collection.id}"
}
