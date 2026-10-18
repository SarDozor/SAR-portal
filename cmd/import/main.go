package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/darkn3sssss/sar-dozor/internal/db"
	"github.com/darkn3sssss/sar-dozor/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
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
	tsf = iota
	country
	latitude
	longitude
	status
	raise_type
	current_maximum_heigh
)

var envParameter string = "DATABASE_URL"

func main() {
	ctx := context.Background()
	var result result
	flags, err := parsingFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, errors.New("Try 'main --help' for more information."))
		os.Exit(1)
	}
	file, err := os.OpenFile(flags.path, os.O_RDONLY, 0766)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	defer file.Close()
	listObject, err := scanCsv(file, &result)
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
		countInsertedObjects, err := postListObject(ctx, pool, listObject)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(6)
		}
		fmt.Println("\ninserted:", countInsertedObjects)
	}
}

func parsingFlags() (flags, error) {
	flags := flags{}
	fileFlag := flag.String("file", "", "file name")
	dsnFlag := flag.String("dsn", "", "database url")
	dryRunFlag := flag.Bool("dry-run", false, "final result")
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

func scanCsv(file *os.File, result *result) ([]model.Object, error) {
	listObject := []model.Object{}
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 7
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("cannot read: %w", err)
	}
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		line, _ := reader.FieldPos(0)
		var parseErr *csv.ParseError
		if errors.As(err, &parseErr) {
			result.totalCount++
			result.countSkip++
			result.wrongs = append(result.wrongs, fmt.Sprintf("[%d] %v", line, parseErr))
			continue
		}
		if err != nil {
			return nil, err
		}
		if len(record) == 0 {
			continue
		}
		object, err := parseObject(record)
		if err != nil {
			result.countSkip++
			result.wrongs = append(result.wrongs, fmt.Sprintf("[%d] %v", line, err))
		} else {
			listObject = append(listObject, object)
			result.countSuccess++
		}
		result.totalCount++
	}
	return listObject, nil
}

func parseObject(record []string) (model.Object, error) {
	var newObject model.Object
	var errs []string
	for i := range record {
		record[i] = strings.TrimSpace(record[i])
	}
	record[tsf] = strings.Join(strings.Fields(record[tsf]), " ")
	if record[tsf] == "" {
		errs = append(errs, "name: text missing")
	} else {
		newObject.Name = record[tsf]
	}
	if record[country] == "" {
		errs = append(errs, "country: text missing")
	} else {
		newObject.Country = record[country]
	}
	if latitude, err := strconv.ParseFloat(record[latitude], 64); err != nil {
		errs = append(errs, "latitude: cannot be format to float")
	} else if latitude < -90 || latitude > 90 {
		errs = append(errs, "latitude: out of range")
	} else {
		newObject.Latitude = &latitude
	}
	if longitude, err := strconv.ParseFloat(record[longitude], 64); err != nil {
		errs = append(errs, "longitude: cannot be format to float")
	} else if longitude < -180 || longitude > 180 {
		errs = append(errs, "longitude: out of range")
	} else {
		newObject.Longitude = &longitude
	}
	if record[status] != "" {
		val := record[status]
		newObject.Status = &val
	}
	if record[raise_type] != "" {
		val := record[raise_type]
		newObject.RaiseType = &val
	}
	if record[current_maximum_heigh] != "" {
		if current_maximum_heigh, err := strconv.ParseFloat(record[current_maximum_heigh], 64); err != nil {
			errs = append(errs, "current_maximum_heigh: cannot be format to float")
		} else if current_maximum_heigh < 0 {
			errs = append(errs, "current_maximum_heigh: cannot be negative")
		} else {
			newObject.MaxHeightM = &current_maximum_heigh
		}
	}
	if len(errs) > 0 {
		return model.Object{}, &validationError{errors: errs}
	}
	return newObject, nil
}

func printReport(result result) {
	for _, wrong := range result.wrongs {
		fmt.Fprintf(os.Stdout, "%s\n", wrong)
	}
	fmt.Fprintf(os.Stdout, "\nlines read: %d\n", result.totalCount)
	fmt.Fprintf(os.Stdout, "ready to download: %d\n", result.countSuccess)
	fmt.Fprintf(os.Stdout, "skipped: %d\n", result.countSkip)
}

func postListObject(ctx context.Context, pool *pgxpool.Pool, listObject []model.Object) (int, error) {
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
}
