resource "luminate_site" "new-site" {
  name = "tfAccSite"
}

resource "luminate_dns_group_resiliency" "new-dns-group" {
  name = "testDNSGroupResiliency"
  sendNotifications = true
  domainSuffixes = ["somedomain.com"]
}

data "luminate_dns_server_resiliency" "new-dns-server-resiliency" {
  name = "testDNSServerResiliency"
  site_id = "${luminate_site.new-site.id}"
  group_id = "${luminate_dns_group_resiliency.new-dns-group.id}"
  internal_address = "udp://20.0.0.1:63"
}
