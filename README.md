# go-csv
CSV files reader written in Golang.

_I wrote this package while I was learning Golang. It should work, but if you need a package to work with CSV files in Golang, please use this one: <https://golang.org/pkg/encoding/csv/>._

### Usage

#### Reading

Read CSV text with reader, for example:

```
csvfile, _ := os.Open(CSVFilePath)
csvReader := gocsv.NewReader(csvfile)
```

or 

```
csvStr := strings.NewReader("foo,bar,blah\nlol,cat,meow")
csvReader := gocsv.NewReader(csvStr)
```

Now you can read from `csvReader` until `Read` returns `EOF`:

```
var str string

csvReader := NewReader(strings.NewReader("foo,bar,baz\nnyan,cat,wat"))

for {
	rec, err := csvReader.Read()
	str += strings.Join(rec, "")

	if err == io.EOF {
		break
	}
}
```

#### Writing

To write you need to create a new Writer:

```
buf := bytes.NewBuffer([]byte{})
csvWriter := gocsv.NewWriter(buf)
```

or

```
file, _ := os.Create(CSVFilePath)
csvWriter := gocsv.NewWriter(file)
```

Now you can write slices of strings into this writer:

```
csvWriter.write([]string{"foo", "bar", "baz"})
csvWriter.write([]string{"nyan", "cat", "wat"})

```
