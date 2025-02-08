package translator

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

type Model struct {
	client *openai.Client
	cache  *cache.Cache
}

func NewModel(apiKey string) *Model {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &Model{
		client: openai.NewClient(apiKey),
		cache:  c,
	}
}

func createCacheKey(text, sourceLanguage, targetLanguage string) string {
	hash := sha1.Sum([]byte(sourceLanguage + targetLanguage + text))
	return hex.EncodeToString(hash[:])
}

func (m *Model) TranslateText(text, sourceLanguage, targetLanguage string) (string, error) {
	cacheKey := createCacheKey(text, sourceLanguage, targetLanguage)

	if cached, found := m.cache.Get(cacheKey); found {
		return cached.(string), nil
	}

	ctx := context.Background()

	systemPrompt := fmt.Sprintf(
		"You are an expert translation assistant specialized in casual chat communications. "+
			"Translate the given text from %s to %s, preserving the original tone and cultural context. "+
			"Provide only the translated text, without any additional commentary.",
		sourceLanguage, targetLanguage,
	)

	userPrompt := fmt.Sprintf(
		"Translate the following text from %s to %s:\n\n%s",
		sourceLanguage, targetLanguage, text,
	)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:       openai.GPT4oMini,
		Messages:    messages,
		Temperature: 0.3,
	}

	maxRetries := 3
	delay := 1 * time.Second

	var resp openai.ChatCompletionResponse
	var err error
	for i := 0; i < maxRetries; i++ {
		resp, err = m.client.CreateChatCompletion(ctx, req)
		if err == nil {
			break
		}

		if strings.Contains(err.Error(), "429") {
			fmt.Printf("Received 429 error, retrying in %v (attempt %d/%d)...\n", delay, i+1, maxRetries)
			time.Sleep(delay)
			delay *= 2
		} else {
			break
		}
	}

	if err != nil {
		return "", fmt.Errorf("API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no translation received")
	}

	translatedText := resp.Choices[0].Message.Content

	m.cache.Set(cacheKey, translatedText, cache.DefaultExpiration)

	return translatedText, nil
}
