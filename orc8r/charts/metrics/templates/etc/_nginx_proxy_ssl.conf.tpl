server {
  listen 443 ssl;
  server_name {{ .Values.prometheus.nginx.endpoints }}; 
  ssl_certificate /etc/nginx/conf.d/{{ .Values.prometheus.nginx.spec.ssl_cert_name }};
  ssl_certificate_key /etc/nginx/conf.d/{{ .Values.prometheus.nginx.spec.ssl_key_name }};
  ssl_verify_client on; 
  ssl_client_certificate /etc/nginx/conf.d/{{ .Values.prometheus.nginx.spec.ssl_certificate }};
  location / { 
     proxy_pass http://orc8r-prometheus:9090;
     proxy_set_header Host $http_host;
     proxy_set_header X-Forwarded-Proto $scheme;
  }
}


