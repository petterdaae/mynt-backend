package sbanken

type Resource struct {
	sub string
}

func Configure(sub string) Resource {
	return Resource{
		sub: sub,
	}
}
