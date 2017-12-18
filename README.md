# gopal

Simple readonly side data storage for Go

Inspired by https://github.com/linkedin/PalDB.

Minimal perfect hash + mmaped data.

## CSV to pal

`data.csv`:

```csv
id,f1,f2
0,v11,v12
1,v21,v22
```

```
$ gopal data.csv data.pal ,
```

## Versions

> There are bug in one of [dependencies](https://github.com/alecthomas/mph/issues/10)

For compatibility Gopal can read serialized data from two different formats: 
`V1` and `V2`. Data format is determined from header and does not require any 
special settings.

Gopal always builds and writes data in `V2`.

