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
  resolution: number;
  sourceAddress: string;
  reader: ReadableStreamDefaultReader;
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
    this.resolution = instanceSettings.jsonData.resolution || 1000.0;
    this.sourceAddress = instanceSettings.jsonData.address || '';
    this.reader = new ReadableStream().getReader();
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
        frame.addField({ name: 'time', type: FieldType.time });
        frame.addField({ name: 'value', type: FieldType.number });

        var request = new Request(`${this.sourceAddress}?numSeries=1`);
        fetch(request)
          .then(response => {
            // In the real world its likely that our json gets chopped into
            // chunks when streamed from the backend. ndjsonStream handles
            // reconstructing the newline-delimmited json for us.
            return ndjsonStream.default(response.body);
          })
          .then(s => {
            this.reader = s.getReader(); // Save the reader so we can cancel it later
            let readHandler;
            this.reader.read().then(
              (readHandler = result => {
                if (result.done) {
                  this.reader.cancel();
                  return;
                }
                frame.add({ time: result.value.timestamp, value: result.value.value });

                subscriber.next({
                  data: [frame],
                  key: query.refId,
                });
                this.reader.read().then(readHandler);
              })
            );
          });
        const intervalId = setInterval(() => {}, 100);

        return () => {
          this.reader.cancel();
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
