resource "aci_tenant" "tenant_for_ap" {
  name        = "tenant_for_ap"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_app_profile" "demo_app_profile" {
  fv_tenant_dn   = "${aci_tenant.tenant_for_ap.id}"
  name        = "test_tf_ap"
  description = "This app profile is created by terraform ACI provider"
}