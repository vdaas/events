# hackday

## Vald Agent の起動

```sh
$ make run
```

## 検索対象データの追加

```sh
$ cd insert
$ wget https://github.com/singletongue/WikiEntVec/releases/download/20190520/jawiki.entity_vectors.100d.txt.bz2
$ bzip2 -d jawiki.entity_vectors.100d.txt.bz2
$ go run main.go
```

## 検索

```sh
$ cd search
$ go run main.go -q Google
```