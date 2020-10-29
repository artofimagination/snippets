import defaults from 'lodash/defaults';

import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './DataSource';
import { defaultQuery, MyDataSourceOptions, MyQuery } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onDataTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, dataText: event.target.value });
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { dataText } = query;

    return (
      <div className="gf-form">
        <FormField
          labelWidth={8}
          value={dataText || ''}
          onChange={this.onDataTextChange}
          label="Data to show"
          tooltip="Data labels to be shown. Separate each name with a comma."
        />
      </div>
    );
  }
}
