package main

import (
    "errors"
    "bytes"
    "io"
    "net/http"
    "fmt"
    "net/http/cookiejar"
    "golang.org/x/net/publicsuffix"
    "log"
    "io/ioutil"
    "net/url"
    "strings"
    "golang.org/x/net/html"
)

func getBody(doc *html.Node) (*html.Node, error) {
    var b *html.Node
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "body" {
            b = n
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)
    if b != nil {
        return b, nil
    }
    return nil, errors.New("Missing <body> in the node tree")
}

func renderNode(n *html.Node) string {
    var buf bytes.Buffer
    w := io.Writer(&buf)
    html.Render(w, n)
    return buf.String()
}

func main() {
    options := cookiejar.Options{
        PublicSuffixList: publicsuffix.List,
    }
    jar, err := cookiejar.New(&options)
    if err != nil {
        log.Fatalln(err)
    }
    client := http.Client{Jar: jar}
    resp, err := client.PostForm("http://www.quakecon.org/wp-content/plugins/age-verification/age-verification.php?redirect_to=http://www.quakecon.org%2F",
        url.Values{
            "age_month": {"01"},
            "age_day": {"01"},
            "age_year": {"1969"},
        })
    if err != nil {
        log.Fatalln(err)
    }

    resp, err = client.Get("http://www.quakecon.org/registration-2/")
    if err != nil {
        log.Fatalln(err)
    }

    data, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
        log.Fatal(err)
    }
    
    htm := string(data)

    doc, _ := html.Parse(strings.NewReader(htm))
    bn, err := getBody(doc)
    if err != nil {
        return
    }
    body := renderNode(bn)
    fmt.Println(body)
}