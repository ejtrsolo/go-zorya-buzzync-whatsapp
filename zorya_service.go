package zoryabuzzyncwhatsapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ejtrsolo/go-zorya-buzzync-whatsapp/tools"
)

// ZoryaService handles interactions with the Zorya API
type ZoryaService struct {
	BaseURL  string
	Token    string
	Username string
	Password string
}

// NewZoryaService creates a new instance of ZoryaService
func NewZoryaService(baseURL string, username string, password string) *ZoryaService {
	return &ZoryaService{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
	}
}

// Login authenticates with the Zorya API and retrieves a token
func (s *ZoryaService) Login() error {
	// If we already have a token, we don't need to login again
	if s.Token != "" {
		return nil
	}

	if s.Username == "" || s.Password == "" {
		return errors.New("zorya credentials not configured")
	}

	reqBody := LoginRequest{
		Username:     s.Username,
		UserPassword: s.Password,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	response := tools.ConsultClient(s.BaseURL+"/api/v1/User/login", "POST", nil, bytes.NewBuffer(jsonData), true)

	if !response.Success {
		return fmt.Errorf("login failed: %s", response.Message)
	}

	if response.Data["status_code"].(int) != http.StatusOK {
		return fmt.Errorf("login failed with status: %d", response.Data["status_code"])
	}

	bodyStr, ok := response.Data["body"].(string)
	if !ok {
		return errors.New("invalid response body format")
	}

	var loginResp LoginResponse
	if err := json.Unmarshal([]byte(bodyStr), &loginResp); err != nil {
		return err
	}

	if !loginResp.Success {
		return fmt.Errorf("login failed: %v", loginResp.Errors)
	}

	if loginResp.Data != "" {
		s.Token = loginResp.Data
	} else {
		return errors.New("no token received from Zorya login")
	}

	return nil
}

// SendWhatsAppTemplateMessage sends a WhatsApp template message via Zorya API
func (s *ZoryaService) SendWhatsAppTemplateMessage(req WhatsAppMessageRequest) (*WhatsAppMessageResponse, error) {
	if s.Token == "" {
		if err := s.Login(); err != nil {
			return nil, err
		}
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Authorization": "Bearer " + s.Token,
	}

	url := s.BaseURL + "/api/v1/WhatsApp/messages"

	// Log de los parámetros de la petición
	log.Printf("[ZORYA] Request URL: %s", url)
	log.Printf("[ZORYA] Request Method: POST")
	log.Printf("[ZORYA] Request Headers: %+v", headers)
	log.Printf("[ZORYA] Request Body: %s", string(jsonData))

	response := tools.ConsultClient(url, "POST", headers, bytes.NewBuffer(jsonData), true)

	// Log de la respuesta (opcional)
	log.Printf("[ZORYA] Response Status: %v", response.Data["status_code"])
	log.Printf("[ZORYA] Response Success: %v", response.Success)

	if !response.Success {
		log.Printf("[ZORYA] Response Error: %s", response.Message)
		return nil, fmt.Errorf("send message failed: %s", response.Message)
	}

	bodyStr, ok := response.Data["body"].(string)
	if !ok {
		return nil, errors.New("invalid response body format")
	}

	var sendResp WhatsAppMessageResponse
	if err := json.Unmarshal([]byte(bodyStr), &sendResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	statusCode := response.Data["status_code"].(int)
	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		log.Printf("[ZORYA] Error Status Code: %d, Response: %s", statusCode, bodyStr)
		return &sendResp, fmt.Errorf("send message failed with status: %d", statusCode)
	}

	log.Printf("[ZORYA] Message sent successfully")
	return &sendResp, nil
}
