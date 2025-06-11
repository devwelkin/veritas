# --- 1. Aşama: Build Aşaması ---
# Go'nun kurulu olduğu, derleme için optimize edilmiş bir imajla başlıyoruz.
FROM golang:1.24-alpine AS builder

# Hangi uygulamayı derleyeceğimizi build sırasında dışarıdan alacağız.
ARG APP_NAME

# Uygulama kaynak kodları için bir çalışma dizini oluşturuyoruz.
WORKDIR /app

# Önce sadece modül dosyalarını kopyalayıp bağımlılıkları indiriyoruz.
# Bu katman, sadece bağımlılıklar değiştiğinde yeniden çalışır, build'leri hızlandırır.
COPY go.mod go.sum ./
RUN go mod download

# Tüm kaynak kodunu kopyalıyoruz.
COPY . .

# Uygulamamızı derliyoruz.
# CGO_ENABLED=0, statik bir binary oluşturarak minimal imajlarda sorunsuz çalışmasını sağlar.
# Çıktıyı /app/bin dizinine, argüman olarak gelen uygulama adıyla kaydediyoruz.
RUN CGO_ENABLED=0 go build -o /app/bin/${APP_NAME} ./cmd/${APP_NAME}

# --- 2. Aşama: Final Aşaması ---
# Sadece uygulamayı çalıştırmak için gerekenleri içeren, ultra minimal bir alpine imajı.
FROM alpine:latest

# Güvenlik güncellemelerini alalım.
RUN apk --no-cache add ca-certificates

# Derlenmiş uygulamamızı 'builder' aşamasından bu imajın içine kopyalıyoruz.
COPY --from=builder /app/bin/* /usr/local/bin/

# Container çalıştığında hangi komutun çalışacağını belirtiyoruz.
# Docker Compose dosyasından gelen command bu satırı ezecek ama varsayılan olarak kalması iyidir.
CMD ["/usr/local/bin/creator"]