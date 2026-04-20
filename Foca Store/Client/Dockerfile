FROM node:20-alpine

WORKDIR /app

# Copy daftar library
COPY package*.json ./

# Instal library (ini bakal cepat karena sudah ada cache dari yang tadi)
RUN npm install

# Copy semua kode
COPY . .

# Beritahu Docker kita pakai port 3000
EXPOSE 3000

# KITA HAPUS BARIS 'RUN npm run build'
# Langsung jalankan mode dev saja
CMD ["npm", "run", "dev"]