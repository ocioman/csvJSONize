package converter

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var ce conversionError

type conversionError struct{}

func (ce conversionError) Error() string {
	return fmt.Sprintf("error during conversion:")
}

func Serialize(is io.Reader, os io.Writer) (err error) {
	csvReader := csv.NewReader(is)
	bw := bufio.NewWriter(os)

	/*
		eseguo il join in modo che se viene catturato un errore prima del flush e anche il flush causa
		errore (che viene sempre eseguito prima di ogni return a causa del defer) li stampo entrambi
	*/
	defer func() {
		if flushErr := bw.Flush(); flushErr != nil {
			if err != nil {
				err = errors.Join(err, fmt.Errorf("%s %w", ce.Error(), flushErr))
			} else {
				err = fmt.Errorf("%s %w", ce.Error(), flushErr)
			}
		}
	}()

	/*
		uso un bufio in modo da poter scrivere ogni volta la stringa in un buffer
		senza usare il builder -> se ho un file da 2gb lo spezzo in parti da 4kb
		e ogni volta che il buffer si riempie faccio il flush -> in questo modo
		non avrò mai più di 4kb in ram -> se usassi uno string builder terrei 2gb di stringa
		all'interno della ram fino al flush
	*/
	encoder := json.NewEncoder(bw)

	//leggo gli headers del csv
	headers, headersErr := csvReader.Read()

	if headersErr == io.EOF {
		return fmt.Errorf("%s stream ended unexpectedly", ce.Error())
	} else if headersErr != nil {
		return fmt.Errorf("%s %w", ce.Error(), headersErr)
	}

	_, writeOpeningErr := bw.WriteString("[\n")

	if writeOpeningErr != nil {
		return fmt.Errorf("%s %w", ce.Error(), writeOpeningErr)
	}

	firstRecord := true

	for {
		data := make(map[string]string)

		record, recordErr := csvReader.Read()

		if recordErr == io.EOF {
			break
		} else if recordErr != nil {
			return fmt.Errorf("%s %w", ce.Error(), recordErr)
		}

		//json non accetta il trailing comma, quindi scrivo ',' in tutti i record tranne all'inizio del primo e la fine dell'ultimo
		if !firstRecord {
			_, writeCommaErr := bw.WriteRune(',')

			if writeCommaErr != nil {
				return fmt.Errorf("%s %w", ce.Error(), writeCommaErr)
			}
		}

		if len(record) < len(headers) {
			return fmt.Errorf("%s missing field", ce.Error())
		}

		for i, v := range record {
			if v == "" {
				data[headers[i]] = "null"
			} else {
				data[headers[i]] = v
			}
		}

		encoderErr := encoder.Encode(data)

		if encoderErr != nil {
			return fmt.Errorf("%s %w", ce.Error(), encoderErr)
		}

		if firstRecord {
			firstRecord = false
		}
	}

	_, writeClosingErr := bw.WriteRune(']')

	if writeClosingErr != nil {
		return fmt.Errorf("%s %w", ce.Error(), writeClosingErr)
	}

	return err
}

func Deserialize(is io.Reader, os io.Writer) error {
	csvWriter := csv.NewWriter(os)
	decoder := json.NewDecoder(is)
	firstRecord := true

	defer csvWriter.Flush()

	//devo estrarre "[" altrimenti deserializza tutto l'array JSON
	token, delError := decoder.Token()

	if delError == io.EOF {
		return fmt.Errorf("%s stream ended unexpectedly", ce.Error())
	} else if delError != nil {
		return fmt.Errorf("%s %w", ce.Error(), delError)
	}

	del, isToken := token.(json.Delim)

	if isToken {
		if del != json.Delim('[') {
			return fmt.Errorf("%s unknown delim: %s", ce.Error(), del)
		}
	} else {
		return fmt.Errorf("%s expected a delim token", ce.Error())
	}

	headers := make([]string, 0)

	for decoder.More() {
		//uso any perché valori numerici/bool potrebbero essere soggetti a type inference e causare errore
		//se avessi string-string
		data := make(map[string]any)
		fields := make([]string, 0)
		decodeErr := decoder.Decode(&data)

		if decodeErr == io.EOF {
			return fmt.Errorf("%s stream ended unexpectedly", ce.Error())
		} else if decodeErr != nil {
			return fmt.Errorf("%s %w", ce.Error(), decodeErr)
		}

		//scrittura header
		if firstRecord {
			for k := range data {
				headers = append(headers, k)
			}
			headersErr := csvWriter.Write(headers)

			if headersErr != nil {
				return fmt.Errorf("%s %w", ce.Error(), headersErr)
			}

			firstRecord = false
		}

		for _, h := range headers {
			fields = append(fields, fmt.Sprintf("%v", data[h]))
		}

		writeErr := csvWriter.Write(fields)

		if writeErr != nil {
			return fmt.Errorf("%s %w", ce.Error(), writeErr)
		}
	}

	return nil
}
