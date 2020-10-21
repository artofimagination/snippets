# Timescale DB example
This example meant to test basic functionalities of postgres TimescaleDB. The example allows to store some json data in a project_data table. Functionalities to delete and query rows are also implemented
The table structure looks as follows

```
project_id {
    created_at -> timestamp,
    project_id -> UUID,
    run_seq_no -> integer,
    data -> JSON Blob
}
```

# Usage
## Build
Run ```docker-compose up --build --force-recreate -d main-server``` to generate and start all containers

In order to access the db run:
- ```docker exec -it user-data-db bash -c "psql data root"```

## Execution examples
Use the browser or ```curl``` command to execute the following
- add new data with predefined project UUID: ```http://localhost:8080/insert?project=408c57ad-134c-11eb-ab0c-0242ac120003&seqNo=1&data={"test":"1"}```
- add new data with random project UUID: ```http://localhost:8080/insert?seqNo=1&data={"test":"1"}```
- get a chunk of data belonging to a specific project run starting at the defined date: ```http://localhost:8080/get-by-project-run?project=408c57ad-134c-11eb-ab0c-0242ac120003&seqNo=1&chunk=3&date=2020-05-07```

At the moment there is no support to search by datetime, but only by date.
- delete data by project run: ```http://localhost:8080/delete-by-project-run?project=408c57ad-134c-11eb-ab0c-0242ac120003&seqNo=1```
- delete data by project: ```http://localhost:8080/delete-by-project?project=408c57ad-134c-11eb-ab0c-0242ac120003```
