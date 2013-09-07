#!/bin/bash
host=$(grep -A3 elasticsearch graphite-ng.conf | sed -n 's/^host = "\(.*\)"/\1/p')
port=$(grep -A3 elasticsearch graphite-ng.conf | sed -n 's/^port = \(.*\)/\1/p')
index=carbon-es

echo "elasticsarch server: http://$host:$port"
echo "delete existing index $index (maybe)"
curl -X DELETE http://$host:$port/$index
echo
echo "create index $index"
curl -XPOST http://$host:$port/$index -d '{
    "settings" : {
        "number_of_shards" : 1
    },
    "mappings" : {
        "datapoint" : {
            "_source" : { "enabled" : true },
            "_id": {"index": "not_analyzed", "store" : "yes"},
            "properties" : {
                "metric": {"type": "string"},
                "ts": {"type": "integer"},
                "value": {"type": "float"}
            }
        }
    }
}'
echo
