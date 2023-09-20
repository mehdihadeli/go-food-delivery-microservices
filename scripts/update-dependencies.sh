# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

readonly service="$1"

echo "start upgrading packages in $service"

if [ "$service" = "pkg" ]; then
    cd "./internal/pkg" && go get -u -t -d -v ./... && go mod tidy
# Check if input is not empty or null
elif [ -n "$service"  ]; then
    cd "./internal/services/$service" && go get -u -t -d -v ./... && go mod tidy
fi
