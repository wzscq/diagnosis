echo create folder for build package ...
if [ ! -e package ]; then
  mkdir package
fi

if [ ! -e package/web ]; then
  mkdir package/web
fi

echo build the code ...
cd ../diag_event_report
npm install
sed -i  's/host=\"*.*\"/host=\"\"/' ./public/index.html
npm run build
cd ../build

echo remove last pacakge if exist
if [ -e package/web/diag_event_report ]; then
  rm -rf package/web/diag_event_report
fi

mv ../diag_event_report/build package/web/diag_event_report

echo diag_event_report package build over.
