package httpclient

type QueryParams struct {
	IdList []int `query:"id_list" queryschema:"list containing the article id u wanna search"`
	Start int `query:"start" queryschema:"apply paging defines the index of the first returned result"`
	MaxResults int `query:"max_results" queryschema:"apply paging being the max number of results returned by the query"`
}

type SearchQuery struct {
	Title string `query:"ti" queryschema:"title of the article u wanna search"`
	Author string `query:"au" queryschema:"author of the articles u wanna search"`
	Abstract string `query:"abs" queryschema:"abstract of the article u wanna search"`
	Comment string `query:"co" queryschema:"comment in the article u wanna search"`
	JournalReference string `query:"jr" queryschema:"journal reference u wanna search"`
	SubjectCategory string `query:"cat" queryschema:"subject category u wanna search"`
	ReportNumber string `query:"rn" queryschema:"report number of the article u wanna search"`
	All string `query:"all" queryschema:"all above"`
}