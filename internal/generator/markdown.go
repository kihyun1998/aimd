package generator

type MarkdownGenerator interface {
	Generate(files []string) error
	SetTemplate(template string) error
}
