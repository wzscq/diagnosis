FROM wangzhsh/crvframe:0.1.0
RUN apk --update add openjdk11-jre
COPY nginx.conf /etc/nginx/nginx.conf
ADD ./package/web /web
ADD ./package/service/diagnosis_service /services/diagnosis/diagnosis_service
ADD ./package/service/DiagEngine/DiagEngine.jar /services/diagengine/DiagEngine.jar
copy entrypoint.sh /entrypoint.sh
ENTRYPOINT ["sh","entrypoint.sh"]