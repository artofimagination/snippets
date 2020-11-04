# Database and backend server load tester utility
This small example is a functioning load tester. 
Main components:
 - python [locust](https://locust.io/) load generation. This produces REST API requests that are sent to the backend
 - golang backend server. Connects to DB and grafana. Receives locust requests
 - Test Databases. 
 - Grafana dashboard
 
 # How it works?
 - Once the backend system is started locust can be initiated. It will start sending REST requests (at the mment DB insert queries) to the backend.
 - The backend will execute the query in the database or send back a failure response if the request has failed
 - Meanwhile if grafana board has been opened it will request data streaming from the backend
 - After the request has been processed the backend will stream grafana query data infinitely back to the grafana server
 
 Grafana dashboard and backend is setup to display
  - CPU/memory load
  - IO wait time
  - single insert query exec time
  - overall runtime of insert sequence. If measurement of total time  of 1 million inserts is needed for example
  - request rate (requests/second)
  - failed request count
  - row count
  
  # How to use?
  - To start the backend run ```docker-compose up --build --force-recreate -d main-server```
  - To start locust install python 3 and run ```pip3 install -r requirements.txt```. This will install the modules the script needs
  - Run locust master ```locust -f loadTester.py --master```
  - Run locust workers ```locust -f loadTester.py --worker --master-host=localhost```
  - Once all these are running both grafana and locust can be accessed through ```http://localhost:8080/show```
    - grafana may need an initial login beforehand through ```http://localhost:3000```
    - locust can be access alternatively through ```http://localhost:8089```
  - Use locust config to setup spawned users and start streaming
  
  # Additional info
  The tool is incomplete see [grafana-streamed-charting](https://github.com/artofimagination/snippets/blob/master/grafana-streamed-charts/README.md) for known issues.
  
    
