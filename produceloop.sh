#!/bin/bash
for i in {1..20}; do 
aws sqs send-message --queue-url https://eu-central-1.queue.amazonaws.com/684464168616/test_queue --message-body "hello there$i"
done
