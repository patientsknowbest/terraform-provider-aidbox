package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/patientsknowbest/terraform-provider-aidbox/aidbox"
)

func TestAccQuestionnaireThemeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccQuestionnaireThemeResourceCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccQuestionnaireThemeResourceConfig("nhs-wayfinder-theme", "Wayfinder Theme", "NHS"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aidbox_questionnaire_theme.nhs_wayfinder_theme", "id", "nhs-wayfinder-theme"),
					resource.TestCheckResourceAttr("aidbox_questionnaire_theme.nhs_wayfinder_theme", "theme_name", "Wayfinder Theme"),
					resource.TestCheckResourceAttr("aidbox_questionnaire_theme.nhs_wayfinder_theme", "design_system", "NHS"),
				),
			},
			{
				ResourceName:      "aidbox_questionnaire_theme.nhs_wayfinder_theme",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccQuestionnaireThemeResourceConfig(id, themeName, designSystem string) string {
	return fmt.Sprintf(`
resource "aidbox_questionnaire_theme" "nhs_wayfinder_theme" {
	id = "%s"
	theme_name = "%s"
	design_system = "%s"
}
`, id, themeName, designSystem)
}

func testAccQuestionnaireThemeResourceCheckDestroy(s *terraform.State) error {
	apiClient := testProvider.Meta().(*aidbox.ApiClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aidbox_questionnaire_theme" {
			continue
		}

		_, err := apiClient.GetQuestionnaireTheme(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("QuestionnaireTheme still exists")
		}
	}
	return nil
}
