package email

import (
	"strings"
	"testing"
)

func TestConvertRunes(t *testing.T) {
	var testData = map[string]string{
		"=??=_.": "=3D=3F=3F=3D=5F.",
		"Příšerně žluťoučký kůn úpěl ďábelské ódy 🐎": "P=C5=99=C3=AD=C5=A1ern=C4=9B_=C5=BElu=C5=A5ou=C4=8Dk=C3=BD_k=C5=AFn_=C3=BAp=C4=9Bl_=C4=8F=C3=A1belsk=C3=A9_=C3=B3dy_=F0=9F=90=8E",
	}
	for input, expected := range testData {
		got := strings.Join(convertRunes(input), "")
		if got != expected {
			t.Errorf("Input: '%s', expected '%s', got: '%s'", input, expected, got)
		}
	}
}

type genHeaderTestData struct {
	name     string
	value    string
	expected string
	maxWidth int
}

func TestGenHeaderQ(t *testing.T) {
	var testData = []genHeaderTestData{
		{
			name:  "Subject",
			value: "Příšerně žluťoučký kůn úpěl ďábelské ódy 🐎",
			expected: "Subject: =?utf-8?q?P=C5=99=C3=AD=C5=A1ern=C4=9B_=C5=BElu=C5=A5ou=C4=8Dk?=\n" +
				"    =?utf-8?q?=C3=BD_k=C5=AFn_=C3=BAp=C4=9Bl_=C4=8F=C3=A1belsk=C3=A9_=C3=B3?=\n" +
				"    =?utf-8?q?dy_=F0=9F=90=8E?=",
			maxWidth: 80,
		},
	}
	for _, data := range testData {
		got := genHeader(data.name, data.value, data.maxWidth)
		if got != data.expected {
			t.Errorf("Input: '%s', expected \n===\n%s\n===, got: \n===\n%s\n==='", data.value, data.expected, got)
		}

	}
}

type genAddressHeaderTestData struct {
	name      string
	addresses []Address
	expected  string
	maxLength int
}

func TestGenAddressHeader(t *testing.T) {
	var testData = []genAddressHeaderTestData{
		{
			name: "To",
			addresses: []Address{
				{
					Name:  "Oldřich Jánský",
					Email: "olrd@example.com",
				},
			},
			expected:  "To: =?utf-8?q?Old=C5=99ich_J=C3=A1nsk=C3=BD?= <olrd@example.com>",
			maxLength: 80,
		},
		{
			name: "Subject",
			addresses: []Address{
				{
					Name:  "Oldřich Jánský",
					Email: "olrd@example.com",
				},
				{
					Name:  "Jan Novák",
					Email: "novak@example.com",
				},
			},
			expected: "Subject: =?utf-8?q?Old=C5=99ich_J=C3=A1nsk=C3=BD?= <olrd@example.com>, \n" +
				"    =?utf-8?q?Jan_Nov=C3=A1k?= <novak@example.com>",
			maxLength: 80,
		},
	}
	for _, data := range testData {
		got := genAddressHeader(data.name, data.addresses, data.maxLength)
		if got != data.expected {
			t.Errorf("Test: '%s', expected \n===\n%s\n===, got: \n===\n%s\n==='", data.name, data.expected, got)
		}

	}
}
