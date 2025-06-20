FROM node:18-alpine

# Create non-root user
RUN adduser -D -s /bin/sh nodeuser

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

# Change ownership and switch to non-root user
RUN chown -R nodeuser:nodeuser /app
USER nodeuser

EXPOSE 3000

CMD ["npm", "run", "dev"]
