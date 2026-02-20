package port

import "context"

// Messenger abstracts outgoing message sending (Telegram or other).
type Messenger interface {
	SendText(ctx context.Context, chatID int64, text string) error
	SendTextWithKeyboard(ctx context.Context, chatID int64, text string, options [][]string) error
}
