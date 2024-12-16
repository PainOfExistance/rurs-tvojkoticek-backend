docker container stop tvoj_koticek_backend_container || true
docker container rm tvoj_koticek_backend_container || true

docker build -t tvoj_koticek_backend .

docker run --name tvoj_koticek_backend_container -p 8080:8080 tvoj_koticek_backend