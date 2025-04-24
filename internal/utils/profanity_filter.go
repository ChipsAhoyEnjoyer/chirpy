package utils

import (
	"fmt"
	"strings"

	"github.com/ChipsAhoyEnjoyer/chirpy/pkg/titlecase"
)

func ProfanityFilter(s string) string {
	profanities := []string{"kerfuffle", "sharbert", "fornax"}
	transformed := []string{}
	for _, p := range profanities {
		transformed = append(transformed, titlecase.Titlecase(p))
		transformed = append(transformed, strings.ToUpper(p))
	}
	profanities = append(profanities, transformed...)
	for _, p := range profanities {
		sep := fmt.Sprintf(" %s ", p)
		s = strings.Join(strings.Split(s, sep), " **** ")
	}
	return s
}
