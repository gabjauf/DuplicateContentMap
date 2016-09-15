package main

import (
    "os"
    "flag"
	"fmt"
    "time"
	"log"
    "strings"
    "math/rand"
	"github.com/PuerkitoBio/goquery"
    "encoding/csv"
    "strconv"
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

func ExportFormat(pages []page) [][]string {
    result := make([][]string, len(pages) + 1)
    for url := range pages {
        result[url] = make([]string, len(pages))
        result[0][url] = pages[url].url
    }
    result[len(result) - 1] = make([]string, len(pages))
    for i := range pages {
        for j := range pages {
            result[i + 1][j] = strconv.FormatFloat(evaluate(pages[i].shingles, pages[j].shingles), 'f', -1, 64)
        }
    }
    return result
}



func main() {
    rand.Seed(time.Now().UTC().UnixNano()) // initiate the seed for the random function
    k := flag.Int("k", 3, "The value of k, which defines the length of the k-shingles")
    exportFormat := flag.String("ExportFormat", "cmd", "Defines the type of export format expected. Choice between cmd, csv and pdf")
    flag.Parse()
    args := flag.Args()
    pages := make([]page, len(args))

    for url := range args {
        fmt.Printf("=================== FETCHING %s  =================\n", args[url])
        pages[url].url = args[url]
        pages[url].shingles = pages[url].Scrape(args[url], *k)
        //pages[url].PrettyPrint()
        fmt.Println("")
        fmt.Printf("=================== END of url %s ================\n", args[url])
        fmt.Println("")
    }
    matrix := ExportFormat(pages)
    switch *exportFormat {
    case "csv":
        file, err := os.Create("DuplicateContentMatrix.csv")
        if err != nil {
            log.Fatal("Cannot create CSV file", err)
        }
        defer file.Close()
        writer := csv.NewWriter(file)
        for _, value := range matrix {
            if err := writer.Write(value); err != nil {
                log.Fatal("Error during writing", err)
            }
        }
        defer writer.Flush()

    case "pdf":
        fmt.Println("TO DO")

    default:
        for i := range matrix {
            for j := range matrix[i] {
                fmt.Printf("%s\t", matrix[i][j])
            }
            fmt.Printf("\n")
        }
    }
}

