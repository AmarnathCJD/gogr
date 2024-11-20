// Copyright (c) 2024 RoseLoverX

package telegram

import "fmt"

type InlineSendOptions struct {
	Gallery      bool   `json:"gallery,omitempty"`
	NextOffset   string `json:"next_offset,omitempty"`
	CacheTime    int32  `json:"cache_time,omitempty"`
	Private      bool   `json:"private,omitempty"`
	SwitchPm     string `json:"switch_pm,omitempty"`
	SwitchPmText string `json:"switch_pm_text,omitempty"`
}

func (c *Client) AnswerInlineQuery(QueryID int64, Results []InputBotInlineResult, Options ...*InlineSendOptions) (bool, error) {
	options := getVariadic(Options, &InlineSendOptions{})
	options.CacheTime = getValue(options.CacheTime, 60)
	request := &MessagesSetInlineBotResultsParams{
		Gallery:    options.Gallery,
		Private:    options.Private,
		QueryID:    QueryID,
		Results:    Results,
		CacheTime:  options.CacheTime,
		NextOffset: options.NextOffset,
	}
	if options.SwitchPm != "" {
		request.SwitchPm = &InlineBotSwitchPm{
			Text:       options.SwitchPm,
			StartParam: getValue(options.SwitchPmText, "start"),
		}
	}
	resp, err := c.MessagesSetInlineBotResults(request)
	if err != nil {
		return false, err
	}
	return resp, nil
}

type CallbackOptions struct {
	Alert     bool   `json:"alert,omitempty"`
	CacheTime int32  `json:"cache_time,omitempty"`
	URL       string `json:"url,omitempty"`
}

func (c *Client) AnswerCallbackQuery(QueryID int64, Text string, Opts ...*CallbackOptions) (bool, error) {
	options := getVariadic(Opts, &CallbackOptions{})
	request := &MessagesSetBotCallbackAnswerParams{
		QueryID: QueryID,
		Message: Text,
		Alert:   options.Alert,
	}
	if options.URL != "" {
		request.URL = options.URL
	}
	if options.CacheTime != 0 {
		request.CacheTime = options.CacheTime
	}
	resp, err := c.MessagesSetBotCallbackAnswer(request)
	if err != nil {
		return false, err
	}
	return resp, nil
}

// BOT COMMANDS

func (c *Client) SetBotCommands(commands []*BotCommand, scope *BotCommandScope, languageCode ...string) (bool, error) {
	resp, err := c.BotsSetBotCommands(*scope, getVariadic(languageCode, "en"), commands)
	if err != nil {
		return false, err
	}
	return resp, nil
}

func (c *Client) SetBotDefaultPrivileges(privileges *ChatAdminRights, ForChannels ...bool) (resp bool, err error) {
	forCh := getVariadic(ForChannels, true)
	if forCh {
		resp, err = c.BotsSetBotBroadcastDefaultAdminRights(privileges)
		return
	}
	resp, err = c.BotsSetBotGroupDefaultAdminRights(privileges)
	return
}

func (c *Client) SetChatMenuButton(userID int64, button *BotMenuButton) (bool, error) {
	peer, err := c.ResolvePeer(userID)
	if err != nil {
		return false, err
	}
	peerUser, ok := peer.(*InputPeerUser)
	if !ok {
		return false, fmt.Errorf("invalid user")
	}
	resp, err := c.BotsSetBotMenuButton(&InputUserObj{AccessHash: peerUser.AccessHash, UserID: peerUser.UserID}, *button)
	if err != nil {
		return false, err
	}
	return resp, nil
}
