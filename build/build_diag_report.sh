echo create folder for build package ...
if [ ! -e package ]; then
  mkdir package
fi

if [ ! -e package/web ]; then
  mkdir package/web
fi

echo build the code ...
cd ../diag_report
npm install
sed -i  's/host=\"*.*\"/host=\"\"/' ./public/index.html
npm run build
cd ../build

echo remove last pacakge if exist
if [ -e package/web/diag_report ]; then
  rm -rf package/web/diag_report
fi

mv ../diag_report/build package/web/diag_report

echo diag_report package build over.
