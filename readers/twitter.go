package readers

import (
	"strconv"
	"strings"
	"time"

	"github.com/tyba/srcd-domain/models/social"
	"github.com/tyba/srcd-rovers/http"

	"github.com/PuerkitoBio/goquery"
)

const TwitterBaseURL = "https://twitter.com/%s"

type TwitterReader struct {
	client *http.Client
}

func NewTwitterReader(client *http.Client) *TwitterReader {
	return &TwitterReader{client}
}

func (t *TwitterReader) GetProfileByURL(url string) (*social.TwitterProfile, error) {
	req, err := http.NewRequest(url)
	if err != nil {
		return nil, err
	}

	doc, _, err := t.client.DoHTML(req)
	if err != nil {
		return nil, err
	}

	p := &social.TwitterProfile{Url: url, Created: time.Now()}
	t.fillBasicInfo(doc, p)
	t.fillStats(doc, p)

	return p, nil
}

func (t *TwitterReader) fillBasicInfo(doc *goquery.Document, p *social.TwitterProfile) {
	p.Handle = doc.Find(".ProfileHeaderCard-screenname span").Text()
	p.FullName = doc.Find(".ProfileHeaderCard-name a").Text()
	p.Location = strings.Trim(doc.Find(".ProfileHeaderCard-locationText").Text(), "\n\r\t ")
	p.Bio = doc.Find(".ProfileHeaderCard-bio").Text()
	p.Web, _ = doc.Find(".ProfileHeaderCard-url a").Attr("title")
}

func (t *TwitterReader) fillStats(doc *goquery.Document, p *social.TwitterProfile) {
	tweets := doc.Find("[data-nav='tweets'] .ProfileNav-value").Text()
	if value, err := strconv.Atoi(tweets); err == nil {
		p.Tweets = value
	}

	following := doc.Find("[data-nav='following'] .ProfileNav-value").Text()
	if value, err := strconv.Atoi(following); err == nil {
		p.Following = value
	}

	followers := doc.Find("[data-nav='followers'] .ProfileNav-value").Text()
	if value, err := strconv.Atoi(followers); err == nil {
		p.Followers = value
	}

	favorites := doc.Find("[data-nav='favorites'] .ProfileNav-value").Text()
	if value, err := strconv.Atoi(favorites); err == nil {
		p.Favorites = value
	}
}
