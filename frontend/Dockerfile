# Используем стабильную версию Node.js
FROM node:18-alpine

WORKDIR /app

# Копируем зависимости
COPY package*.json ./
RUN npm install

# Копируем весь код фронтенда
COPY . .

# Открываем порт React
EXPOSE 3000

# Запускаем проект в режиме разработки
CMD ["npm", "start"]