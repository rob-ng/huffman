package huffman

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

// TestWriteRead tests that writing then reading returns the initial input.
func TestWriteRead(t *testing.T) {
	writeReadTests := []struct {
		input string
	}{
		// Latin characters
		{input: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur metus elit, varius pellentesque porta sit amet, malesuada vel metus. Donec laoreet maximus enim, vel semper augue vehicula sed. Etiam vel ex in velit tempus euismod a non libero. Aenean aliquam sodales metus eget convallis. Curabitur sit amet eleifend dui, a tristique est. Vestibulum vehicula a tortor id ultricies. Nunc quis justo nec dui aliquet convallis. Curabitur condimentum eu massa eu efficitur. Phasellus lacinia nibh metus, vitae faucibus augue hendrerit vitae. Maecenas at tortor eu enim gravida pharetra a a arcu. Mauris sem orci, rhoncus sed cursus non, posuere vehicula nunc. Donec quam nulla, elementum sed rutrum quis, auctor sit amet nisi."},
		// Non-Latin characters
		{input: "優表然択間文井著謙然作連答言遣率。日愛三郡使愛設卒内薄予委禁違掛。点円紙帰局都海以市動注千津堀。意質会再山控柴変無提金燃型光老。表生覚優取掲裁大終図大従理招権追仕。報漢禁質防転出演車価半場紀津止地関聞周。慣宗後不投大競要半読動息中地式。世逆量乱読反群身環無賀予壇。種重道畔掲比罪仕謙出刊終著明月改装社件況。"},
		// Input with new lines
		{input: "This strings consists\nof multiple lines.\nNew lines are\nOK!"},
		// Input with carriage returns
	}

	for _, test := range writeReadTests {
		var wrote bytes.Buffer
		r := strings.NewReader(test.input)

		header, err := NewHeaderFromReader(r)
		if err != nil {
			t.Errorf("Header.NewHeaderFromReader: %v", err)
		}

		hw := NewWriter(&wrote, header)

		_, err = hw.Write([]byte(test.input))
		if err != nil {
			t.Errorf("Header.Write: %v", err)
		}
		_, err = hw.Flush()
		if err != nil {
			t.Errorf("Header.Flush: %v", err)
		}

		read := make([]byte, len(test.input))
		br := bufio.NewReader(&wrote)
		hr := NewReader(br)

		// Once errors are added, add check to ensure Read does not
		// return error.
		hr.Read(read)

		if strings.Compare(test.input, string(read)) != 0 {
			t.Errorf("Input and final result do not match. Input: %s\nResult: %s\n",
				test.input, read)
		}
	}
}

// TestEmpty tests that empty input still results in valid Writer.
func TestEmpty(t *testing.T) {
	emptyInput := ""
	r := strings.NewReader(emptyInput)

	header, err := NewHeaderFromReader(r)
	if err != nil {
		t.Errorf("Header.NewHeaderFromReader: %v", err)
	}

	var wrote bytes.Buffer
	hw := NewWriter(&wrote, header)

	_, err = hw.Write([]byte(emptyInput))
	if err != nil {
		t.Errorf("Header.Write: %v", err)
	}
	_, err = hw.Flush()
	if err != nil {
		t.Errorf("Header.Flush: %v", err)
	}
}
