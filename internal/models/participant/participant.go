package participant

import (
	"fmt"
)

type Participant struct {
	Timestamp       string `csv:"Timestamp"`
	DiscordUsername string `csv:"Discord Username"`
	Address         string `csv:"Shipping Address (include country)"`
	Name            string `csv:"Shipping Name"`
	Over21          string `csv:"Are you at least 21 years old?"`
	TShirtSize      string `csv:"US T-Shirt Size"`
	Allergies       string `csv:"Any Food Allergies"`
	Requests        string `csv:"Special requests or notes (for either your gifter or the henchmen)"`
}

func (p *Participant) DumpGifterMessage(eventName string, gifter *Participant) string {
	message := fmt.Sprintf("%s\n", gifter.DiscordUsername)
	message += fmt.Sprintf("You are to send a gift to %s for the %s event in the DERPGRU discord server.\n", p.DiscordUsername, eventName)
	message += fmt.Sprintf("Here is all of the information they submitted:\n```\n")
	message += fmt.Sprintf("Address: %s\n", p.Address)
	if p.Name != "" {
		message += fmt.Sprintf("Name: %s\n", p.Name)
	}
	message += fmt.Sprintf("Over 21: %s\n", p.Over21)
	message += fmt.Sprintf("T Shirt Size: %s\n", p.TShirtSize)
	if p.Allergies != "" {
		message += fmt.Sprintf("Allergies: %s\n", p.Allergies)
	}
	if p.Requests != "" {
		message += fmt.Sprintf("Requests: \n%s\n", p.Requests)
	}
	message += "```\n\n"
	return message
}
