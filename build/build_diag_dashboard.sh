echo create folder for build package ...
if [ ! -e package ]; then
  mkdir package
fi

if [ ! -e package/web ]; then
  mkdir package/web
fi

echo build the code ...
cd ../diag_dashboard
npm install
sed -i  's/host=\"*.*\"/host=\"\"/' ./public/index.html
npm run build
cd ../build

echo remove last pacakge if exist
if [ -e package/web/diag_dashboard ]; then
  rm -rf package/web/diag_dashboard
fi

mv ../web/diag_dashboard/build package/web/diag_dashboard

echo diag_dashboard package build over.
