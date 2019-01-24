MultiExec executes multiple applications as configured.

The main reason this tool exists is because running PHP-FPM application as a single docker container is not enough, it still requires Nginx (or webserver of your choice).
To make them as an easily deployable container, Nginx & PHP-FPM need to be in the same container. Docker documentation has information on how to run multiple service in single container, but I found bash solution is too hard to read (just lazy, actually) and Supervisor is too much.

References:
- https://docs.docker.com/config/containers/multi-service_container/
- https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html