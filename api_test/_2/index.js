const http = require('http');
const url = require('url');

const PORT = 3000;

const server = http.createServer((req, res) => {
  const parsedUrl = url.parse(req.url, true);
  const path = parsedUrl.pathname;

  if (req.method === 'GET' && path === '/') {

    res.writeHead(201, { 'Content-Type': 'text/plain; charset=utf-8' });
    res.end('Olá Mundo da segunda');

  } else if (req.method === 'POST' && path === '/test') {

    res.end(JSON.stringify({ message: 'pong' }));
  } else if (req.method === 'GET' && path === '/ping') {


    res.writeHead(200, { 'Content-Type': 'application/json; charset=utf-8' });
    res.end(JSON.stringify({ message: 'pong 2' }));
  } else {
    res.writeHead(404, { 'Content-Type': 'text/plain; charset=utf-8' });
    res.end('Rota não encontrada');
  }
});

server.listen(PORT, () => {
  console.log(`Servidor rodando na porta ${PORT}`);
});
