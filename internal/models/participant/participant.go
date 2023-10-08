package participant

type Participant struct {
	Timestamp       string `csv:"Timestamp"`
	DiscordUsername string `csv:"Discord Username (formatted like kaiiorg#5074)"`
	Address         string `csv:"Shipping Address (include country)"`
	Name            string `csv:"Shipping Name"`
	Over21          string `csv:"Are you at least 21 years old?"`
	TShirtSize      string `csv:"US T-Shirt Size"`
	Allergies       string `csv:"Any Food Allergies"`
	Requests        string `csv:"Special requests or notes (for either your gifter or the henchmen)"`
}
