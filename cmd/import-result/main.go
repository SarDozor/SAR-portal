package importresult

import (
	"errors"
	"strings"
	"fmt"
	"flag"
	"os"
	"context"
)

type flags struct {
	path string
	dsn string
	isDryRun bool
}	

type result struct {
	totalCount int
	countSuccess int
	countSkip int
	wrongs []string
}

type emptyStruct struct{}

type validationError struct {
	errors []string
}

func (v validationError) Error() string{
	return strings.Join(v.errors, ", ")
}

const position = map[string]int{
	"JobType": 0, 
	"ReferenceGranule": 1, 
	"SecondaryGranule":2, 
	"ReferenceDate":3, 
	"SecondaryDate":4, 
	"Status":5,
	"Latitude":6,
	"Longitude":7,
	"DisplacementMM":8,
	"Coherence":9,
}

func main(){
	ctx := context.Background()
	var result result
	flags, err := parsingFlags()
	if err != nil{
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, errors.New("Try 'main --help' for more information"))
		os.Exit(1)
	}
	file, err := os.OpenFile(flags.path, os.O_RDONLY, 0766)
	if err != nil{
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	defer file.Close()
	emptyStruct, err := scanSCV()
	if err != nil{
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
	printReport(result)
	if !flags.isDryRun {
		// if flags.dsn == "" {
		// 	if flags.dsn, err = searchDSN(); err != nil {
		// 		fmt.Fprintln(os.Stderr, err)
		// 		os.Exit(4)
		// 	}
		// }
		// pool, err := db.New(ctx, flags.dsn)
		// if err != nil {
		// 	fmt.Fprintln(os.Stderr, err)
		// 	os.Exit(5)
		// }
		// defer pool.Close()
		// countInsertedObjects, err := postListObject(ctx, pool, listObject)
		// if err != nil {
		// 	fmt.Fprintln(os.Stderr, err)
		// 	os.Exit(6)
		// }
		// fmt.Println("\ninserted:", countInsertedObjects)
		fmt.Println("query")
	}
}

func parsingFlags() (flags, error){
	flags := flags{}
	fileFlag := flag.String("file", "", "file name")
	dsnFlag := flag.String("dsn", "", "database url")
	dryRunFlag := flag.Bool("dry-run", false, "test")
	flag.Parse()
	if *fileFlag == ""{
		return flags, fmt.Errorf("flag --file missing")
	}
	flags.path = *fileFlag
	flags.dsn = *dsnFlag
	flags.isDryRun = *dryRunFlag
	return  flags, nil
}

func scanCSV(file *os.File, result *result) (emptyStruct, error) {
	emptyStruct := emptyStruct{}
	fmt.Println("scanning csv file")
	return emptyStruct, nil
}

func printReport(result result) {
	fmt.Println("print the result")
}
