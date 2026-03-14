package httpclient

import (
	"fmt"
	"net/http"

	atomparser "github.com/wbernest/atom-parser"
	"golang.org/x/tools/blog/atom"
)

func ParseAtom(res *http.Response) (*atom.Feed, error){
	data, err := ReadBody(res);
	if err != nil {
		return nil, fmt.Errorf("error reading body response: %s", err)
	}

	feed, err := atomparser.ParseString(string(data));
	if err != nil {
		return nil, fmt.Errorf("error parsing atom xml: %s", err)
	}

	return feed, nil
}