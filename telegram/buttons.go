// Copyright (c) 2024 RoseLoverX

package telegram

import (
	"bytes"
	"fmt"
	"strings"
)

type Button struct{}

type KeyboardBuilder struct {
	rows []*KeyboardButtonRow
}

// NewKeyboard initializes a new keyboard builder.
func NewKeyboard() *KeyboardBuilder {
	return &KeyboardBuilder{}
}

// AddRow adds a new row of buttons to the keyboard.
func (kb *KeyboardBuilder) AddRow(buttons ...KeyboardButton) *KeyboardBuilder {
	kb.rows = append(kb.rows, &KeyboardButtonRow{Buttons: buttons})
	return kb
}

// NewGrid arranges buttons into a grid based on specified rows (x) and columns (y).
// If there are fewer buttons than x*y, the last row may contain fewer buttons.
func (kb *KeyboardBuilder) NewGrid(x, y int, buttons ...KeyboardButton) *KeyboardBuilder {
	totalButtons := len(buttons)
	for i := 0; i < x && i*y < totalButtons; i++ {
		endIndex := (i + 1) * y
		if endIndex > totalButtons {
			endIndex = totalButtons
		}
		rowButtons := buttons[i*y : endIndex]
		kb.AddRow(rowButtons...)
	}

	if totalButtons > x*y {
		kb.AddRow(buttons[x*y:]...)
	}

	return kb
}

// NewColumn arranges buttons into a grid based on specified number of buttons (x) per column.
func (kb *KeyboardBuilder) NewColumn(x int, buttons ...KeyboardButton) *KeyboardBuilder {
	// i.e x buttons per column
	for i := 0; i < len(buttons); i += x {
		endIndex := i + x
		if endIndex > len(buttons) {
			endIndex = len(buttons)
		}
		kb.AddRow(buttons[i:endIndex]...)
	}
	return kb
}

// NewRow arranges buttons into a grid based on specified number of buttons (y) per row.
func (kb *KeyboardBuilder) NewRow(y int, buttons ...KeyboardButton) *KeyboardBuilder {
	// i.e y buttons per row
	for i := 0; i < y; i++ {
		var rowButtons []KeyboardButton
		for j := i; j < len(buttons); j += y {
			rowButtons = append(rowButtons, buttons[j])
		}
		kb.AddRow(rowButtons...)
	}
	return kb
}

// Build finalizes the keyboard and returns the inline markup.
func (kb *KeyboardBuilder) Build() *ReplyInlineMarkup {
	return &ReplyInlineMarkup{Rows: kb.rows}
}

func (Button) Force(placeHolder string) *ReplyKeyboardForceReply {
	return &ReplyKeyboardForceReply{Placeholder: placeHolder}
}

func (Button) Auth(text, url, forwardText string, buttonID int32) *KeyboardButtonURLAuth {
	return &KeyboardButtonURLAuth{Text: text, URL: url, FwdText: forwardText, ButtonID: buttonID}
}

func (Button) URL(text, url string) *KeyboardButtonURL {
	return &KeyboardButtonURL{Text: text, URL: url}
}

func (Button) Data(text, data string) *KeyboardButtonCallback {
	return &KeyboardButtonCallback{Text: text, Data: []byte(data)}
}

func (Button) RequestLocation(text string) *KeyboardButtonRequestGeoLocation {
	return &KeyboardButtonRequestGeoLocation{Text: text}
}

func (Button) Buy(text string) *KeyboardButtonBuy {
	return &KeyboardButtonBuy{Text: text}
}

func (Button) Game(text string) *KeyboardButtonGame {
	return &KeyboardButtonGame{Text: text}
}

func (Button) RequestPhone(text string) *KeyboardButtonRequestPhone {
	return &KeyboardButtonRequestPhone{Text: text}
}

func (Button) RequestPeer(text string, buttonID int32, peerType RequestPeerType, max ...int32) *KeyboardButtonRequestPeer {
	return &KeyboardButtonRequestPeer{Text: text, ButtonID: buttonID, PeerType: peerType, MaxQuantity: getVariadic(max, int32(0))}
}

func (Button) RequestPoll(text string, quiz bool) *KeyboardButtonRequestPoll {
	return &KeyboardButtonRequestPoll{Text: text, Quiz: quiz}
}

