# Example snippet for grafana streamed json datasource from http backend
This example represents an http backend server that streams json data to a locally hosted grafana server.
Very crude example which is still heavily under development.
Examples used:
[Build streaming datasource for grafana](https://grafana.com/docs/grafana/latest/developers/plugins/build-a-streaming-data-source-plugin/)
[Simple streaming datasrouce](https://github.com/seanlaff/simple-streaming-datasource)

# How it works?
1. The example consist of an http backend and a grafana docker container.
2. When the example DataStream Dashboard is opened the datasource will request the streaming from the backend
3. Once the response starts to arrive to the data source the data will be displayed in the charts

# How to run?
## Setup
- Run ```docker-compose up --build --force-recreate -d main-server```
- Access grafana on `http://localhost:3000`
- Set data-stream datasource with the host address `http://localhost:8080`
- Import the demo dashboard from under the Dashboards tab

## Running
- Once the setup steps are done the line chart panel can be accessed from the backend through `http://localhost:8080/show`
- The data name can be changed in the panel edit page under the dataText field. If multiple series need to be shown, names need to be separated by a comma (For example: testData1,testData2).

## Build the datasource
- navigate in `./grafana/json-data-stream`
- run `yarn install` (may requires to install some linux node.js packages)
- run `yarn dev`

# Known issues
This example is meant to be a POC and it is in a very early stage yet. 
1. Only handles panels in a single dashboard

2. Data display gets slow after a while

3. The panel visualization is cleared every time you update the dashboard. If you have access to historical data, you can add, or backfill, it to the data frame before the first call to subscriber.next().

4. There are building issues with the datasource (the fail does not prevent the usage of the datasource)
```
peter@peter-GL63-8RCS:~/Git/snippets/grafana-charting/json-data-stream$ yarn dev
yarn run v1.22.5
$ grafana-toolkit plugin:dev
✔ Linting
⠋ Bundling plugin in dev mode  Starting type checking service...
  Using 1 worker with 2048MB memory limit
⠴ Bundling plugin in dev mode  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts(27,3):
  TS2416: Property 'query' in type 'DataSource' is not assignable to the same property in base type 'DataSourceApi<MyQuery, MyDataSourceOptions>'.
    Type '(request: DataQueryRequest<MyQuery>) => Observable<DataQueryResponse>' is not assignable to type '(request: DataQueryRequest<MyQuery>) => Promise<DataQueryResponse> | Observable<DataQueryResponse>'.
      Type 'Observable<DataQueryResponse>' is not assignable to type 'Promise<DataQueryResponse> | Observable<DataQueryResponse>'.
        Type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/rxjs/internal/Observable").Observable<import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/types/datasource").DataQueryResponse>' is not assignable to type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/node_modules/rxjs/internal/Observable").Observable<import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/types/datasource").DataQueryResponse>'.
          The types of 'source.operator.call' are incompatible between these types.
            Type '(subscriber: import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/rxjs/internal/Subscriber").Subscriber<any>, source: any) => import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/rxjs/internal/types").TeardownLogic' is not assignable to type '(subscriber: import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/node_modules/rxjs/internal/Subscriber").Subscriber<any>, source: any) => import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/node_modules/...'.
              Types of parameters 'subscriber' and 'subscriber' are incompatible.
                Type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/node_modules/rxjs/internal/Subscriber").Subscriber<any>' is not assignable to type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/rxjs/internal/Subscriber").Subscriber<any>'.
                  Property 'isStopped' is protected but type 'Subscriber<T>' is not a class derived from 'Subscriber<T>'.
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts(57,17):
  TS7034: Variable 'readHandler' implicitly has type 'any' in some locations where its type cannot be determined.
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts(59,30):
  TS7006: Parameter 'result' implicitly has an 'any' type.
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts(65,38):
  TS7006: Parameter 'element' implicitly has an 'any' type.
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts(75,41):
  TS7005: Variable 'readHandler' implicitly has an 'any' type.
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/module.ts(7,44):
  TS2344: Type 'DataSource' does not satisfy the constraint 'DataSourceApi<MyQuery, MyDataSourceOptions>'.
    The types returned by 'query(...)' are incompatible between these types.
      Type 'Observable<DataQueryResponse>' is not assignable to type 'Promise<DataQueryResponse> | Observable<DataQueryResponse>'.
        Type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/rxjs/internal/Observable").Observable<import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/types/datasource").DataQueryResponse>' is not assignable to type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/node_modules/rxjs/internal/Observable").Observable<import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/types/datasource").DataQueryResponse>'.
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/module.ts(8,20):
  TS2345: Argument of type 'typeof ConfigEditor' is not assignable to parameter of type 'ComponentType<DataSourcePluginOptionsEditorProps<MyDataSourceOptions, unknown>>'.
    Type 'typeof ConfigEditor' is not assignable to type 'ComponentClass<DataSourcePluginOptionsEditorProps<MyDataSourceOptions, unknown>, any>'.
      Types of parameters 'props' and 'props' are incompatible.
        Type 'DataSourcePluginOptionsEditorProps<MyDataSourceOptions, unknown>' is not assignable to type 'Readonly<Props>'.
          Types of property 'options' are incompatible.
            Type 'DataSourceSettings<MyDataSourceOptions, unknown>' is not assignable to type 'DataSourceSettings<MyDataSourceOptions, {}>'.
              Type 'unknown' is not assignable to type '{}'.
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/QueryEditor.tsx(11,31):
  TS2344: Type 'DataSource' does not satisfy the constraint 'DataSourceApi<MyQuery, MyDataSourceOptions>'.
    The types returned by 'query(...)' are incompatible between these types.
      Type 'Observable<DataQueryResponse>' is not assignable to type 'Promise<DataQueryResponse> | Observable<DataQueryResponse>'.
        Type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/rxjs/internal/Observable").Observable<import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/types/datasource").DataQueryResponse>' is not assignable to type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/node_modules/rxjs/internal/Observable").Observable<import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/types/datasource").DataQueryResponse>'.
  
   Hash: 55255503148a1f4648a0
  Version: webpack 4.41.5
  Time: 4037ms
  Built at: 10/29/2020 5:37:29 PM
                       Asset       Size  Chunks                   Chunk Names
                     LICENSE   11.1 KiB          [emitted]        
                   README.md   1.33 KiB          [emitted]        
  dashboards/DataStream.json    7.2 KiB          [emitted]        
                img/logo.svg   1.55 KiB          [emitted]        
                   module.js   86.3 KiB  module  [emitted]        module
               module.js.map   80.5 KiB  module  [emitted] [dev]  module
        partials/config.html   78 bytes          [emitted]        
  partials/query.editor.html  468 bytes          [emitted]        
                 plugin.json  952 bytes          [emitted]        
  Entrypoint module = module.js module.js.map
  [../node_modules/lodash/_baseRest.js] 559 bytes {module} [built]
  [../node_modules/lodash/_isIterateeCall.js] 877 bytes {module} [built]
  [../node_modules/lodash/defaults.js] 1.71 KiB {module} [built]
  [../node_modules/lodash/eq.js] 799 bytes {module} [built]
  [../node_modules/tslib/tslib.es6.js] 10 KiB {module} [built]
  [./ConfigEditor.tsx] 3.88 KiB {module} [built]
  [./DataSource.ts] 3.28 KiB {module} [built]
  [./QueryEditor.tsx] 1.2 KiB {module} [built]
  [./module.ts] 296 bytes {module} [built]
  [./types.ts] 45 bytes {module} [built]
  [./vendor/ndjson.js] 1.82 KiB {module} [built]
  [@grafana/data] external "@grafana/data" 42 bytes {module} [built]
  [@grafana/ui] external "@grafana/ui" 42 bytes {module} [built]
  [react] external "react" 42 bytes {module} [built]
  [rxjs] external "rxjs" 42 bytes {module} [built]
      + 43 hidden modules
  
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts(27,3):
  TS2416: Property 'query' in type 'DataSource' is not assignable to the same property in base type 'DataSourceApi<MyQuery, MyDataSourceOptions>'.
    Type '(request: DataQueryRequest<MyQuery>) => Observable<DataQueryResponse>' is not assignable to type '(request: DataQueryRequest<MyQuery>) => Promise<DataQueryResponse> | Observable<DataQueryResponse>'.
      Type 'Observable<DataQueryResponse>' is not assignable to type 'Promise<DataQueryResponse> | Observable<DataQueryResponse>'.
        Type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/rxjs/internal/Observable").Observable<import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/types/datasource").DataQueryResponse>' is not assignable to type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/node_modules/rxjs/internal/Observable").Observable<import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/types/datasource").DataQueryResponse>'.
          The types of 'source.operator.call' are incompatible between these types.
            Type '(subscriber: import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/rxjs/internal/Subscriber").Subscriber<any>, source: any) => import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/rxjs/internal/types").TeardownLogic' is not assignable to type '(subscriber: import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/node_modules/rxjs/internal/Subscriber").Subscriber<any>, source: any) => import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/node_modules/...'.
              Types of parameters 'subscriber' and 'subscriber' are incompatible.
                Type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/node_modules/rxjs/internal/Subscriber").Subscriber<any>' is not assignable to type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/rxjs/internal/Subscriber").Subscriber<any>'.
                  Property 'isStopped' is protected but type 'Subscriber<T>' is not a class derived from 'Subscriber<T>'.
  
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts(57,17):
  TS7034: Variable 'readHandler' implicitly has type 'any' in some locations where its type cannot be determined.
  
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts(59,30):
  TS7006: Parameter 'result' implicitly has an 'any' type.
  
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts(65,38):
  TS7006: Parameter 'element' implicitly has an 'any' type.
  
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/DataSource.ts(75,41):
  TS7005: Variable 'readHandler' implicitly has an 'any' type.
  
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/module.ts
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/module.ts(7,44):
  TS2344: Type 'DataSource' does not satisfy the constraint 'DataSourceApi<MyQuery, MyDataSourceOptions>'.
    The types returned by 'query(...)' are incompatible between these types.
      Type 'Observable<DataQueryResponse>' is not assignable to type 'Promise<DataQueryResponse> | Observable<DataQueryResponse>'.
        Type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/rxjs/internal/Observable").Observable<import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/types/datasource").DataQueryResponse>' is not assignable to type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/node_modules/rxjs/internal/Observable").Observable<import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/types/datasource").DataQueryResponse>'.
  
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/module.ts
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/module.ts(8,20):
  TS2345: Argument of type 'typeof ConfigEditor' is not assignable to parameter of type 'ComponentType<DataSourcePluginOptionsEditorProps<MyDataSourceOptions, unknown>>'.
    Type 'typeof ConfigEditor' is not assignable to type 'ComponentClass<DataSourcePluginOptionsEditorProps<MyDataSourceOptions, unknown>, any>'.
      Types of parameters 'props' and 'props' are incompatible.
        Type 'DataSourcePluginOptionsEditorProps<MyDataSourceOptions, unknown>' is not assignable to type 'Readonly<Props>'.
          Types of property 'options' are incompatible.
            Type 'DataSourceSettings<MyDataSourceOptions, unknown>' is not assignable to type 'DataSourceSettings<MyDataSourceOptions, {}>'.
              Type 'unknown' is not assignable to type '{}'.
  
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/QueryEditor.tsx
  ERROR in /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/src/QueryEditor.tsx(11,31):
  TS2344: Type 'DataSource' does not satisfy the constraint 'DataSourceApi<MyQuery, MyDataSourceOptions>'.
    The types returned by 'query(...)' are incompatible between these types.
      Type 'Observable<DataQueryResponse>' is not assignable to type 'Promise<DataQueryResponse> | Observable<DataQueryResponse>'.
        Type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/rxjs/internal/Observable").Observable<import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/types/datasource").DataQueryResponse>' is not assignable to type 'import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/node_modules/rxjs/internal/Observable").Observable<import("/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/data/types/datasource").DataQueryResponse>'. 
  
  Trace: Build failed
      at /home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/toolkit/src/cli/utils/useSpinner.js:25:29
      at step (/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/toolkit/node_modules/tslib/tslib.js:140:27)
      at Object.throw (/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/toolkit/node_modules/tslib/tslib.js:121:57)
      at rejected (/home/peter/Git/snippets/grafana-streamed-charts/grafana/json-data-stream/node_modules/@grafana/toolkit/node_modules/tslib/tslib.js:112:69)
      at processTicksAndRejections (internal/process/task_queues.js:97:5)
✖ Build failed
error Command failed with exit code 1.
info Visit https://yarnpkg.com/en/docs/cli/run for documentation about this command.
```
