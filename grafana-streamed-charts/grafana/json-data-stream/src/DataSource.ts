import defaults from 'lodash/defaults';
import * as ndjsonStream from './vendor/ndjson.js';

import {
  DataQueryRequest,
  DataQueryResponse,
  DataSourceApi,
  DataSourceInstanceSettings,
  CircularDataFrame,
  FieldType,
} from '@grafana/data';

import { MyQuery, MyDataSourceOptions, defaultQuery } from './types';
import { Observable, merge } from 'rxjs';

export class DataSource extends DataSourceApi<MyQuery, MyDataSourceOptions> {
  sourceAddress: string;
  readers: Map<string, ReadableStreamDefaultReader>;
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
    this.sourceAddress = instanceSettings.jsonData.address || '';
    this.readers = new Map();
  }

  query(request: DataQueryRequest<MyQuery>): Observable<DataQueryResponse> {
    const streams = request.targets.map(target => {
      const query = defaults(target, defaultQuery);

      return new Observable<DataQueryResponse>(subscriber => {
        const frame = new CircularDataFrame({
          append: 'tail',
          capacity: 1000,
        });

        frame.refId = query.refId;
        frame.addField({ name: 'timestamp', type: FieldType.time });
        var dataRows = query.dataText.split(',');
        dataRows.forEach(element => {
          frame.addField({ name: element, type: FieldType.number });
        });

        var restRequest = new Request(
          `${this.sourceAddress}?panelid=${request.panelId}&refid=${query.refId}&data-rows=${query.dataText}&start=${request.range.from}&end=${request.range.to}&datapoints=${request.maxDataPoints}`
        );
        fetch(restRequest)
          .then(response => {
            // In the real world its likely that our json gets chopped into
            // chunks when streamed from the backend. ndjsonStream handles
            // reconstructing the newline-delimmited json for us.
            return ndjsonStream.default(response.body);
          })
          .then(s => {
            this.readers.set(`${request.panelId}-${query.refId}`, s.getReader()); // Save the reader so we can cancel it later
            let readHandler;
            if (this.readers.has(`${request.panelId}-${query.refId}`)) {
              var reader = this.readers.get(`${request.panelId}-${query.refId}`) || new ReadableStream().getReader();
              reader.read().then(
                (readHandler = result => {
                  if (result.done) {
                    reader.cancel();
                    return;
                  }
                  result.value.forEach(element => {
                    if (request.panelId === element.panelid && query.refId === element.refid) {
                      frame.add(element.values);

                      subscriber.next({
                        data: [frame],
                        key: query.refId,
                      });
                    }
                  });
                  reader.read().then(readHandler);
                })
              );
            }
          });
        const intervalId = setInterval(() => {}, 100);

        return () => {
          var reader = this.readers.get(`${request.panelId}-${query.refId}`) || new ReadableStream().getReader();
          reader.cancel();
          clearInterval(intervalId);
        };
      });
    });
    return merge(...streams);
  }

  async testDatasource() {
    // Implement a health check for your data source.
    return {
      status: 'success',
      message: this.sourceAddress + ' success',
    };
  }
}
