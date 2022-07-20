nginx

cd /services/diagnosis
./diagnosis_service &

cd /services/crvframe
./frame &

cd /services/diagengine
java -jar ./DiagEngine.jar