package main

type Token struct {
	typ   TokenType
	value string
	line  int
	col   int
}
