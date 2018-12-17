resource "aci_tenant" "tenant_for_ap" {
  name        = "tenant_for_ap"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_application_profile" "demo_app_profile" {
  tenant_dn                 = "${aci_tenant.tenant_for_ap.id}"
  name                      = "test_tf_ap"
  description               = "This app profile is created by terraform ACI provider"
  prio                      = "unspecified"
  relation_fv_rs_ap_mon_pol = "test1mon" # Relation to class monEPGPol. Cardinality - N_TO_ONE
}
