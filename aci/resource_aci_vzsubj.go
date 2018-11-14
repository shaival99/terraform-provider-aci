package aci

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceAciContractSubject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciContractSubjectCreate,
		Update: resourceAciContractSubjectUpdate,
		Read:   resourceAciContractSubjectRead,
		Delete: resourceAciContractSubjectDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciContractSubjectImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"contract_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"cons_match_t": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "consumer subject match criteria",

				ValidateFunc: validation.StringInSlice([]string{
					"All",
					"AtleastOne",
					"AtmostOne",
					"None",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"prio": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "priority level specifier",

				ValidateFunc: validation.StringInSlice([]string{
					"level1",
					"level2",
					"level3",
					"unspecified",
				}, false),
			},

			"prov_match_t": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "consumer subject match criteria",

				ValidateFunc: validation.StringInSlice([]string{
					"All",
					"AtleastOne",
					"AtmostOne",
					"None",
				}, false),
			},

			"rev_flt_ports": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "enables filter to apply on ingress and egress traffic",

				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"target_dscp": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "target dscp",

				ValidateFunc: validation.StringInSlice([]string{
					"AF11",
					"AF12",
					"AF13",
					"AF21",
					"AF22",
					"AF23",
					"AF31",
					"AF32",
					"AF33",
					"AF41",
					"AF42",
					"AF43",
					"CS0",
					"CS1",
					"CS2",
					"CS3",
					"CS4",
					"CS5",
					"CS6",
					"CS7",
					"EF",
					"VA",
					"unspecified",
				}, false),
			},

			"relation_vz_rs_subj_graph_att": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to vnsAbsGraph",
			},
			"relation_vz_rs_subj_filt_att": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vzFilter",
				Set:         schema.HashString,
			},
		}),
	}
}

func getRemoteContractSubject(client *client.Client, dn string) (*models.ContractSubject, error) {
	vzSubjCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzSubj := models.ContractSubjectFromContainer(vzSubjCont)

	if vzSubj.DistinguishedName == "" {
		return nil, fmt.Errorf("Bridge Domain %s not found", vzSubj.DistinguishedName)
	}

	return vzSubj, nil
}

func setContractSubjectAttributes(vzSubj *models.ContractSubject, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(vzSubj.DistinguishedName)
	d.Set("description", vzSubj.Description)
	d.Set("contract_dn", GetParentDn(vzSubj.DistinguishedName))
	vzSubjMap, _ := vzSubj.ToMap()

	d.Set("cons_match_t", vzSubjMap["consMatchT"])
	d.Set("name_alias", vzSubjMap["nameAlias"])
	d.Set("prio", vzSubjMap["prio"])
	d.Set("prov_match_t", vzSubjMap["provMatchT"])
	d.Set("rev_flt_ports", vzSubjMap["revFltPorts"])
	d.Set("target_dscp", vzSubjMap["targetDscp"])
	return d
}

func resourceAciContractSubjectImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	aciClient := m.(*client.Client)

	dn := d.Id()

	vzSubj, err := getRemoteContractSubject(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setContractSubjectAttributes(vzSubj, d)
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciContractSubjectCreate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	ContractDn := d.Get("contract_dn").(string)

	vzSubjAttr := models.ContractSubjectAttributes{}
	if ConsMatchT, ok := d.GetOk("cons_match_t"); ok {
		vzSubjAttr.ConsMatchT = ConsMatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzSubjAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		vzSubjAttr.Prio = Prio.(string)
	}
	if ProvMatchT, ok := d.GetOk("prov_match_t"); ok {
		vzSubjAttr.ProvMatchT = ProvMatchT.(string)
	}
	if RevFltPorts, ok := d.GetOk("rev_flt_ports"); ok {
		vzSubjAttr.RevFltPorts = RevFltPorts.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		vzSubjAttr.TargetDscp = TargetDscp.(string)
	}
	vzSubj := models.NewContractSubject(fmt.Sprintf("subj-%s", name), ContractDn, desc, vzSubjAttr)

	err := aciClient.Save(vzSubj)
	if err != nil {
		return err
	}

	if relationTovzRsSubjGraphAtt, ok := d.GetOk("relation_vz_rs_subj_graph_att"); ok {
		relationParam := relationTovzRsSubjGraphAtt.(string)
		err = aciClient.CreateRelationvzRsSubjGraphAtt(vzSubj.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	if relationTovzRsSubjFiltAtt, ok := d.GetOk("relation_vz_rs_subj_filt_att"); ok {
		relationParamList := toStringList(relationTovzRsSubjFiltAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationvzRsSubjFiltAtt(vzSubj.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
		}
	}

	d.SetId(vzSubj.DistinguishedName)
	return resourceAciContractSubjectRead(d, m)
}

func resourceAciContractSubjectUpdate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)
	ContractDn := d.Get("contract_dn").(string)

	vzSubjAttr := models.ContractSubjectAttributes{}
	if ConsMatchT, ok := d.GetOk("cons_match_t"); ok {
		vzSubjAttr.ConsMatchT = ConsMatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzSubjAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		vzSubjAttr.Prio = Prio.(string)
	}
	if ProvMatchT, ok := d.GetOk("prov_match_t"); ok {
		vzSubjAttr.ProvMatchT = ProvMatchT.(string)
	}
	if RevFltPorts, ok := d.GetOk("rev_flt_ports"); ok {
		vzSubjAttr.RevFltPorts = RevFltPorts.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		vzSubjAttr.TargetDscp = TargetDscp.(string)
	}
	vzSubj := models.NewContractSubject(fmt.Sprintf("subj-%s", name), ContractDn, desc, vzSubjAttr)

	vzSubj.Status = "modified"

	err := aciClient.Save(vzSubj)

	if err != nil {
		return err
	}
	if d.HasChange("relation_vz_rs_subj_graph_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_subj_graph_att")
		err = aciClient.DeleteRelationvzRsSubjGraphAtt(vzSubj.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvzRsSubjGraphAtt(vzSubj.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_vz_rs_subj_filt_att") {
		oldRel, newRel := d.GetChange("relation_vz_rs_subj_filt_att")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationvzRsSubjFiltAtt(vzSubj.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationvzRsSubjFiltAtt(vzSubj.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}
	}

	d.SetId(vzSubj.DistinguishedName)
	return resourceAciContractSubjectRead(d, m)

}

func resourceAciContractSubjectRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	dn := d.Id()
	vzSubj, err := getRemoteContractSubject(aciClient, dn)

	if err != nil {
		return err
	}
	setContractSubjectAttributes(vzSubj, d)
	return nil
}

func resourceAciContractSubjectDelete(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzSubj")
	if err != nil {
		return err
	}

	d.SetId("")
	return err
}
