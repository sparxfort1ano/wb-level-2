package main

import "testing"

func TestUnpackRune(t *testing.T) {
	testTable := []struct {
		str         string
		expectedStr string
		expectedErr error
	}{
		{str: "а4бс2д5е",
			expectedStr: "аааабссддддде",
			expectedErr: nil},
		{str: "abcd",
			expectedStr: "abcd",
			expectedErr: nil},
		{str: "45",
			expectedStr: "",
			expectedErr: ErrDigitsInvalidFormat},
		{str: "",
			expectedStr: "",
			expectedErr: nil},
		{str: "a10",
			expectedStr: "aaaaaaaaaa",
			expectedErr: nil},
		{str: "a10б4c",
			expectedStr: "aaaaaaaaaaббббc",
			expectedErr: nil},
		{str: "б\\110",
			expectedStr: "б1111111111",
			expectedErr: nil},
		{str: "qwe\\4\\5",
			expectedStr: "qwe45",
			expectedErr: nil},
		{str: "а\\1",
			expectedStr: "а1",
			expectedErr: nil},
		{str: "qwe\\\\10",
			expectedStr: "qwe\\\\\\\\\\\\\\\\\\\\",
			expectedErr: nil},
		{str: "qwe\\\\\\",
			expectedStr: "",
			expectedErr: ErrEscapeInvalidFormat},
		{str: "a1b1c1e",
			expectedStr: "abce",
			expectedErr: nil},
		{str: "ээ4б0к",
			expectedStr: "эээээк",
			expectedErr: nil},
	}

	for _, tt := range testTable {
		gotStr, gotErr := UnpackRune([]rune(tt.str))
		if gotStr != tt.expectedStr || gotErr != tt.expectedErr {
			if gotStr != tt.expectedStr {
				t.Errorf("Unpacking: %s -> %s, expected: %s -> %s",
					tt.str, gotStr, tt.str, tt.expectedStr)
			}
			if gotErr != tt.expectedErr {
				t.Errorf("Got err: %s, expected err: %s",
					tt.expectedErr.Error(), gotErr.Error())
			}
		}
	}
}
