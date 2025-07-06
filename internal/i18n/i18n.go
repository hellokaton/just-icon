package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed locales/*.json
var localesFS embed.FS

// Language represents supported languages
type Language string

const (
	English Language = "en"
	Chinese Language = "zh"
)

// Messages contains all translatable text
type Messages struct {
	data map[string]string
}

// NewMessages creates a new Messages instance from JSON data
func NewMessages(jsonData []byte) (*Messages, error) {
	var data map[string]string
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}
	return &Messages{data: data}, nil
}

// Get returns the message for the given key
func (m *Messages) Get(key string) string {
	if value, exists := m.data[key]; exists {
		return value
	}
	return key // Return key as fallback if not found
}

// Localizer handles internationalization
type Localizer struct {
	currentLanguage Language
	messages        map[Language]*Messages
}

// NewLocalizer creates a new localizer instance
func NewLocalizer(language Language) *Localizer {
	l := &Localizer{
		currentLanguage: language,
		messages:        make(map[Language]*Messages),
	}
	
	// Load messages for all supported languages
	l.loadMessages()
	
	return l
}

// GetCurrentLanguage returns the current language
func (l *Localizer) GetCurrentLanguage() Language {
	return l.currentLanguage
}

// SetLanguage changes the current language
func (l *Localizer) SetLanguage(language Language) {
	l.currentLanguage = language
}

// T returns the translated text for the given key
func (l *Localizer) T(key string) string {
	messages := l.messages[l.currentLanguage]
	if messages == nil {
		// Fallback to English if current language is not available
		messages = l.messages[English]
		if messages == nil {
			return key // Return key if no translation found
		}
	}

	return messages.Get(key)
}

// Tf returns the translated text with formatting
func (l *Localizer) Tf(key string, args ...interface{}) string {
	template := l.T(key)
	return fmt.Sprintf(template, args...)
}

// loadMessages loads all language files
func (l *Localizer) loadMessages() {
	// Try to load from external files first (for development/customization)
	l.tryLoadExternalFile(English, "internal/i18n/locales/en.json")
	l.tryLoadExternalFile(Chinese, "internal/i18n/locales/zh.json")

	// If external files failed, load from embedded files
	if l.messages[English] == nil {
		l.loadEmbeddedFile(English, "locales/en.json")
	}
	if l.messages[Chinese] == nil {
		l.loadEmbeddedFile(Chinese, "locales/zh.json")
	}
}

// tryLoadExternalFile tries to load a language file from external path
func (l *Localizer) tryLoadExternalFile(lang Language, filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return // File doesn't exist or can't be read
	}

	messages, err := NewMessages(data)
	if err != nil {
		return // Parsing failed
	}

	l.messages[lang] = messages
}

// loadEmbeddedFile loads a language file from embedded filesystem
func (l *Localizer) loadEmbeddedFile(lang Language, filePath string) {
	data, err := localesFS.ReadFile(filePath)
	if err != nil {
		// This should never happen if files are properly embedded
		fmt.Printf("Warning: Failed to load embedded file %s: %v\n", filePath, err)
		return
	}

	messages, err := NewMessages(data)
	if err != nil {
		fmt.Printf("Warning: Failed to parse embedded file %s: %v\n", filePath, err)
		return
	}

	l.messages[lang] = messages
}



// Global localizer instance
var globalLocalizer *Localizer

// InitLocalizer initializes the global localizer
func InitLocalizer(language Language) {
	globalLocalizer = NewLocalizer(language)
}

// GetLocalizer returns the global localizer instance
func GetLocalizer() *Localizer {
	if globalLocalizer == nil {
		globalLocalizer = NewLocalizer(English) // Default to English
	}
	return globalLocalizer
}

// T is a convenience function for translation
func T(key string) string {
	return GetLocalizer().T(key)
}

// Tf is a convenience function for formatted translation
func Tf(key string, args ...interface{}) string {
	return GetLocalizer().Tf(key, args...)
}
