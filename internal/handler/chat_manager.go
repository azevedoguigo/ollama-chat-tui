package handler

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/azevedoguigo/ollama-chat-tui/internal/storage"
	"github.com/google/uuid"
)

type ChatManager struct {
	chats     map[string]*storage.ChatSession
	mutex     sync.RWMutex
	configDir string
	chatsDir  string
}

func NewChatManager(configDir, chatsDir string) (*ChatManager, error) {
	chats, err := storage.LoadChats(configDir, chatsDir)
	if err != nil {
		chats = make(map[string]*storage.ChatSession)
	}
	return &ChatManager{
		chats:     chats,
		configDir: configDir,
		chatsDir:  chatsDir,
	}, nil
}

func (cm *ChatManager) AddChat(title, model string) *storage.ChatSession {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	chat := &storage.ChatSession{
		ID:        uuid.New(),
		Title:     title,
		Messages:  []storage.Message{},
		Model:     model,
		CreatedAt: time.Now(),
	}
	cm.chats[chat.ID.String()] = chat

	if err := storage.SaveChat(cm.configDir, cm.chatsDir, chat); err != nil {
		fmt.Printf("Error when saving new chat: %v\n", err)
	}

	return chat
}

func (cm *ChatManager) GetChatByID(id string) (*storage.ChatSession, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	chat, exists := cm.chats[id]
	return chat, exists
}

func (cm *ChatManager) GetAllChats() []*storage.ChatSession {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	var chatsList []*storage.ChatSession
	for _, chat := range cm.chats {
		chatsList = append(chatsList, chat)
	}
	sort.Slice(chatsList, func(i, j int) bool {
		return chatsList[i].CreatedAt.Before(chatsList[j].CreatedAt)
	})

	return chatsList
}

func (cm *ChatManager) AppendMessage(chatID, role, content string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	chat, exists := cm.chats[chatID]
	if !exists {
		return fmt.Errorf("chat with ID %s not found", chatID)
	}

	chat.Messages = append(chat.Messages, storage.Message{
		Role:    role,
		Content: content,
	})

	return storage.SaveChat(cm.configDir, cm.chatsDir, chat)
}

func (cm *ChatManager) UpdateLastMessage(chatID, content string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	chat, exists := cm.chats[chatID]
	if !exists {
		return fmt.Errorf("chat with ID %s not found", chatID)
	}
	if len(chat.Messages) == 0 {
		return fmt.Errorf("the chat has no messages to update")
	}
	chat.Messages[len(chat.Messages)-1].Content += content

	return storage.SaveChat(cm.configDir, cm.chatsDir, chat)
}

func (cm *ChatManager) DeleteChat(chatID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	chat, exists := cm.chats[chatID]
	if !exists {
		return fmt.Errorf("chat with ID %s not found", chatID)
	}
	if err := storage.DeleteChat(chat, cm.chats); err != nil {
		return err
	}

	delete(cm.chats, chatID)
	return nil
}
