nginx

cd /diagnosis
./diagnosis_service &

cd ../frame_service
./frame &

cd ../DiagEngine
java -jar ./DiagEngine.jar