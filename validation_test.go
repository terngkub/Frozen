package main

import "testing"

func TestNickname(t *testing.T) {
	t.Run("valid_00", nicknameAlpha)
	t.Run("valid_01", nicknameNumber)
	t.Run("valid_02", nicknameSymbol)
	t.Run("valid_03", nicknameCombine)
	t.Run("fail_00", nicknameTooLong)
	t.Run("fail_01", nicknameForbidden)
}

func nicknameAlpha(t *testing.T) {
	nick := "Nick"
	if !isValidNickname(nick) {
		t.Errorf("'%s' should be valid", nick)
	}
}

func nicknameNumber(t *testing.T) {
	nick := "0123"
	if !isValidNickname(nick) {
		t.Errorf("'%s' should be valid", nick)
	}
}

func nicknameSymbol(t *testing.T) {
	nick := "-[]\\`^{}"
	if !isValidNickname(nick) {
		t.Errorf("'%s' should be valid", nick)
	}
}

func nicknameCombine(t *testing.T) {
	nick := "[N0]{n\\}"
	if !isValidNickname(nick) {
		t.Errorf("'%s' should be valid", nick)
	}
}

func nicknameTooLong(t *testing.T) {
	nick := "ThisIsMyNickName"
	if isValidNickname(nick) {
		t.Errorf("'%s' shouldn't be valid", nick)
	}
}

func nicknameForbidden(t *testing.T) {
	nick := "OMG!!!"
	if isValidNickname(nick) {
		t.Errorf("'%s' shouldn't be valid", nick)
	}
}
