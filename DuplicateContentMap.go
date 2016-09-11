package main

import (
    "strconv"
	"fmt"
    "time"
	"log"
    "strings"
    "math/rand"
    "os"
	"github.com/PuerkitoBio/goquery"
)

type page struct {
    url             string
    shingles        [][]string
}

// Scrape content in url "url" with defined shingle length k
func (page *page) Scrape(url string, k int) {
    doc, err := goquery.NewDocument(url) 
    if err != nil {
        log.Fatal(err)
    }

    // Find the review items
    doc.Each(func(i int, s *goquery.Selection) {
        // Get only the paragraphs "p"
        content := s.Find("p").Text()
        tokens := strings.Fields(content)
        result := make([][]string, 200)
        for shingle := range result { 
            result[shingle] = make([]string, k)
            index := rand.Intn(len(tokens) - k)
            for i := 0; i < k; i++ {
                result[shingle][i] = tokens[index + i]
            }
        }
        page.shingles = result
    })

}

func (page *page) PrettyPrint() {
    for shingle := range page.shingles {
        fmt.Println(page.shingles[shingle])
    }
}

func evaluate(page1 page, page2 page) int {
    result := 0
    for shingle1 := range page1.shingles {
        for shingle2 := range page2.shingles {
            if StringSliceCompare(page1.shingles[shingle1], page2.shingles[shingle2]) == true {
                result += 1
                //fmt.Println(page1.shingles[shingle1], page2.shingles[shingle2])
            }
        }
    }
    return result
    
}

// Compare two slices of string of same size same index
func StringSliceCompare(slice1 []string, slice2 []string) bool {
    for i := range slice1 {
        if strings.Compare(slice1[i], slice2[i]) != 0 {
            return false
        }
    }
    return true
}

// compare two slices of string of same size but full index
func StringSliceCompareFull(slice1 []string, slice2 []string) bool {
    for i := range slice1 {
        for j := range slice2 {
            if strings.Compare(slice1[i], slice2[j]) != 0 {
                return false
            }
        }
    }
    return true
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())
    args := os.Args[2:]
    k, err := strconv.Atoi(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    pages := make([]page, len(args))

    for url := range args {
        fmt.Printf("=================== FETCHING %s  =================\n", args[url])
        pages[url].url = args[url]
        pages[url].Scrape(args[url], k)
        //pages[url].PrettyPrint()
        fmt.Println("")
        fmt.Printf("========== END of url %s ===========\n", args[url])
        fmt.Println("")
    }
    fmt.Println("duplicate content:", evaluate(pages[0], pages[1]))
}
