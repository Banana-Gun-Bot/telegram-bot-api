package tgbotapi

import "testing"

// A request parameter the server does not recognise is ignored rather than
// rejected, so these pin the names of the parameters added in Bot API 10.3.

func TestRichMessageConfigSendsEphemeralMessageParameters(t *testing.T) {
	config := RichMessageConfig{
		ChatID:      76918703,
		RichMessage: InputRichMessage{HTML: "<p>filled</p>"},
		EphemeralMessageParameters: &EphemeralMessageParameters{
			ReceiverUserID:              76918703,
			CallbackQueryID:             "4382bfdwdsb323b2d9",
			ReplaceCallbackQueryMessage: true,
		},
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("building the parameters failed: %v", err)
	}

	want := `{"receiver_user_id":76918703,"callback_query_id":"4382bfdwdsb323b2d9","replace_callback_query_message":true}`
	if got := params["ephemeral_message_parameters"]; got != want {
		t.Errorf("ephemeral_message_parameters = %s, want %s", got, want)
	}

	// An ordinary rich message must not turn into an ephemeral one by default.
	params, err = RichMessageConfig{ChatID: 76918703}.params()
	if err != nil {
		t.Fatalf("building the parameters failed: %v", err)
	}

	if _, ok := params["ephemeral_message_parameters"]; ok {
		t.Error("a message with no ephemeral parameters must not send the field")
	}
}

func TestRichMessageDraftConfigSendsTheStopParameters(t *testing.T) {
	config := RichMessageDraftConfig{
		ChatID:      76918703,
		DraftID:     7,
		RichMessage: InputRichMessage{Blocks: []InputRichBlock{InputRichBlockThinking{Type: "thinking"}}},
		CanStop:     true,
		KeepOnStop:  true,
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("building the parameters failed: %v", err)
	}

	for key, want := range map[string]string{"can_stop": "true", "keep_on_stop": "true"} {
		if got := params[key]; got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}

	params, err = RichMessageDraftConfig{ChatID: 76918703, DraftID: 7}.params()
	if err != nil {
		t.Fatalf("building the parameters failed: %v", err)
	}

	for _, key := range []string{"can_stop", "keep_on_stop"} {
		if _, ok := params[key]; ok {
			t.Errorf("a draft that offers no stop button must not send %s", key)
		}
	}
}

func TestPromoteChatMemberConfigSendsCanSendWelcomeMessages(t *testing.T) {
	config := PromoteChatMemberConfig{
		ChatMemberConfig:       ChatMemberConfig{ChatID: 76918703, UserID: 76918703},
		CanSendWelcomeMessages: true,
	}

	params, err := config.params()
	if err != nil {
		t.Fatalf("building the parameters failed: %v", err)
	}

	if got := params["can_send_welcome_messages"]; got != "true" {
		t.Errorf("can_send_welcome_messages = %q, want true", got)
	}
}
