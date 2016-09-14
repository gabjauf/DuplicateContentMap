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
func (page *page) Scrape(url string, k int) [][]string {
    doc, err := goquery.NewDocument(url) 
    if err != nil {
        log.Fatal(err)
    }

    result := make([][]string, 200)
    // Find the review items
    doc.Each(func(i int, s *goquery.Selection) {
        // Get only the paragraphs "p"
        content := s.Find("p").Text()
        tokens := strings.Fields(content)
        for shingle := range result { 
            result[shingle] = make([]string, k)
            index := rand.Intn(len(tokens) - k)
            for i := 0; i < k; i++ {
                result[shingle][i] = tokens[index + i]
            }
        }
        page.shingles = result
    })
    return result

}

func (page *page) PrettyPrint() {
    for shingle := range page.shingles {
        fmt.Println(page.shingles[shingle])
    }
}

func evaluate(shingles1 [][]string, shingles2 [][]string) float64 {
    intersect := make([][]string, 0)
    union := make([][]string, 0)
    for shingle1 := range shingles1 {
        for shingle2 := range shingles2 {
            if StringSliceCompare(shingles1[shingle1], shingles2[shingle2]) == true {
                intersect = AppendIfMissing(intersect, shingles1[shingle1])
                union = AppendIfMissing(union, shingles1[shingle1])
                //fmt.Println(page1.shingles[shingle1], page2.shingles[shingle2])
                break
            } else {
                union = AppendIfMissing(union, shingles1[shingle1])
               union = AppendIfMissing(union, shingles2[shingle2])
            }
        }
    }
    //fmt.Println()
    //fmt.Println("Union",len(union))
    //fmt.Println("intersect", len(intersect))
    return float64(len(intersect)) / float64(len(union))
    
}

func AppendIfMissing(slice [][]string, elt []string) [][]string {
    for ele := range slice {
        if StringSliceCompare(slice[ele], elt) {
            return slice
        }
    }
    return append(slice, elt)
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
        pages[url].shingles = pages[url].Scrape(args[url], k)
        //pages[url].PrettyPrint()
        fmt.Println("")
        fmt.Printf("========== END of url %s ===========\n", args[url])
        fmt.Println("")
    }
    
    buf := make([]float64, len(pages))
    for i := 0; i < len(pages); i ++ {
        //pages[0].PrettyPrint()

        //pages[1].PrettyPrint()
        for j := 0; j < len(pages); j++ {

            //pages[0].PrettyPrint()
            buf[j] = evaluate(pages[i].shingles, pages[j].shingles)
            //pages[0].PrettyPrint()
            fmt.Printf("%f ", buf[j])
        }
        fmt.Printf("\n")
    }
}