func (Button) SwitchInline(text string, samePeer bool, query string) *KeyboardButtonSwitchInline {
	return &KeyboardButtonSwitchInline{Text: text, SamePeer: samePeer, Query: query}
}

func (Button) WebView(text, url string) *KeyboardButtonSimpleWebView {
	return &KeyboardButtonSimpleWebView{Text: text, URL: url}
}

func (Button) Mention(text string, userID int64) *KeyboardButtonUserProfile {
	return &KeyboardButtonUserProfile{Text: text, UserID: userID}
}

func (Button) Row(Buttons ...KeyboardButton) *KeyboardButtonRow {
	return &KeyboardButtonRow{Buttons: Buttons}
}

func (Button) Keyboard(Rows ...*KeyboardButtonRow) *ReplyInlineMarkup {
	return &ReplyInlineMarkup{Rows: Rows}
}

func (Button) Clear() *ReplyKeyboardHide {
	return &ReplyKeyboardHide{}
}

type ClickOptions struct {
	Game     bool
	Password string
}

// Click clicks a button in a message.
//
// If no argument is given, the first button will be clicked.
//
// If an argument is provided, it can be one of the following:
//   - The text of the button to click.
//   - The data of the button to click.
//   - The coordinates of the button to click as a slice of integers [x, y].
func (m *NewMessage) Click(options ...any) (*MessagesBotCallbackAnswer, error) {
	requestParams := &MessagesGetBotCallbackAnswerParams{
		Peer:  m.Peer,
		MsgID: m.ID,
		Game:  false,
	}

	if len(options) > 0 {
		if opt, ok := options[0].(*ClickOptions); ok {
			requestParams.Game = opt.Game
			if opt.Password != "" {
				accountPasswordSrp, err := m.Client.AccountGetPassword()
				if err != nil {
					return nil, err
				}

				password, err := GetInputCheckPassword(opt.Password, accountPasswordSrp)
				if err != nil {
					return nil, err
				}

				requestParams.Password = password
			}
		}
	}

	if m.ReplyMarkup() == nil {
		return nil, fmt.Errorf("replyMarkup: message has no buttons")
	}

	if messageButtons, ok := (*m.ReplyMarkup()).(*ReplyInlineMarkup); ok {
		if len(messageButtons.Rows) == 0 {
			return nil, fmt.Errorf("replyMarkup: rows are empty")
		}

		switch len(options) {
		case 0:
			if len(messageButtons.Rows[0].Buttons) == 0 {
				return nil, fmt.Errorf("replyMarkup: row(0) has no buttons")
			}

			if button, ok := messageButtons.Rows[0].Buttons[0].(*KeyboardButtonCallback); ok {
				requestParams.Data = button.Data
			}

		case 1:
			currentX := 0
			currentY := 0
			for _, row := range messageButtons.Rows {
				for _, button := range row.Buttons {
					switch opt := options[0].(type) {
					case string:
						if button, ok := button.(*KeyboardButtonCallback); ok && strings.EqualFold(button.Text, opt) {
							requestParams.Data = button.Data
						}
					case []byte:
						if button, ok := button.(*KeyboardButtonCallback); ok && bytes.Equal(button.Data, opt) {
							requestParams.Data = button.Data
						}
					case int, int32, int64:
						if optInt, ok := opt.(int); ok && optInt == currentX {
							if button, ok := button.(*KeyboardButtonCallback); ok {
								requestParams.Data = button.Data
							}
						}

					case []int, []int32, []int64:
						if optInts, ok := opt.([]int); ok && len(optInts) == 2 {
							if optInts[0] == currentX && optInts[1] == currentY {
								if button, ok := button.(*KeyboardButtonCallback); ok {
									requestParams.Data = button.Data
								}
							}
						}

					default:
						return nil, fmt.Errorf("replyMarkup: invalid argument type (expected string, []byte, int, or []int)")
					}
					currentY++
				}
				currentX++
				currentY = 0
			}
		}
	}

	if requestParams.Data == nil {
		return nil, fmt.Errorf("replyMarkup: button with given (text, data, or coordinates) not found")
	}

	return m.Client.MessagesGetBotCallbackAnswer(requestParams)
}
