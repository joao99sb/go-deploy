const http = require('http');
const url = require('url');
const WebSocket = require('ws')

const PORT = 3000;

const server = http.createServer((req, res) => {
  const parsedUrl = url.parse(req.url, true);
  const path = parsedUrl.pathname;

  if (req.method === 'GET' && path === '/') {

    res.writeHead(201, { 'Content-Type': 'text/plain; charset=utf-8' });
    res.end('Olá Mundo da segunda');

  } else if (req.method === 'POST' && path === '/test') {

    res.end(JSON.stringify({ message: 'pong' }));
  } else if (req.method === 'GET' && path === '/health') {


    res.writeHead(200, { 'Content-Type': 'application/json; charset=utf-8' });
    res.end(JSON.stringify({ ok: true }));
  } else {
    res.writeHead(404, { 'Content-Type': 'text/plain; charset=utf-8' });
    res.end('Rota não encontrada');
  }
});
const wss = new WebSocket.Server({ server });

wss.on('connection', (ws) => {
  console.log('Cliente conectado');

  ws.on('message', (message) => {
    console.log('Recebido:', message);

    // Echo de volta para o cliente
    ws.send('Echo: ' + message);
  });

  ws.on('close', () => {
    console.log('Conexão fechada');
  });
});

server.listen(PORT, () => {
  console.log(`Servidor rodando na porta ${PORT}`);
});
