#!/usr/bin/env sh

cmdname=$(basename $0)

usage()
{
    cat << USAGE >&2
Usage:
    ./$cmdname -h host -p port [-c command] [-t timeout] [-q quiet]
    -h          Host or IP under test
    -p          TCP port under test
    -c          Command to execute if test succeeds
    -t          Timeout in seconds, zero for no timeout
    -q          Don't output any status messages
USAGE
    exit 1
}

# process arguments
while [[ $# -gt 0 ]]
do
    case "$1" in
        -h)
        HOST="$2"
        shift 2
        ;;
        -p)
        PORT="$2"
        shift 2
        ;;
        -c)
        COMMAND="$2"
        shift 2
        ;;
        -t)
        TIMEOUT="$2"
        shift 2
        ;;
        -q)
        QUIET=1
        shift 1
        ;;
        --help)
        usage
        ;;
        *)
        echoerr "Unknown argument: $1"
        usage
        ;;
    esac
done

if [[ "$HOST" == "" || "$PORT" == "" ]]; then
    echo "Error: you need to provide a host and port to test."
    usage
fi

TIMEOUT=${TIMEOUT:-0}
QUIET=${QUIET:-0}

wait=1
step=1
run=0

if [[ ! $QUIET -eq 1 ]]; then
    echo Testing $HOST:$PORT
fi

while :
do
    result=$(nc $HOST $PORT -w $wait -z &> /dev/null; echo $?)
    if [ "$result" -eq "0" ]; then
        run=1
        break
    fi
    
    if [[ ! $QUIET -eq 1 ]]; then
        echo $step
    fi

    sleep 1

    if [ $step -eq $TIMEOUT ] && [ ! $TIMEOUT -eq 0 ]; then
        break
    fi

    step=$((step+1))
done

if [ $run -eq 1 ]; then
    if [[ ! $QUIET -eq 1 ]]; then
        echo Sucess
    fi

    if [[ ! "$COMMAND" == "" ]]; then
        exec $COMMAND
    fi
fi
