package titlecase

import "strings"

func Titlecase(sentence string) string {
	words := strings.Fields(sentence)
	f_sentence := []string{}
	for _, word := range words {
		valOfFirstLetter := word[0]
		theRest := word[1:]

		if 97 <= valOfFirstLetter && valOfFirstLetter <= 122 {
			valOfFirstLetter -= 32
			f_sentence = append(f_sentence, strings.Join(
				[]string{
					string(valOfFirstLetter),
					strings.ToLower(theRest),
				},
				"",
			),
			)

		} else if 65 <= valOfFirstLetter && valOfFirstLetter <= 90 {
			f_sentence = append(f_sentence, strings.Join(
				[]string{
					string(valOfFirstLetter),
					strings.ToLower(theRest),
				},
				"",
			),
			)

		} else { // if the first char ain't a letter, just return the word
			f_sentence = append(f_sentence, word)
		}
	}
	return strings.Join(f_sentence, " ")
}
