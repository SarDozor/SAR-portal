package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/darkn3sssss/sar-dozor/internal/db"
	"github.com/darkn3sssss/sar-dozor/internal/model"
	"github.com/joho/godotenv"
)

type flags struct {
	path     string
	dsn      string
	isDryRun bool
}

type result struct {
	totalCount   int
	countSuccess int
	countSkip    int
	wrongs       []string
}

type validationError struct {
	errors []string
}

func (v validationError) Error() string {
	return strings.Join(v.errors, ", ")
}

const (
	objectCode = iota
	pairName
	referenceDate
	secondaryDate
	latitude
	longitude
	damCoherenceMean
	damCoherenceValidPixels
	referenceMethod
	displacementMM
	referenceGranule
	secondaryGranule
)

var envParameter string = "DATABASE_URL"

func main() {
	ctx := context.Background()
	var result result
	flags, err := parsingFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, errors.New("Try 'main --help' for more information"))
		os.Exit(1)
	}
	file, err := os.OpenFile(flags.path, os.O_RDONLY, 0766)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	defer file.Close()
	//listResult, err := scanSCV(file, &result)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
	printReport(result)
	if !flags.isDryRun {
		if flags.dsn == "" {
			if flags.dsn, err = searchDSN(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(4)
			}
		}
		pool, err := db.New(ctx, flags.dsn)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(5)
		}
		defer pool.Close()
		/*countInsertedResults, err := postListResult(ctx, pool, listResult)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(6)
		}
		fmt.Println("\ninserted:", countInsertedResults)*/
		fmt.Println("query")
	}
}

func parsingFlags() (flags, error) {
	flags := flags{}
	fileFlag := flag.String("file", "", "file name")
	dsnFlag := flag.String("dsn", "", "database url")
	dryRunFlag := flag.Bool("dry-run", false, "test")
	flag.Parse()
	if *fileFlag == "" {
		return flags, fmt.Errorf("flag --file missing")
	}
	flags.path = *fileFlag
	flags.dsn = *dsnFlag
	flags.isDryRun = *dryRunFlag
	return flags, nil
}

func searchDSN() (string, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "cannot load .env: %s", err.Error())
	}
	dsn := os.Getenv(envParameter)
	if dsn == "" {
		return "", errors.New("variable database url is not set")
	}
	return dsn, nil
}

func scanSCV(file *os.File, result *result) ([]model.ResultRow, error) {
	listResult := []model.ResultRow{}
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 12
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("cannot read: %w", err)
	}
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		var parseErr *csv.ParseError
		if errors.As(err, &parseErr) {
			result.totalCount++
			result.countSkip++
			result.wrongs = append(result.wrongs, fmt.Sprintf("[%d] %v", parseErr.Line, parseErr))
			continue
		}
		if err != nil {
			return nil, err
		}
		if len(record) == 0 {
			continue
		}
		line, _ := reader.FieldPos(0)
		resultRow, err := parseResult(record)
		if err != nil {
			result.countSkip++
			result.wrongs = append(result.wrongs, fmt.Sprintf("[%d] %v", line, err))
		} else {
			listResult = append(listResult, resultRow)
			result.countSuccess++
		}
		result.totalCount++
	}
	return listResult, nil
}

