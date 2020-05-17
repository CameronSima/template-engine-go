package main

type Origin struct {
	TemplateName string
	Template     string
}

type Loader interface {
	GetContents()
}
