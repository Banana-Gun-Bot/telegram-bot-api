package tgbotapi

import (
	"encoding/json"
	"strings"
	"testing"
)

// The blocks form exists because the html form is lossy: the server maps HTML
// into rich nodes and silently drops tags and attributes it does not model, so
// structure expressed as HTML may or may not survive. A block carries the same
// structure as typed fields. These tests pin the wire shape, since a wrong
// field name fails the same silent way.

func TestInputRichMessageBlocksMarshalToTheWireShape(t *testing.T) {
	message := InputRichMessage{
		Blocks: []InputRichBlock{
			InputRichBlockParagraph{Type: "paragraph", Text: RichTextOf("Positions")},
			InputRichBlockDivider{Type: "divider"},
		},
	}

	encoded, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("marshalling a block message failed: %v", err)
	}

	got := string(encoded)
	for _, want := range []string{
		`"blocks":[`,
		`{"type":"paragraph","text":"Positions"}`,
		`{"type":"divider"}`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("encoded message is missing %s\ngot: %s", want, got)
		}
	}

	// html and markdown stay absent, so a blocks-only message never sends two
	// competing descriptions of its content.
	if strings.Contains(got, `"html"`) || strings.Contains(got, `"markdown"`) {
		t.Errorf("a blocks-only message must not carry html or markdown\ngot: %s", got)
	}
}

func TestInputRichBlockTableKeepsCellSpansAndAlignment(t *testing.T) {
	table := InputRichBlockTable{
		Type: "table",
		Cells: [][]RichBlockTableCell{{
			{Text: RichTextOf("BONK"), Colspan: 2, Align: "left", Valign: "middle"},
			{Text: RichTextOf("12.40"), Colspan: 3, Align: "right", Valign: "middle"},
		}},
	}

	encoded, err := json.Marshal(table)
	if err != nil {
		t.Fatalf("marshalling a table failed: %v", err)
	}

	got := string(encoded)
	for _, want := range []string{
		`"type":"table"`,
		`"cells":[[`,
		`"text":"BONK"`,
		`"colspan":2`,
		`"colspan":3`,
		`"align":"right"`,
		`"valign":"middle"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("encoded table is missing %s\ngot: %s", want, got)
		}
	}
}

// RichText is polymorphic on the wire -- a bare string, an array, or a typed
// node -- so RichTextOf has to preserve whichever form it is handed rather
// than stringifying everything.
func TestRichTextOfPreservesEachForm(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value interface{}
		want  string
	}{
		{name: "plain string", value: "BONK", want: `"BONK"`},
		{
			name:  "typed node",
			value: RichTextUrl{Type: "url", Text: RichTextOf("BONK"), URL: "https://t.me/bot"},
			want:  `{"type":"url","text":"BONK","url":"https://t.me/bot"}`,
		},
		{
			name: "run of nodes",
			value: []interface{}{
				RichTextCustomEmoji{Type: "custom_emoji", CustomEmojiID: "5", AlternativeText: "🅢"},
				"BONK",
			},
			want: `[{"type":"custom_emoji","custom_emoji_id":"5","alternative_text":"🅢"},"BONK"]`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := string(RichTextOf(tc.value)); got != tc.want {
				t.Errorf("RichTextOf(%v) = %s, want %s", tc.value, got, tc.want)
			}
		})
	}
}

// An unmarshallable value yields a JSON null rather than a broken document:
// the message still sends, with that one text absent.
func TestRichTextOfFallsBackToNull(t *testing.T) {
	if got := string(RichTextOf(make(chan int))); got != "null" {
		t.Errorf("RichTextOf(unmarshallable) = %s, want null", got)
	}
}
