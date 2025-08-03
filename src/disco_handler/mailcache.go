package discohandler

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	mailtmpl "github.com/yamaki-87/mailbot/src/mail_tmpl"
)

type SessionState struct {
	mail     *mailtmpl.Mail
	lastTime time.Time
}

func NewSessionState(mail *mailtmpl.Mail) SessionState {
	return SessionState{
		mail:     mail,
		lastTime: time.Now(),
	}
}

func (s *SessionState) GetMail() *mailtmpl.Mail {
	return s.mail
}

type MailStore struct {
	mu    sync.RWMutex
	store map[string]SessionState
}

func NewMailStore() *MailStore {
	return &MailStore{
		store: make(map[string]SessionState),
	}
}

func (ms *MailStore) Get(userID string) (SessionState, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	val, ok := ms.store[userID]
	return val, ok
}

func (ms *MailStore) Set(userID string, state SessionState) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.store[userID] = state
}

func (ms *MailStore) Delete(userID string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	delete(ms.store, userID)
}

func (ms *MailStore) IsEmpty() bool {
	return len(ms.store) == 0
}

// MailChache監視関数
//
// isTestから送られて物を監視し、2分経過したら削
func (ms *MailStore) TimeoutCheck(timeout time.Duration, onTimeout func(userID string)) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	now := time.Now()
	for userID, state := range ms.store {
		if now.Sub(state.lastTime) > timeout {
			delete(ms.store, userID)
			onTimeout(userID)
		}
	}
}

func StartSessionTimeoutWatcher(mailStore *MailStore) {
	go func() {
		ticker := time.NewTicker(10 * time.Second) // 10秒おきに監視
		defer ticker.Stop()

		for {
			<-ticker.C
			if mailStore.IsEmpty() {
				continue
			}
			mailStore.TimeoutCheck(2*time.Minute, func(userID string) {
				log.Info().Msgf("タイムアウト発生 userID:%s", userID)
			})
		}
	}()
}
