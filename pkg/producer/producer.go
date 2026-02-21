package producer

import (
	"context"
	"encoding/csv"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/itsanindyak/email-campaign/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
)

var (
	meter              = otel.Meter("email-engine")
	loadDuration, _    = meter.Float64Histogram("csv_load_duration_second", metric.WithDescription("Time taken to load csv file"))
	rowsReadTotal, _   = meter.Int64Counter("csv_rows_read_total", metric.WithDescription("Total number of rows read from CSV"))
	rowsFailedTotal, _ = meter.Int64Counter("csv_rows_failed_total", metric.WithDescription("Total number of rows failed to read from CSV"))
)

// LoadFile reads a CSV file containing recipient data and sends each recipient
// to the provided channel for processing. The CSV file should have at least two
// columns: Name and Email. The first row is treated as a header and skipped.
//
// Note: The caller is responsible for closing the channel after LoadFile returns.
func LoadFile(ctx context.Context, path string, ch chan types.Recipient, wg *sync.WaitGroup) error {

	tracer := otel.Tracer("email-engine")

	ctx, span := tracer.Start(ctx, "producer.Loadfile")

	defer span.End()

	span.SetAttributes(attribute.String("file.path", path))

	startTime := time.Now()

	file, err := os.Open(path)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to open file")
		slog.ErrorContext(ctx, "Error opening file", "filepath", path, "Error", err)
		return err
	}

	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read() // skip header
	if err == io.EOF {
		span.SetStatus(codes.Ok, "Empty file")
		return nil
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to read CSV header")
		slog.ErrorContext(ctx, "Error reading CSV header", "error", err)
		return err
	}

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			rowsFailedTotal.Add(ctx, 1)
			slog.ErrorContext(ctx, "Error reading CSV record", "Error", err)
			break
		}

		// Ensure the record has at least name and email columns before accessing.
		if len(record) < 2 {
			rowsFailedTotal.Add(ctx, 1)
			slog.ErrorContext(ctx, "Skipping malformed row", "record", record)
			continue
		}

		wg.Add(1)

		ch <- types.Recipient{
			Name:     record[0],
			Email:    record[1],
			Attempts: 0,
		}

		rowsReadTotal.Add(ctx, 1)
	}

	duration := time.Since(startTime).Seconds()
	loadDuration.Record(ctx, duration)

	span.SetStatus(codes.Ok, "File loaded succesfully")
	slog.InfoContext(ctx, "Finished loading CSV", "path", path, "duration_seconds", duration)

	return nil
}
