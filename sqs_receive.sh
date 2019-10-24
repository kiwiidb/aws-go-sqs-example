#!/bin/bash

#queue url should be first arg
QURL=$1

while true; do
MSG=$(aws sqs receive-message --queue-url $QURL)
if [ ! -z "$MSG" ]; then
echo "$MSG" | jq -r '.Messages[] | .ReceiptHandle' \
| xargs -I {} aws sqs delete-message --queue-url $QURL --receipt-handle {}
echo $"$MSG" | jq -r '.Messages[] | .Body'
fi
done

