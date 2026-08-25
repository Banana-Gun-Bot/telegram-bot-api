package tgbotapi

import (
	"encoding/json"
	"strings"
	"testing"
)

// Bot API 10.3 additions to the reply markup. Both are optional flags the
// server ignores when the field name is wrong, so the wire shape is pinned.

func TestInlineKeyboardMarkupCarriesForceReply(t *testing.T) {
	markup := InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{{NewInlineKeyboardButtonData("Buy", "buy:BONK")}},
		ForceReply:     true,
	}

	encoded, err := json.Marshal(markup)
	if err != nil {
		t.Fatalf("marshalling an inline keyboard failed: %v", err)
	}

	if got := string(encoded); !strings.Contains(got, `"force_reply":true`) {
		t.Errorf("an inline keyboard that forces a reply must say so\ngot: %s", got)
	}

	encoded, err = json.Marshal(InlineKeyboardMarkup{})
	if err != nil {
		t.Fatalf("marshalling an empty inline keyboard failed: %v", err)
	}

	if got := string(encoded); strings.Contains(got, "force_reply") {
		t.Errorf("a keyboard that does not force a reply must stay silent about it\ngot: %s", got)
	}
}

func TestReplyKeyboardMarkupCarriesForceReply(t *testing.T) {
	markup := ReplyKeyboardMarkup{
		Keyboard:   [][]KeyboardButton{{NewKeyboardButton("Buy")}},
		ForceReply: true,
	}

	encoded, err := json.Marshal(markup)
	if err != nil {
		t.Fatalf("marshalling a reply keyboard failed: %v", err)
	}

	if got := string(encoded); !strings.Contains(got, `"force_reply":true`) {
		t.Errorf("a reply keyboard that forces a reply must say so\ngot: %s", got)
	}
}

// A disabled button is signalled by the presence of an empty object, not by a
// boolean, so an absent field and a present-but-empty one mean opposite things.
func TestInlineKeyboardButtonCarriesADisabledButton(t *testing.T) {
	button := NewInlineKeyboardButtonData("Sold out", "noop")
	button.Disabled = new(DisabledButton)

	encoded, err := json.Marshal(button)
	if err != nil {
		t.Fatalf("marshalling a disabled button failed: %v", err)
	}

	if got := string(encoded); !strings.Contains(got, `"disabled":{}`) {
		t.Errorf("a disabled button must carry an empty object\ngot: %s", got)
	}

	encoded, err = json.Marshal(NewInlineKeyboardButtonData("Buy", "buy:BONK"))
	if err != nil {
		t.Fatalf("marshalling an enabled button failed: %v", err)
	}

	if got := string(encoded); strings.Contains(got, "disabled") {
		t.Errorf("an enabled button must not carry the field at all\ngot: %s", got)
	}
}
