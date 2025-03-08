# BusyHttp

This is a basic image for debugging kubernetes/container setups. It is based on busybox (hence the name), but it adds
some capabilities via a convenient http api.

You can get it from https://hub.docker.com/repository/docker/elmendez/busyhttp

| Key                       | Description                                                                           |
|---------------------------|---------------------------------------------------------------------------------------|
| **BUSY_STARTUP_TIME_MS**  | The startup time for the http server                                                  |
| **BUSY_CRASH**            | Immediately crashes after startup time. Anything other than 0 activates it.           |
| **BUSY_SHUTDOWN_TIME_MS** | The time after SIGINT needed to start shutting down the http server                   |
| **BUSY_TRUSTED_PROXIES**  | Trusted proxies CIDR                                                                  |
| **BUSY_READY_TIME_MS**    | Time after the http server is up needed for the ready endpoint to return 200.         |
| **BUSY_TRUSTED_PLATFORM** | The header for IP mapping with trusted proxies                                        |
| **BUSY_SECRET**           | If auth is needed, you can set this variable and it must be sent in the Bearer format |
| **BUSY_ADDRESS**          | The server address, default :8080                                                     |

| Endpoint                     | Description                                                      |
|------------------------------|------------------------------------------------------------------|
| **/ping**                    | Returns 200 Ok                                                   |
| **/help**                    | Shows this message                                               |
| **/ready**                   | Returns 200 (depends on BUSY_READY_TIME_MS)                      |
| **/info**                    | Returns info about this server                                   |
| **/time**                    | Returns info about time data                                     |
| **/echo**                    | Returns info about the request received                          |
| **/crash**                   | Immediately exits with 1.                                        |
| **/wait/:ms**                | Waits for :ms and then returns 200.                              |
| **/exit**                    | Immediately exits with 0                                         |
| **/status/:code**            | Returns the :code status code                                    |
| **/file?filename=:filename** | Gets or creates a file                                           |
| **/specify/:any/Endpoint**   | A reexports of the same endpoints as an somewhat arbitrary route |
