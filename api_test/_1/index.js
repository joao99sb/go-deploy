const http = require('http');
const url = require('url');
const WebSocket = require('ws')


const PORT = 3000;

const server = http.createServer(async (req, res) => {
  const parsedUrl = url.parse(req.url, true);
  const path = parsedUrl.pathname;

  if (req.method === 'GET' && path === '/') {

    res.writeHead(201, { 'Content-Type': 'text/plain; charset=utf-8' });


    const result = await new Promise((resolve, reject) => {

      setTimeout(() => {
        // Simula uma operação assíncrona bem-sucedida
        return resolve('Operação assíncrona concluída com sucesso');
      }, 100); // 1
    })

    res.end('Olá Mundo da primeira');
  } else if (req.method === 'POST' && path === '/test') {
    
    res.end(JSON.stringify({ message: 'pong' }));
  } else if (req.method === 'GET' && path === '/health') {


    res.writeHead(200, { 'Content-Type': 'application/json; charset=utf-8' });
    res.end(JSON.stringify({ ok: true }));
  } else if (req.method === 'GET' && path === '/html'){
    res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' });
    res.end('<DOCTYPE html><html><head><title>Olá Mundo</title></head><body><h1>Olá Mundo</h1></body></html>');
  }else {
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

setTimeout(()=>{server.listen(PORT, () => {
  console.log(`Servidor rodando na porta ${PORT}`);
});
},5000)

