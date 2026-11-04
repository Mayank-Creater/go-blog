package content

type PostMeta struct{
	Title string
	Date string
	Tags []string
	Draft bool
}

type Content struct{
	Metadata PostMeta
	Content string
}