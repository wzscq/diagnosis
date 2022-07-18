FROM wangzhsh/crvframe:0.1.0
RUN apk --update add openjdk11-jre
COPY nginx.conf /etc/nginx/nginx.conf
ADD ./package/web /web
ADD ./package/service /diagnosis
ADD ./package/DiagEngine /DiagEngine
copy entrypoint.sh /entrypoint.sh
ENTRYPOINT ["sh","entrypoint.sh"]