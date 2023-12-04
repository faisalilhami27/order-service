package template

import "order-service/config"

func GetTemplateIDByName(name string) *string {
	templates := config.Config.InternalService.Notification.Templates
	for _, template := range templates {
		if template.Name == name {
			return &template.TemplateID
		}
	}
	return nil
}
