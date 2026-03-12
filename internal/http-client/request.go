package httpclient

import (
	"fmt"
	"reflect"
)

type Parsable interface {
	Parse() string
}

type QueryParams struct {
	IdList []int `query:"id_list" queryschema:"list containing the article id u wanna search"`
	Start int `query:"start" queryschema:"apply paging defines the index of the first returned result"`
	MaxResults int `query:"max_results" queryschema:"apply paging being the max number of results returned by the query"`
	Search SearchQuery `query:"search_query" queryschema:"apply search filter to results"`
}

func (q *QueryParams) Parse() string {
	var query string = "";
	var first bool = true;

	v := reflect.ValueOf(q).Elem();
	t := reflect.TypeOf(*q);

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i);
		value := v.FieldByName(field.Name);

		key := field.Tag.Get("query");
		
		if value.IsZero() {  continue } 

		// IMPROVE TYPE ASSERTION
		var val string;
		parsableType := (*Parsable)(nil);
		
		if reflect.PointerTo(field.Type).Implements(reflect.TypeOf(parsableType).Elem()) {
			parsable := value.Addr().Interface().(Parsable)
			val = parsable.Parse();
		} else {
			switch field.Type.Kind(){
			case reflect.Int:
				val = fmt.Sprintf("%d", value.Interface().(int))
			case reflect.String:
				val = value.Interface().(string)
			case reflect.Slice: // IMPROVE TYPE ASSERTION
				val = fmt.Sprintf("%v", value.Interface()) 
			}
		}

		if first {
			query = fmt.Sprintf("%s=%s", key, val)
		}else {
			query = fmt.Sprintf("%s&%s=%s", query, key, val)
		}

		first = false;

	}

	return query;
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

func (q *SearchQuery) Parse() string {
	var query string = "";
	var first bool = true;

	v := reflect.ValueOf(*q);
	t := reflect.TypeOf(*q);

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i);
		value := v.FieldByName(field.Name);
		
		key := field.Tag.Get("query");
		val := value.Interface().(string);
		if value.IsZero() { continue }

		if first {
			query = fmt.Sprintf("%s:%s", key, val)
		} else { 
			query = fmt.Sprintf("%s+AND+%s:%s", query, key, val)
		}

		first = false;
	}

	return query;
}