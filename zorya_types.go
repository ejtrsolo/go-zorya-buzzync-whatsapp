package zoryabuzzyncwhatsapp

// LoginRequest represents the payload for Zorya login
type LoginRequest struct {
	Username     string `json:"username"`
	UserPassword string `json:"userPassword"`
}

// LoginResponse represents the response from Zorya login
type LoginResponse struct {
	Success        bool          `json:"success"`
	Errors         interface{}   `json:"errors"`
	DetailedErrors []interface{} `json:"detailedErrors"`
	Data           string        `json:"data"` // The JWT token
}

// SendMessageRequest represents the payload for sending a message
type SendMessageRequest struct {
	Phone      string            `json:"phone"`
	Message    string            `json:"message"`
	TemplateID string            `json:"template_id,omitempty"`
	Variables  map[string]string `json:"variables,omitempty"`
}

// SendMessageResponse represents the response from sending a message
type SendMessageResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id"`
	Error     string `json:"error,omitempty"`
}

// WhatsAppMessageRequest represents the payload for sending a WhatsApp template message
type WhatsAppMessageRequest struct {
	From            string          `json:"from"`
	To              string          `json:"to"`
	TrackingMessage bool            `json:"trackingMessage,omitempty"`
	Content         WhatsAppContent `json:"content"`
}

// WhatsAppContent represents the content of the WhatsApp message
type WhatsAppContent struct {
	Type         string       `json:"type"`
	TemplateName string       `json:"templateName"`
	LanguageCode string       `json:"languageCode"`
	TemplateData TemplateData `json:"templateData"`
}

// TemplateData represents the data for the template
type TemplateData struct {
	Header           *TemplateHeader           `json:"header,omitempty"`
	Body             TemplateBody              `json:"body"`
	Footer           *TemplateFooter           `json:"footer,omitempty"`
	Buttons          []TemplateButton          `json:"buttons,omitempty"`
	Carousel         *TemplateCarousel         `json:"carousel,omitempty"`
	LimitedTimeOffer *TemplateLimitedTimeOffer `json:"limitedTimeOffer,omitempty"`
}

// TemplateHeader represents the header component
type TemplateHeader struct {
	Type     string `json:"type,omitempty"`
	MediaURL string `json:"mediaUrl,omitempty"`
	Filename string `json:"filename,omitempty"`
}

// TemplateBody represents the body component
type TemplateBody struct {
	Placeholders      []string          `json:"placeholders,omitempty"`
	NamedPlaceholders map[string]string `json:"namedPlaceholders,omitempty"`
}

// TemplateFooter represents the footer component
type TemplateFooter struct {
	Text string `json:"text,omitempty"`
}

// TemplateButton represents a button component
type TemplateButton struct {
	Type         string   `json:"type,omitempty"`
	Text         string   `json:"text,omitempty"`
	URL          string   `json:"url,omitempty"`
	MediaURL     string   `json:"mediaUrl,omitempty"`
	Placeholders []string `json:"placeholders,omitempty"`
	Parameter    string   `json:"parameter,omitempty"`
}

// TemplateCarousel represents a carousel component
type TemplateCarousel struct {
	Cards []CarouselCard `json:"cards"`
}

// CarouselCard represents a card in the carousel
type CarouselCard struct {
	Header  CarouselHeader   `json:"header"`
	Body    *CarouselBody    `json:"body,omitempty"`
	Buttons []TemplateButton `json:"buttons,omitempty"`
}

// CarouselHeader represents the header of a carousel card
type CarouselHeader struct {
	Type     string `json:"type"`
	MediaURL string `json:"mediaUrl"`
}

// CarouselBody represents the body of a carousel card
type CarouselBody struct {
	Placeholders []string `json:"placeholders,omitempty"`
}

// TemplateLimitedTimeOffer represents the limited time offer component
type TemplateLimitedTimeOffer struct {
	ExpirationTime string `json:"expirationTime,omitempty"`
}

// WhatsAppMessageResponse represents the response from sending a WhatsApp message
type WhatsAppMessageResponse struct {
	Success        bool                   `json:"success"`
	Errors         interface{}            `json:"errors"`
	DetailedErrors []interface{}          `json:"detailedErrors"`
	Data           WhatsAppMessageResData `json:"data"`
}

type WhatsAppMessageResData struct {
	TransactionID string `json:"transactionId"`
}
