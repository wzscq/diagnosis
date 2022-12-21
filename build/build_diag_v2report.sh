echo create folder for build package ...
if [ ! -e package ]; then
  mkdir package
fi

if [ ! -e package/web ]; then
  mkdir package/web
fi

echo build the code ...
cd ../diag_v2report
npm install
sed -i  's/host=\"*.*\"/host=\"\"/' ./public/index.html
npm run build
cd ../build

echo remove last pacakge if exist
if [ -e package/web/diag_v2report ]; then
  rm -rf package/web/diag_v2report
fi

mv ../diag_v2report/build package/web/diag_v2report

echo diag_v2report package build over.
