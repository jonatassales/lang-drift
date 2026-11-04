# The console and extub in one image: extub on 127.0.0.1:3401, the console
# on 0.0.0.0:3001. Illustrative: the planner reads the port and the
# readiness, it never builds this file.
FROM node:24-alpine
WORKDIR /srv
COPY . .
EXPOSE 3001
HEALTHCHECK --interval=5s --timeout=2s --retries=12 CMD wget -q --spider http://127.0.0.1:3001/icon.svg
CMD ["sh", "-c", "extub & node apps/console/server.js"]
