const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://backend:8080',
      changeOrigin: true,
      pathRewrite: {
        '^/api': '/api'
      },
      onError: (err, req, res) => {
        console.error('Proxy error:', err);
      }
    })
  );
};
