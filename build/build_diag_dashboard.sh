cd ..

echo create folder for build package ...
if [ ! -e package ]; then
  mkdir package
fi

if [ ! -e package/web ]; then
  mkdir package/web
fi

echo build the code ...
cd web/diag_dashboard
npm install
npm run build
cd ../..

echo remove last pacakge if exist
if [ -e package/web/diag_dashboard ]; then
  rm -rf package/web/diag_dashboard
fi

mv web/diag_dashboard/build package/web/diag_dashboard

echo diag_dashboard package build over.
