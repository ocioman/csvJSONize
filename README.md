# csvJSONize

## English

`csvJSONize` is a small Go project to convert:

- CSV → JSON (`Serialize`)
- JSON → CSV (`Deserialize`)

### Project structure

- `/converter/converter.go`: conversion logic
- `/main/main.go`: simple example entry point

### How to run

From the repository root:

```bash
go run ./main < ./main/MOCK_DATA.csv
```

This produces:

- `res.json` (serialized from CSV)
- `res.csv` (deserialized back from JSON)

## Italiano

`csvJSONize` è un piccolo progetto Go per convertire:

- CSV → JSON (`Serialize`)
- JSON → CSV (`Deserialize`)

### Struttura del progetto

- `/converter/converter.go`: logica di conversione
- `/main/main.go`: entry point di esempio

### Come eseguirlo

Dalla root del repository:

```bash
go run ./main < ./main/MOCK_DATA.csv
```

Questo comando genera:

- `res.json` (serializzato da CSV)
- `res.csv` (deserializzato da JSON)
