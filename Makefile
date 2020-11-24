spellcheck:
	@find . -type f -name '*.*' | grep -v vendor/ | xargs misspell -error -w