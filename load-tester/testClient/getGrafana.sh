git clone https://github.com/artofimagination/grafana-json-streaming-datasource ./grafana-streaming
cp -r ./grafana-streaming/grafana .
rm -frv ./grafana-streaming

cd grafana/json-data-stream
# May require to install nodejs related packages (npm package).
# At the moment the plugin has to be built from source, there is no downloadable version.
yarn install
yarn dev