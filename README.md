#Api-Router application

Application analyse the incoming api key provided with the user,\
and according to the algorithm adds the corresponding http header,\
which Nginx uses for routing connection to specific server.

Config example is :
    config.toml

Api key example:
    "x-api" : "BQY123456789a123456789b123456789"