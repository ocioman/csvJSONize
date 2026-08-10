package converter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type conversionError struct {
	defaultMess string
}

func (ce conversionError) Error() string {
	return fmt.Sprintf("error during conversion:")
}

func Serialize(is io.Reader, os io.Writer) error {
	var ce conversionError
	csvReader := csv.NewReader(is)
	builder := strings.Builder{}
	encoder := json.NewEncoder(&builder)

	//leggo gli headers del csv
	headers, headersErr := csvReader.Read()

	if headersErr == io.EOF {
		return fmt.Errorf("%s stream ended unexpectedly", ce.Error())
	} else if headersErr != nil {
		return fmt.Errorf("%s %w", ce.Error(), headersErr)
	}

	builder.WriteRune('[')
	builder.WriteRune('\n')

	for {
		data := make(map[string]string)

		record, recordErr := csvReader.Read()

		if recordErr == io.EOF {
			break
		} else if recordErr != nil {
			return fmt.Errorf("%s %w", ce.Error(), recordErr)
		}

		//json non accetta il trailing comma, quindi scrivo ',' in tutti i record tranne all'inizio del primo e la fine dell'ultimo
		if len([]rune(builder.String())) > 2 {
			builder.WriteRune(',')
		}

		if len(record) < len(headers) {
			return fmt.Errorf("%s missing field", ce.Error())
		}

		for i, v := range record {
			if len([]rune(v)) == 0 {
				data[headers[i]] = "null"
			} else {
				data[headers[i]] = v
			}
		}

		encoderErr := encoder.Encode(data)

		if encoderErr != nil {
			return fmt.Errorf("%s %w", ce.Error(), encoderErr)
		}
	}

	builder.WriteRune(']')

	_, writeErr := os.Write([]byte(builder.String()))

	if writeErr != nil {
		return fmt.Errorf("%s %w", ce.Error(), writeErr)
	}

	return nil
}
