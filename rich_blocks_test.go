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

// Bot API 10.3 put buttons inside the message body. A button block is the only
// way to offer a control that scrolls with the content instead of hanging off
// the bottom of the message, so its wire shape is pinned for the same reason as
// the blocks above.
func TestInputRichBlockButtonsMarshalToTheWireShape(t *testing.T) {
	buttons := InputRichBlockButtons{
		Type: "buttons",
		Buttons: []RichMessageButton{
			{Text: RichTextOf("Buy"), Style: "primary", CallbackData: "buy:BONK"},
			{Text: RichTextOf("Sold out"), Disabled: new(DisabledButton)},
		},
		Align: "center",
	}

	encoded, err := json.Marshal(buttons)
	if err != nil {
		t.Fatalf("marshalling a button row failed: %v", err)
	}

	got := string(encoded)
	for _, want := range []string{
		`"type":"buttons"`,
		`"buttons":[`,
		`{"text":"Buy","style":"primary","callback_data":"buy:BONK"}`,
		`{"text":"Sold out","disabled":{}}`,
		`"align":"center"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("encoded button row is missing %s\ngot: %s", want, got)
		}
	}
}

// An empty switch_inline_query means "insert just the bot's username", which is
// a different instruction from not switching at all. The field is a pointer so
// the two stay apart on the wire.
func TestRichMessageButtonKeepsAnEmptySwitchInlineQuery(t *testing.T) {
	empty := ""

	encoded, err := json.Marshal(RichMessageButton{Text: RichTextOf("Share"), SwitchInlineQuery: &empty})
	if err != nil {
		t.Fatalf("marshalling a switch button failed: %v", err)
	}

	if got := string(encoded); !strings.Contains(got, `"switch_inline_query":""`) {
		t.Errorf("an empty switch query must survive\ngot: %s", got)
	}

	encoded, err = json.Marshal(RichMessageButton{Text: RichTextOf("Share")})
	if err != nil {
		t.Fatalf("marshalling a plain button failed: %v", err)
	}

	if got := string(encoded); strings.Contains(got, "switch_inline_query") {
		t.Errorf("an unset switch query must stay absent\ngot: %s", got)
	}
}

// A button is also a text node, so it can sit inside a run of text rather than
// only in a row of its own.
func TestRichTextButtonMarshalsAsATextNode(t *testing.T) {
	node := RichTextButton{
		Type:   "button",
		Button: RichMessageButton{Text: RichTextOf("Buy"), CallbackData: "buy:BONK"},
	}

	want := `{"type":"button","button":{"text":"Buy","callback_data":"buy:BONK"}}`
	if got := string(RichTextOf(node)); got != want {
		t.Errorf("RichTextOf(button) = %s, want %s", got, want)
	}
}

func TestInputRichBlockDocumentCarriesTheDocumentMedia(t *testing.T) {
	block := InputRichBlockDocument{
		Type: "document",
		Document: InputMediaDocument{
			BaseInputMedia: BaseInputMedia{Type: "document", Media: FileID("BQADAgADOQADjMcoCcioX1GrDvp3Ag")},
		},
	}

	encoded, err := json.Marshal(block)
	if err != nil {
		t.Fatalf("marshalling a document block failed: %v", err)
	}

	got := string(encoded)
	for _, want := range []string{
		`"type":"document"`,
		`"document":{"type":"document","media":"BQADAgADOQADjMcoCcioX1GrDvp3Ag"}`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("encoded document block is missing %s\ngot: %s", want, got)
		}
	}
}

// The collapsible quotation carries flat text where the plain blockquote nests
// blocks. Sending one shaped like the other loses the whole quotation.
func TestInputRichBlockExpandableBlockQuotationCarriesFlatText(t *testing.T) {
	block := InputRichBlockExpandableBlockQuotation{
		Type:   "expandable_blockquote",
		Text:   RichTextOf("Simulated 12 routes"),
		Credit: RichTextOf("router"),
	}

	encoded, err := json.Marshal(block)
	if err != nil {
		t.Fatalf("marshalling an expandable quotation failed: %v", err)
	}

	got := string(encoded)
	want := `{"type":"expandable_blockquote","text":"Simulated 12 routes","credit":"router"}`
	if got != want {
		t.Errorf("encoded expandable quotation = %s, want %s", got, want)
	}
}

func TestInputRichBlockTableMarksACompactTable(t *testing.T) {
	compact := InputRichBlockTable{Type: "table", IsCompact: true}

	encoded, err := json.Marshal(compact)
	if err != nil {
		t.Fatalf("marshalling a compact table failed: %v", err)
	}

	if got := string(encoded); !strings.Contains(got, `"is_compact":true`) {
		t.Errorf("a compact table must say so\ngot: %s", got)
	}

	encoded, err = json.Marshal(InputRichBlockTable{Type: "table"})
	if err != nil {
		t.Fatalf("marshalling a table failed: %v", err)
	}

	if got := string(encoded); strings.Contains(got, "is_compact") {
		t.Errorf("a table that is not compact must stay silent about it\ngot: %s", got)
	}
}

// Received blocks arrive as raw JSON discriminated by "type", so the decode
// side needs the same field names as the send side.
func TestReceivedBotAPI103BlocksUnmarshal(t *testing.T) {
	t.Run("buttons", func(t *testing.T) {
		var block RichBlockButtons
		if err := json.Unmarshal([]byte(`{"type":"buttons","buttons":[{"text":"Buy","callback_data":"buy:BONK","disabled":{}}],"align":"right"}`), &block); err != nil {
			t.Fatalf("decoding a button row failed: %v", err)
		}

		if len(block.Buttons) != 1 {
			t.Fatalf("decoded %d buttons, want 1", len(block.Buttons))
		}
		if got := string(block.Buttons[0].Text); got != `"Buy"` {
			t.Errorf("button text = %s, want \"Buy\"", got)
		}
		if block.Buttons[0].CallbackData != "buy:BONK" {
			t.Errorf("callback data = %q, want buy:BONK", block.Buttons[0].CallbackData)
		}
		if block.Buttons[0].Disabled == nil {
			t.Error("a disabled button must decode to a non-nil DisabledButton")
		}
		if block.Align != "right" {
			t.Errorf("align = %q, want right", block.Align)
		}
	})

	t.Run("document", func(t *testing.T) {
		var block RichBlockDocument
		if err := json.Unmarshal([]byte(`{"type":"document","document":{"file_id":"BQADAgADOQAD","file_unique_id":"AgADOQAD","file_name":"report.pdf"}}`), &block); err != nil {
			t.Fatalf("decoding a document block failed: %v", err)
		}

		if block.Document.FileID != "BQADAgADOQAD" || block.Document.FileName != "report.pdf" {
			t.Errorf("decoded document = %+v", block.Document)
		}
	})

	t.Run("expandable blockquote", func(t *testing.T) {
		var block RichBlockExpandableBlockQuotation
		if err := json.Unmarshal([]byte(`{"type":"expandable_blockquote","text":"Simulated 12 routes"}`), &block); err != nil {
			t.Fatalf("decoding an expandable quotation failed: %v", err)
		}

		if got := string(block.Text); got != `"Simulated 12 routes"` {
			t.Errorf("quotation text = %s", got)
		}
	})

	t.Run("compact table", func(t *testing.T) {
		var block RichBlockTable
		if err := json.Unmarshal([]byte(`{"type":"table","cells":[],"is_compact":true}`), &block); err != nil {
			t.Fatalf("decoding a compact table failed: %v", err)
		}

		if !block.IsCompact {
			t.Error("is_compact must decode to IsCompact")
		}
	})
}