func parseResult(record []string) (model.ResultRow, error) {
	var newResult model.ResultRow
	var errs []string
	for i := range record {
		record[i] = strings.TrimSpace(record[i])
	}
	record[objectCode] = strings.Join(strings.Fields(record[objectCode]), " ")
	if record[objectCode] == "" {
		errs = append(errs, "object_code: text missing")
	} else {
		newResult.ObjectCode = record[objectCode]
	}
	if record[pairName] == "" {
		errs = append(errs, "pair_name: text missing")
	} else {
		newResult.PairName = record[pairName]
	}
	if record[referenceDate] == "" {
		errs = append(errs, "reference_date: text missing")
	} else {
		time, err := time.Parse("2006-01-02", record[referenceDate])
		if err != nil {
			errs = append(errs, "reference_date: wrong time")
		}
		newResult.ReferenceDate = time
	}
	if record[secondaryDate] == "" {
		errs = append(errs, "secondary_date: text missing")
	} else {
		time, err := time.Parse("2006-01-02", record[secondaryDate])
		if err != nil {
			errs = append(errs, "secondary_date: wrong time")
		}
		if time.After(newResult.ReferenceDate) {
			newResult.SecondaryDate = time
		} else {
			errs = append(errs, "secondary_date: earlier reference_date")
		}
	}
	if record[referenceGranule] == "" {
		errs = append(errs, "reference_granule: text missing")
	} else {
		newResult.ReferenceGranule = record[referenceGranule]
	}
	if record[secondaryGranule] == "" {
		errs = append(errs, "secondary_granule: text missing")
	} else {
		newResult.SecondaryGranule = record[secondaryGranule]
	}
	if latitude, err := strconv.ParseFloat(record[latitude], 64); err != nil {
		errs = append(errs, "latitude: cannot be format to float")
	} else if math.IsNaN(latitude) || math.IsInf(latitude, 0) {
		errs = append(errs, "latitude: must be finite")
	} else if latitude < -90 || latitude > 90 {
		errs = append(errs, "latitude: must be between -90 and 90")
	} else {
		newResult.Latitude = latitude
	}
	if longitude, err := strconv.ParseFloat(record[longitude], 64); err != nil {
		errs = append(errs, "longitude: cannot be format to float")
	} else if math.IsNaN(longitude) || math.IsInf(longitude, 0) {
		errs = append(errs, "longitude: must be finite")
	} else if longitude < -180 || longitude > 180 {
		errs = append(errs, "longitude: must be between -180 and 180")
	} else {
		newResult.Longitude = longitude
	}
	if damCoherenceMean, err := strconv.ParseFloat(record[damCoherenceMean], 64); err != nil {
		errs = append(errs, "dam_coherence_mean: cannot be format to float")
	} else if math.IsNaN(damCoherenceMean) || math.IsInf(damCoherenceMean, 0) {
		errs = append(errs, "dam_coherence_mean: must be finite")
	} else if damCoherenceMean < 0 || damCoherenceMean > 1 {
		errs = append(errs, "dam_coherence_mean: must be between 0 and 1")
	} else {
		newResult.DamCoherenceMean = damCoherenceMean
	}
	if damCoherenceValidPixels, err := strconv.ParseInt(record[damCoherenceValidPixels], 10, 32); err != nil {
		errs = append(errs, "dam_coherence_valid_pixels: cannot be format to int")
	} else if damCoherenceValidPixels < 0 {
		errs = append(errs, "dam_coherence_valid_pixels: cannot be negative")
	} else {
		newResult.DamCoherenceValidPixels = int(damCoherenceValidPixels)
	}
	if record[referenceMethod] == "" {
		errs = append(errs, "reference_method: text missing")
	} else {
		newResult.ReferenceMethod = record[referenceMethod]
	}
	if displacementMM, err := strconv.ParseFloat(record[displacementMM], 64); err != nil {
		errs = append(errs, "displacement_mm: cannot be format to float")
	} else if math.IsNaN(displacementMM) || math.IsInf(displacementMM, 0) {
		errs = append(errs, "displacement_mm: must be finite")
	} else {
		newResult.DisplacementMM = displacementMM
	}
	if len(errs) > 0 {
		return model.ResultRow{}, &validationError{errors: errs}
	}
	return newResult, nil
}

func printReport(result result) {
	for _, wrong := range result.wrongs {
		fmt.Fprintf(os.Stdout, "%s\n", wrong)
	}
	fmt.Fprintf(os.Stdout, "\nlines read: %d\n", result.totalCount)
	fmt.Fprintf(os.Stdout, "ready to download: %d\n", result.countSuccess)
	fmt.Fprintf(os.Stdout, "skipped: %d\n", result.countSkip)
}

/*
func postListResult(ctx context.Context, pool *pgxpool.Pool, listObject []model.ResultRow) (int, error) {
	var countInsertedObjects int
	query := "INSERT INTO objects (name, country, latitude, longitude, status, raise_type, max_height_m) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	for index, obj := range listObject {
		_, err := tx.Exec(
			ctx, query, obj.Name, obj.Country, obj.Latitude, obj.Longitude,
			obj.Status, obj.RaiseType, obj.MaxHeightM,
		)
		if err != nil {
			return 0, fmt.Errorf("[%d] insert %q: %w", index+1, obj.Name, err)
		}
		countInsertedObjects++
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return countInsertedObjects, nil
}*/
