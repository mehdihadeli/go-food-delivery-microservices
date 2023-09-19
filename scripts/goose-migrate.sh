#!/bin/bash

#!/bin/bash


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
  goose -s -v -dir "$path" create "$name" sql
elif [ "$command" == "up" ]; then
  goose -s -v -dir "$path" postgres "$connection" up
elif [ "$command" == "down" ]; then
  goose -s -v -dir "$path" postgres "$connection" down
else
  echo "Invalid command. Supported commands: create, up, down."
fi
