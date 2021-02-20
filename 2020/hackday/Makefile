
.PHONY: download
download:
	curl -fsSL -o insert/jawiki.entity_vectors.100d.txt.bz2 https://github.com/singletongue/WikiEntVec/releases/download/20190520/jawiki.entity_vectors.100d.txt.bz2
	bzip2 -d insert/jawiki.entity_vectors.100d.txt.bz2

.PHONY: run
run:
	docker run -v $(PWD)/config/config.yaml:/etc/server/config.yaml -p 8081:8081 --rm -it vdaas/vald-agent-ngt

.PHONY: insert
insert:
	go run insert/main.go

.PHONY: search
search:
	go run search/main.go
