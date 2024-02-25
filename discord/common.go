package discord

import "github.com/bwmarrin/discordgo"

func (d *Discord) interactionRespondText(i *discordgo.InteractionCreate, context string, components []discordgo.MessageComponent) error {
	return d.s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    context,
			Components: components,
		},
	})
}

func (d *Discord) editInteractionResponseText(i *discordgo.InteractionCreate, context string, components ...[]discordgo.MessageComponent) error {
	var c *[]discordgo.MessageComponent
	if len(components) > 0 {
		c = &components[0]
	}
	_, err := d.s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:    &context,
		Components: c,
	})
	return err
}

func (d *Discord) isMaster(i *discordgo.InteractionCreate) bool {
	for _, id := range d.masters {
		if i.Member != nil {
			if id == i.Member.User.Username {
				return true
			}
		}
		if i.User != nil {
			if i.User.Username == id {
				return true
			}
		}
	}
	return false
}
