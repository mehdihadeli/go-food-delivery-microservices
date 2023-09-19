#!/bin/bash

# https://github.com/golang-migrate/migrate/tree/856ea12df9d230b0145e23d951b7dbd6b86621cb/cmd/migrate#usage

echo "Running go migrate command..."
while getopts ":p:c:n:o:" opt; do
  case $opt in
    p)
      path="$OPTARG"
      ;;
    c)
      command="$OPTARG"
      ;;
    n)
      name="$OPTARG"
      ;;
    o)
      connection="$OPTARG"
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      exit 1
      ;;
  esac
done

if [ "$command" == "create" ]; then
  migrate create -ext sql -dir "$path" -seq "$name"
elif [ "$command" == "up" ]; then
  migrate -database "$connection" -verbose -path "$path" up
elif [ "$command" == "down" ]; then
  migrate -database "$connection" -verbose -path "$path" down
else
  echo "Invalid command. Supported commands: create, up, down."
fi
