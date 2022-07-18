nginx

cd /diagnosis
./diagnosis_service &

cd ../frame_service
./frame &

cd ../diagnosis/DiagEngine
java -jar ./DiagEngine.jar